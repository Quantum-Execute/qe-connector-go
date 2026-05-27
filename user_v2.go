package qe_connector

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
	"github.com/Quantum-Execute/qe-connector-go/handlers"
)

// V2 endpoints under `/strategy-api/user/.../v2/...` (the `/strategy-api`
// prefix is added by the gateway so the SDK keeps the V1 convention of
// writing `/user/...` here).
const (
	v2ExchangeApisEndpoint   = "/user/exchange/v2/exchange-apis"
	v2MasterOrdersEndpoint   = "/user/trading/v2/master-orders"
	v2OrderFillsEndpoint     = "/user/trading/v2/order-fills"
	v2TCAAnalysisEndpoint    = "/user/trading/v2/tca-analysis"
	v2BatchCancelEndpoint    = "/user/trading/v2/master-orders/batch-cancel"
	v2MasterOrdersByClientId = "/user/trading/v2/master-orders/by-client-order-id"
)

// MasterOrderStatusV2 enumerates the V2 master order statuses documented in
// `frontend-v2-api-upgrade.md` §2.
//
// 母单状态在 V2 接口里的语义按用途分两类：
//
//  1. 详情 / 推送返回：保留全部 9 个细分状态（NEW / WAITING / PROCESSING /
//     PAUSED / CANCELLED / COMPLETED / COMPLETED_WITHTAIL / REJECTED /
//     EXPIRED），SDK 接收端按需展示。
//  2. 列表查询过滤（GetMasterOrdersV2Service.Status）：后端只接受 2 个聚合
//     值：
//     - `NEW`       → 所有"运行中"母单（NEW / WAITING / PROCESSING / PAUSED
//     以及内部的 CANCEL / CANCEL_REJECT / CLEANING 中间态）
//     - `COMPLETED` → 所有"非运行中"母单（CANCELLED / COMPLETED /
//     COMPLETED_WITHTAIL / REJECTED / EXPIRED）
//     传其它细分状态在列表查询里会退化成"按字面值精确匹配"，结果通常为空，
//     **不要这样用**。推荐显式使用 `MasterOrderStatusV2New` 或
//     `MasterOrderStatusV2Completed`。
//
// 其它 V2 接口（batch-cancel / update 等）保持细分状态语义，不受影响。
type MasterOrderStatusV2 string

const (
	// MasterOrderStatusV2New 在列表查询里表示"运行中"聚合（包含 NEW / WAITING
	// / PROCESSING / PAUSED 等所有未终态状态）；在详情/推送里只代表"刚创建未
	// 启动"。
	MasterOrderStatusV2New        MasterOrderStatusV2 = "NEW"
	MasterOrderStatusV2Waiting    MasterOrderStatusV2 = "WAITING"
	MasterOrderStatusV2Processing MasterOrderStatusV2 = "PROCESSING"
	MasterOrderStatusV2Paused     MasterOrderStatusV2 = "PAUSED"
	MasterOrderStatusV2Cancelled  MasterOrderStatusV2 = "CANCELLED"
	// MasterOrderStatusV2Completed 在列表查询里表示"非运行中"聚合（包含
	// CANCELLED / COMPLETED / COMPLETED_WITHTAIL / REJECTED / EXPIRED）；
	// 在详情/推送里只代表"正常完成"。
	MasterOrderStatusV2Completed         MasterOrderStatusV2 = "COMPLETED"
	MasterOrderStatusV2CompletedWithTail MasterOrderStatusV2 = "COMPLETED_WITHTAIL"
	MasterOrderStatusV2Rejected          MasterOrderStatusV2 = "REJECTED"
	MasterOrderStatusV2Expired           MasterOrderStatusV2 = "EXPIRED"
)

// pageSizeMaxV2 is the V2 list-endpoint page size cap. Values above this
// limit are rejected by V2 APIs instead of being silently clamped.
const pageSizeMaxV2 = 100

func validatePageSizeV2(pageSize *int32) error {
	if pageSize != nil && *pageSize > pageSizeMaxV2 {
		return fmt.Errorf("pageSize %d exceeds V2 limit %d", *pageSize, pageSizeMaxV2)
	}
	return nil
}

// callAPIV2WithJSONBody sends a V2 POST/PUT request with a JSON body and
// signs the request the same way the backend's `apiAuth.CollectParamsAndBodyForSign`
// middleware verifies it: signature = HMAC-SHA256(secret, urlValues.Encode())
// where urlValues is the merge of URL query keys (timestamp, recvWindow, ...)
// and the JSON body's top-level keys, sorted by key.
//
// Use this for any V2 endpoint that requires a body. Pure GETs go through
// the existing callAPI flow.
func (c *Client) callAPIV2WithJSONBody(ctx context.Context, method, endpoint string, body params, opts ...RequestOption) ([]byte, error) {
	r := &request{secType: secTypeSigned}
	for _, opt := range opts {
		opt(r)
	}

	timestamp := currentTimestamp() - c.TimeOffset
	tsStr := strconv.FormatInt(timestamp, 10)

	// Build JSON body from non-nil params.
	var bodyBytes []byte
	if len(body) > 0 {
		// json.Marshal preserves int/float/bool/string/array/map shapes; the
		// backend uses dec.UseNumber() so numbers we emit as numbers are
		// converted via json.Number.String() for signing — exactly what we
		// reproduce below.
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	// Re-parse body bytes to mirror the backend's signing logic. This is the
	// safest way to keep the SDK and server in lock-step even when the body
	// contains nested arrays/objects (e.g. batch-cancel `masterOrderIds`).
	signValues := url.Values{}
	if len(bodyBytes) > 0 {
		var obj map[string]interface{}
		dec := json.NewDecoder(bytes.NewReader(bodyBytes))
		dec.UseNumber()
		if err := dec.Decode(&obj); err != nil {
			return nil, err
		}
		for k, v := range obj {
			if strings.EqualFold(k, "signature") || strings.EqualFold(k, "timestamp") {
				continue
			}
			if s, ok := scalarToSignString(v); ok {
				signValues.Add(k, s)
				continue
			}
			// Arrays / nested objects — backend stringifies via json.Marshal.
			if encoded, err := json.Marshal(v); err == nil {
				signValues.Add(k, string(encoded))
			}
		}
	}
	signValues.Set("timestamp", tsStr)
	if r.recvWindow > 0 {
		signValues.Set(recvWindowKey, strconv.FormatInt(r.recvWindow, 10))
	}

	signature := signWithSecret(c.SecretKey, signValues.Encode())

	// Compose URL — timestamp/recvWindow/signature go into the query string
	// alongside the JSON body, matching the backend signing middleware which
	// reads timestamp from query first.
	q := url.Values{}
	q.Set(timestampKey, tsStr)
	if r.recvWindow > 0 {
		q.Set(recvWindowKey, strconv.FormatInt(r.recvWindow, 10))
	}
	q.Set(signatureKey, signature)

	fullURL := fmt.Sprintf("%s%s?%s", c.BaseURL, endpoint, q.Encode())

	var bodyReader io.Reader
	if len(bodyBytes) > 0 {
		bodyReader = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", Name, Version))
	req.Header.Set("X-MBX-APIKEY", c.APIKey)
	if len(bodyBytes) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	c.debug("V2 request: %s %s body=%s sign=%s", method, fullURL, string(bodyBytes), signValues.Encode())

	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := res.Body.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	c.debug("V2 response status=%d body=%s", res.StatusCode, string(data))

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(handlers.APIError)
		_ = json.Unmarshal(data, apiErr)
		return nil, apiErr
	}
	respData := new(handlers.APISuccess)
	if err := json.Unmarshal(data, respData); err != nil {
		return nil, err
	}
	if respData.Code != 200 {
		return nil, &handlers.APIError{
			Code:       respData.Code,
			Reason:     respData.Reason,
			Message:    respData.Message,
			TraceId:    respData.TraceId,
			ServerTime: respData.ServerTime,
		}
	}
	return json.Marshal(respData.Message)
}

// signWithSecret HMAC-SHA256 signs the given payload using secret.
func signWithSecret(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// scalarToSignString mirrors backend `scalarToString` (gin/middleware.go) so
// the SDK and server agree on the canonical form of every scalar JSON value.
// Returns false for non-scalars (arrays, objects) — callers fall back to
// json.Marshal in that case.
func scalarToSignString(v interface{}) (string, bool) {
	switch tv := v.(type) {
	case string:
		return tv, true
	case json.Number:
		return tv.String(), true
	case float64:
		if math.IsNaN(tv) || math.IsInf(tv, 0) {
			return "", false
		}
		return strconv.FormatFloat(tv, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(tv), true
	case int:
		return strconv.FormatInt(int64(tv), 10), true
	case int32:
		return strconv.FormatInt(int64(tv), 10), true
	case int64:
		return strconv.FormatInt(tv, 10), true
	case nil:
		return fmt.Sprint(tv), true
	default:
		return "", false
	}
}

func defaultPovLimitForAlgorithmV2(algorithm trading_enums.Algorithm) string {
	if algorithm == trading_enums.AlgorithmPOV {
		return "0.05"
	}
	return "1"
}

func validatePovLimitV2(value string) error {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil || f < 0 || f > 1 {
		return errors.New("povLimit must be between 0 and 1")
	}
	return nil
}

// =============================================================================
//  /user/exchange/v2/exchange-apis (GET)
// =============================================================================

// ListExchangeApisV2Service queries V2 exchange API key bindings.
type ListExchangeApisV2Service struct {
	c        *Client
	page     *int32
	pageSize *int32
	exchange *trading_enums.Exchange
}

// Page sets the 1-based page number.
func (s *ListExchangeApisV2Service) Page(page int32) *ListExchangeApisV2Service {
	s.page = &page
	return s
}

// PageSize sets the number of items per page. Values above 100 are rejected.
func (s *ListExchangeApisV2Service) PageSize(pageSize int32) *ListExchangeApisV2Service {
	s.pageSize = &pageSize
	return s
}

// Exchange filters by exchange name.
func (s *ListExchangeApisV2Service) Exchange(exchange trading_enums.Exchange) *ListExchangeApisV2Service {
	s.exchange = &exchange
	return s
}

// Do sends the request.
func (s *ListExchangeApisV2Service) Do(ctx context.Context, opts ...RequestOption) (res *ListExchangeApisV2Reply, err error) {
	if err := validatePageSizeV2(s.pageSize); err != nil {
		return nil, err
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: v2ExchangeApisEndpoint,
		secType:  secTypeSigned,
	}
	m := params{}
	if s.page != nil {
		m["page"] = *s.page
	}
	if s.pageSize != nil {
		m["pageSize"] = *s.pageSize
	}
	if s.exchange != nil {
		m["exchange"] = *s.exchange
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(ListExchangeApisV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListExchangeApisV2Reply is the response of `GET /user/exchange/v2/exchange-apis`.
type ListExchangeApisV2Reply struct {
	Items    []ExchangeApiV2Info `json:"items"`
	Total    int32               `json:"total"`
	Page     int32               `json:"page"`
	PageSize int32               `json:"pageSize"`
}

// ExchangeApiV2Info is the per-row payload for V2 exchange API keys. V2 hides
// `verificationMethod` and `balance` compared with V1.
type ExchangeApiV2Info struct {
	ApiKeyId         string `json:"apiKeyId"`
	ApiKeyUuid       string `json:"-"` // Deprecated: use ApiKeyId.
	Id               string `json:"-"` // Deprecated: use ApiKeyId.
	CreatedAt        string `json:"createdAt"`
	AccountName      string `json:"accountName"`
	Exchange         string `json:"exchange"`
	ApiKey           string `json:"apiKey"`
	Status           string `json:"status"`
	IsValid          bool   `json:"isValid"`
	IsTradingEnabled bool   `json:"isTradingEnabled"`
	IsDefault        bool   `json:"isDefault"`
	IsPm             bool   `json:"isPm"`
}

func (i *ExchangeApiV2Info) UnmarshalJSON(data []byte) error {
	type alias ExchangeApiV2Info
	var aux struct {
		*alias
		LegacyApiKeyUuid string `json:"apiKeyUuid"`
		LegacyId         string `json:"id"`
	}
	aux.alias = (*alias)(i)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if i.ApiKeyId == "" {
		if aux.LegacyApiKeyUuid != "" {
			i.ApiKeyId = aux.LegacyApiKeyUuid
		} else {
			i.ApiKeyId = aux.LegacyId
		}
	}
	if i.ApiKeyUuid == "" {
		i.ApiKeyUuid = i.ApiKeyId
	}
	if i.Id == "" {
		if aux.LegacyId != "" {
			i.Id = aux.LegacyId
		} else {
			i.Id = i.ApiKeyId
		}
	}
	return nil
}

// =============================================================================
//  /user/trading/v2/master-orders (POST)
// =============================================================================

// CreateMasterOrderV2Service creates a V2 master order. Required fields are
// `apiKeyId`, `exchange`, `marketType`, `symbol`, `side`, `algorithm`, and
// `executionDurationSeconds` (> 10). Exactly one of `totalQuantity` or
// `orderNotional` must be set; when `isTargetPosition` is true `totalQuantity`
// is mandatory.
type CreateMasterOrderV2Service struct {
	c                        *Client
	apiKeyId                 string
	exchange                 trading_enums.Exchange
	marketType               trading_enums.MarketType
	symbol                   string
	side                     trading_enums.OrderSide
	algorithm                trading_enums.Algorithm
	executionDurationSeconds *int64
	startTimeMs              *int64
	totalQuantity            *string
	orderNotional            *string
	marginType               *trading_enums.MarginType
	reduceOnly               *bool
	isMargin                 *bool
	worstPrice               *string
	mustComplete             *bool
	makerRateLimit           *string
	povLimit                 *string
	povMinLimit              *string
	upTolerance              *string
	lowTolerance             *string
	strictUpBound            *bool
	tailOrderProtection      *bool
	enableMake               *bool
	isTargetPosition         *bool
	clientOrderId            *string
	notes                    *string
}

// ApiKeyId sets the required exchange API Key binding ID.
func (s *CreateMasterOrderV2Service) ApiKeyId(apiKeyId string) *CreateMasterOrderV2Service {
	s.apiKeyId = apiKeyId
	return s
}

// Exchange sets the required exchange.
func (s *CreateMasterOrderV2Service) Exchange(exchange trading_enums.Exchange) *CreateMasterOrderV2Service {
	s.exchange = exchange
	return s
}

// MarketType sets the required market type (`SPOT` or `PERP`).
func (s *CreateMasterOrderV2Service) MarketType(marketType trading_enums.MarketType) *CreateMasterOrderV2Service {
	s.marketType = marketType
	return s
}

// Symbol sets the required trading pair (e.g. `BTCUSDT`).
func (s *CreateMasterOrderV2Service) Symbol(symbol string) *CreateMasterOrderV2Service {
	s.symbol = symbol
	return s
}

// Side sets the required side. In normal mode it's the trade direction; in
// target-position mode it's the position side.
func (s *CreateMasterOrderV2Service) Side(side trading_enums.OrderSide) *CreateMasterOrderV2Service {
	s.side = side
	return s
}

// Algorithm sets the required trading algorithm (`TWAP`, `VWAP`, `POV`).
func (s *CreateMasterOrderV2Service) Algorithm(algorithm trading_enums.Algorithm) *CreateMasterOrderV2Service {
	s.algorithm = algorithm
	return s
}

// ExecutionDurationSeconds sets the maximum execution duration in seconds.
// Must be > 10. V2 always uses seconds (V1 mixed minutes and seconds).
func (s *CreateMasterOrderV2Service) ExecutionDurationSeconds(seconds int64) *CreateMasterOrderV2Service {
	s.executionDurationSeconds = &seconds
	return s
}

// StartTimeMs sets the execution start time in epoch milliseconds. Omit for
// immediate execution.
func (s *CreateMasterOrderV2Service) StartTimeMs(ms int64) *CreateMasterOrderV2Service {
	s.startTimeMs = &ms
	return s
}

// TotalQuantity sets the trade quantity as a decimal string (V2 sends Decimal
// values as strings to avoid JS float precision issues).
func (s *CreateMasterOrderV2Service) TotalQuantity(qty string) *CreateMasterOrderV2Service {
	s.totalQuantity = &qty
	return s
}

// OrderNotional sets the trade notional as a decimal string.
func (s *CreateMasterOrderV2Service) OrderNotional(notional string) *CreateMasterOrderV2Service {
	s.orderNotional = &notional
	return s
}

// MarginType sets the contract margin type (`U` / `C`); required for `PERP`.
func (s *CreateMasterOrderV2Service) MarginType(mt trading_enums.MarginType) *CreateMasterOrderV2Service {
	s.marginType = &mt
	return s
}

// ReduceOnly sets reduce-only mode (PERP only).
func (s *CreateMasterOrderV2Service) ReduceOnly(reduceOnly bool) *CreateMasterOrderV2Service {
	s.reduceOnly = &reduceOnly
	return s
}

// IsMargin enables spot-margin mode (SPOT only).
func (s *CreateMasterOrderV2Service) IsMargin(isMargin bool) *CreateMasterOrderV2Service {
	s.isMargin = &isMargin
	return s
}

// WorstPrice sets the worst acceptable price as a decimal string.
func (s *CreateMasterOrderV2Service) WorstPrice(price string) *CreateMasterOrderV2Service {
	s.worstPrice = &price
	return s
}

// MustComplete toggles whether the order must finish within executionDurationSeconds.
func (s *CreateMasterOrderV2Service) MustComplete(mustComplete bool) *CreateMasterOrderV2Service {
	s.mustComplete = &mustComplete
	return s
}

// MakerRateLimit sets the minimum maker fill ratio as a decimal string (0-1; -1 for auto).
func (s *CreateMasterOrderV2Service) MakerRateLimit(rate string) *CreateMasterOrderV2Service {
	s.makerRateLimit = &rate
	return s
}

// PovLimit sets the participation rate cap as a decimal string.
func (s *CreateMasterOrderV2Service) PovLimit(rate string) *CreateMasterOrderV2Service {
	s.povLimit = &rate
	return s
}

// PovMinLimit sets the participation rate floor as a decimal string (POV only).
func (s *CreateMasterOrderV2Service) PovMinLimit(rate string) *CreateMasterOrderV2Service {
	s.povMinLimit = &rate
	return s
}

// UpTolerance sets the upper progress tolerance as a decimal string.
func (s *CreateMasterOrderV2Service) UpTolerance(tol string) *CreateMasterOrderV2Service {
	s.upTolerance = &tol
	return s
}

// LowTolerance sets the lower progress tolerance as a decimal string.
func (s *CreateMasterOrderV2Service) LowTolerance(tol string) *CreateMasterOrderV2Service {
	s.lowTolerance = &tol
	return s
}

// StrictUpBound enables strict upper-bound enforcement.
func (s *CreateMasterOrderV2Service) StrictUpBound(strict bool) *CreateMasterOrderV2Service {
	s.strictUpBound = &strict
	return s
}

// TailOrderProtection toggles tail-order protection (default true).
func (s *CreateMasterOrderV2Service) TailOrderProtection(enabled bool) *CreateMasterOrderV2Service {
	s.tailOrderProtection = &enabled
	return s
}

// EnableMake toggles maker orders (default true; false → all taker).
func (s *CreateMasterOrderV2Service) EnableMake(enabled bool) *CreateMasterOrderV2Service {
	s.enableMake = &enabled
	return s
}

// IsTargetPosition switches into target-position mode (forces TotalQuantity).
func (s *CreateMasterOrderV2Service) IsTargetPosition(isTarget bool) *CreateMasterOrderV2Service {
	s.isTargetPosition = &isTarget
	return s
}

// ClientOrderId sets the user-defined order ID (for idempotency / tracking).
func (s *CreateMasterOrderV2Service) ClientOrderId(clientOrderId string) *CreateMasterOrderV2Service {
	s.clientOrderId = &clientOrderId
	return s
}

// Notes sets a free-form note attached to the order.
func (s *CreateMasterOrderV2Service) Notes(notes string) *CreateMasterOrderV2Service {
	s.notes = &notes
	return s
}

// Do sends the request.
func (s *CreateMasterOrderV2Service) Do(ctx context.Context, opts ...RequestOption) (res *CreateMasterOrderV2Reply, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	m := params{
		"apiKeyId":   s.apiKeyId,
		"exchange":   string(s.exchange),
		"marketType": string(s.marketType),
		"symbol":     s.symbol,
		"side":       string(s.side),
		"algorithm":  string(s.algorithm),
	}
	if s.executionDurationSeconds != nil {
		m["executionDurationSeconds"] = *s.executionDurationSeconds
	}
	if s.startTimeMs != nil {
		m["startTimeMs"] = *s.startTimeMs
	}
	if s.totalQuantity != nil {
		m["totalQuantity"] = *s.totalQuantity
	}
	if s.orderNotional != nil {
		m["orderNotional"] = *s.orderNotional
	}
	if s.marginType != nil {
		m["marginType"] = string(*s.marginType)
	}
	if s.reduceOnly != nil {
		m["reduceOnly"] = *s.reduceOnly
	}
	if s.isMargin != nil {
		m["isMargin"] = *s.isMargin
	}
	if s.worstPrice != nil {
		m["worstPrice"] = *s.worstPrice
	}
	if s.mustComplete != nil {
		m["mustComplete"] = *s.mustComplete
	}
	if s.makerRateLimit != nil {
		m["makerRateLimit"] = *s.makerRateLimit
	}
	if s.povLimit == nil {
		defaultPovLimit := defaultPovLimitForAlgorithmV2(s.algorithm)
		s.povLimit = &defaultPovLimit
	}
	if s.povLimit != nil {
		m["povLimit"] = *s.povLimit
	}
	if s.povMinLimit != nil {
		m["povMinLimit"] = *s.povMinLimit
	}
	if s.upTolerance != nil {
		m["upTolerance"] = *s.upTolerance
	}
	if s.lowTolerance != nil {
		m["lowTolerance"] = *s.lowTolerance
	}
	if s.strictUpBound != nil {
		m["strictUpBound"] = *s.strictUpBound
	}
	if s.tailOrderProtection != nil {
		m["tailOrderProtection"] = *s.tailOrderProtection
	}
	if s.enableMake != nil {
		m["enableMake"] = *s.enableMake
	}
	if s.isTargetPosition != nil {
		m["isTargetPosition"] = *s.isTargetPosition
	}
	if s.clientOrderId != nil {
		m["clientOrderId"] = *s.clientOrderId
	}
	if s.notes != nil {
		m["notes"] = *s.notes
	}

	data, err := s.c.callAPIV2WithJSONBody(ctx, http.MethodPost, v2MasterOrdersEndpoint, m, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateMasterOrderV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *CreateMasterOrderV2Service) validate() error {
	if s.apiKeyId == "" {
		return errors.New("apiKeyId is required")
	}
	if s.exchange == "" {
		return errors.New("exchange is required")
	}
	if s.marketType == "" {
		return errors.New("marketType is required")
	}
	if s.symbol == "" {
		return errors.New("symbol is required")
	}
	if s.side == "" {
		return errors.New("side is required")
	}
	if s.algorithm == "" {
		return errors.New("algorithm is required")
	}
	if s.executionDurationSeconds == nil {
		return errors.New("executionDurationSeconds is required (must be > 10)")
	}
	if *s.executionDurationSeconds <= 10 {
		return errors.New("executionDurationSeconds must be greater than 10")
	}
	if s.povLimit != nil {
		if err := validatePovLimitV2(*s.povLimit); err != nil {
			return err
		}
	}
	hasQty := s.totalQuantity != nil
	hasNotional := s.orderNotional != nil
	if hasQty == hasNotional {
		return errors.New("exactly one of totalQuantity / orderNotional must be provided")
	}
	if s.isTargetPosition != nil && *s.isTargetPosition {
		if !hasQty {
			return errors.New("totalQuantity is required when isTargetPosition is true")
		}
		if hasNotional {
			return errors.New("orderNotional is not allowed when isTargetPosition is true")
		}
	}
	// Same Deribit / Binance perp_cm guards as the V1 service.
	if s.exchange == trading_enums.ExchangeDeribit &&
		(strings.EqualFold(s.symbol, "BTCUSD") || strings.EqualFold(s.symbol, "ETHUSD")) {
		if hasNotional {
			return errors.New("orderNotional is not allowed when exchange is Deribit and symbol is BTCUSD or ETHUSD; use totalQuantity (unit: USD) instead")
		}
		if !hasQty {
			return errors.New("totalQuantity is required when exchange is Deribit and symbol is BTCUSD or ETHUSD (unit: USD)")
		}
	}
	if s.exchange == trading_enums.ExchangeBinance &&
		s.marketType == trading_enums.MarketTypePerp &&
		s.marginType != nil && *s.marginType == trading_enums.MarginTypeC {
		if hasNotional {
			return errors.New("orderNotional is not allowed when exchange is Binance and marginType is C for PERP orders; use totalQuantity (unit: contracts) instead")
		}
		if !hasQty {
			return errors.New("totalQuantity is required when exchange is Binance and marginType is C for PERP orders (unit: contracts)")
		}
	}
	return nil
}

// CreateMasterOrderV2Reply is the response of `POST /user/trading/v2/master-orders`.
type CreateMasterOrderV2Reply struct {
	MasterOrderId string `json:"masterOrderId"`
	Status        string `json:"status"`
	ClientOrderId string `json:"clientOrderId"`
}

// =============================================================================
//  /user/trading/v2/master-orders (GET) — list & detail
// =============================================================================

// GetMasterOrdersV2Service queries the V2 master order list.
type GetMasterOrdersV2Service struct {
	c             *Client
	page          *int32
	pageSize      *int32
	status        *MasterOrderStatusV2
	exchange      *string
	symbol        *string
	algorithm     *trading_enums.Algorithm
	apiKeyId      *string
	startTime     *string
	endTime       *string
	masterOrderId *string
}

// Page sets the 1-based page number.
func (s *GetMasterOrdersV2Service) Page(page int32) *GetMasterOrdersV2Service {
	s.page = &page
	return s
}

// PageSize sets the per-page count. Values above 100 are rejected.
func (s *GetMasterOrdersV2Service) PageSize(pageSize int32) *GetMasterOrdersV2Service {
	s.pageSize = &pageSize
	return s
}

// Status filters by master order status.
//
// 注意：列表查询的过滤只接受聚合值 MasterOrderStatusV2New（=运行中所有状态）
// 或 MasterOrderStatusV2Completed（=非运行中所有状态）。传其它细分状态会被
// 后端按字面值匹配，结果通常为空——这一点与详情/推送里返回的 9 种细分状态
// 不同（详见 `MasterOrderStatusV2` 类型注释）。
func (s *GetMasterOrdersV2Service) Status(status MasterOrderStatusV2) *GetMasterOrdersV2Service {
	s.status = &status
	return s
}

// Exchange filters by exchange.
func (s *GetMasterOrdersV2Service) Exchange(exchange string) *GetMasterOrdersV2Service {
	s.exchange = &exchange
	return s
}

// Symbol filters by trading pair.
func (s *GetMasterOrdersV2Service) Symbol(symbol string) *GetMasterOrdersV2Service {
	s.symbol = &symbol
	return s
}

// Algorithm filters by algorithm.
func (s *GetMasterOrdersV2Service) Algorithm(algorithm trading_enums.Algorithm) *GetMasterOrdersV2Service {
	s.algorithm = &algorithm
	return s
}

// ApiKeyId filters by exchange API key binding ID.
func (s *GetMasterOrdersV2Service) ApiKeyId(apiKeyId string) *GetMasterOrdersV2Service {
	s.apiKeyId = &apiKeyId
	return s
}

// ApiKeyUuid is kept for source compatibility. New code should use ApiKeyId.
func (s *GetMasterOrdersV2Service) ApiKeyUuid(uuid string) *GetMasterOrdersV2Service {
	return s.ApiKeyId(uuid)
}

// StartTime sets the lower bound (RFC3339 / ISO 8601).
func (s *GetMasterOrdersV2Service) StartTime(startTime string) *GetMasterOrdersV2Service {
	s.startTime = &startTime
	return s
}

// EndTime sets the upper bound (RFC3339 / ISO 8601).
func (s *GetMasterOrdersV2Service) EndTime(endTime string) *GetMasterOrdersV2Service {
	s.endTime = &endTime
	return s
}

// MasterOrderId filters by exact master order ID.
func (s *GetMasterOrdersV2Service) MasterOrderId(id string) *GetMasterOrdersV2Service {
	s.masterOrderId = &id
	return s
}

// Do sends the request.
func (s *GetMasterOrdersV2Service) Do(ctx context.Context, opts ...RequestOption) (res *GetMasterOrdersV2Reply, err error) {
	if err := validatePageSizeV2(s.pageSize); err != nil {
		return nil, err
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: v2MasterOrdersEndpoint,
		secType:  secTypeSigned,
	}
	m := params{}
	if s.page != nil {
		m["page"] = *s.page
	}
	if s.pageSize != nil {
		m["pageSize"] = *s.pageSize
	}
	if s.status != nil {
		m["status"] = string(*s.status)
	}
	if s.exchange != nil {
		m["exchange"] = *s.exchange
	}
	if s.symbol != nil {
		m["symbol"] = *s.symbol
	}
	if s.algorithm != nil {
		m["algorithm"] = string(*s.algorithm)
	}
	if s.apiKeyId != nil {
		m["apiKeyId"] = *s.apiKeyId
	}
	if s.startTime != nil {
		m["startTime"] = *s.startTime
	}
	if s.endTime != nil {
		m["endTime"] = *s.endTime
	}
	if s.masterOrderId != nil {
		m["masterOrderId"] = *s.masterOrderId
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(GetMasterOrdersV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetMasterOrdersV2Reply is the response of `GET /user/trading/v2/master-orders`.
type GetMasterOrdersV2Reply struct {
	Items    []MasterOrderV2Info `json:"items"`
	Total    int32               `json:"total"`
	Page     int32               `json:"page"`
	PageSize int32               `json:"pageSize"`
}

// MasterOrderV2Info is the V2 master order DTO. Fields hidden by V2
// (`apiKey`, `apiKeyName`, `ticktimeInt`, `ticktimeMs`, `submitTimeMs`,
// `algoStartTimeMs`, ...) are intentionally absent.
type MasterOrderV2Info struct {
	CreatedAt                string            `json:"createdAt"`
	UpdatedAt                string            `json:"updatedAt"`
	MasterOrderId            string            `json:"masterOrderId"`
	ClientOrderId            string            `json:"clientOrderId"`
	ApiKeyId                 string            `json:"apiKeyId"`
	ApiKeyUuid               string            `json:"-"` // Deprecated: use ApiKeyId.
	TradingAccount           string            `json:"tradingAccount"`
	Exchange                 string            `json:"exchange"`
	MarketType               string            `json:"marketType"`
	Category                 string            `json:"category"`
	Symbol                   string            `json:"symbol"`
	BaseCurrency             string            `json:"baseCurrency"`
	QuoteCurrency            string            `json:"quoteCurrency"`
	Side                     string            `json:"side"`
	MarginType               *string           `json:"marginType,omitempty"`
	ReduceOnly               *bool             `json:"reduceOnly,omitempty"`
	IsMargin                 *bool             `json:"isMargin,omitempty"`
	Algorithm                string            `json:"algorithm"`
	TotalQuantity            *string           `json:"totalQuantity,omitempty"`
	OrderNotional            *string           `json:"orderNotional,omitempty"`
	StartTimeMs              *FlexInt64        `json:"startTimeMs,omitempty"`
	ExecutionDurationSeconds *FlexInt64        `json:"executionDurationSeconds,omitempty"`
	WorstPrice               *string           `json:"worstPrice,omitempty"`
	MustComplete             *bool             `json:"mustComplete,omitempty"`
	MakerRateLimit           *string           `json:"makerRateLimit,omitempty"`
	PovLimit                 *string           `json:"povLimit,omitempty"`
	PovMinLimit              *string           `json:"povMinLimit,omitempty"`
	UpTolerance              *string           `json:"upTolerance,omitempty"`
	LowTolerance             *string           `json:"lowTolerance,omitempty"`
	StrictUpBound            *bool             `json:"strictUpBound,omitempty"`
	TailOrderProtection      *bool             `json:"tailOrderProtection,omitempty"`
	EnableMake               *bool             `json:"enableMake,omitempty"`
	IsTargetPosition         *bool             `json:"isTargetPosition,omitempty"`
	Notes                    string            `json:"notes"`
	Status                   string            `json:"status"`
	RejectReason             string            `json:"rejectReason"`
	FinishedMs               *FlexInt64        `json:"finishedMs,omitempty"`
	CumFilledQty             *string           `json:"cumFilledQty,omitempty"`
	CumFilledNotional        *string           `json:"cumFilledNotional,omitempty"`
	AvgFilledPrice           *string           `json:"avgFilledPrice,omitempty"`
	MakerRate                *string           `json:"makerRate,omitempty"`
	CompletedQuantity        *string           `json:"completedQuantity,omitempty"`
	Commission               map[string]string `json:"commission"`
}

func (i *MasterOrderV2Info) UnmarshalJSON(data []byte) error {
	type alias MasterOrderV2Info
	var aux struct {
		*alias
		LegacyApiKeyUuid string `json:"apiKeyUuid"`
	}
	aux.alias = (*alias)(i)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if i.ApiKeyId == "" {
		i.ApiKeyId = aux.LegacyApiKeyUuid
	}
	if i.ApiKeyUuid == "" {
		i.ApiKeyUuid = i.ApiKeyId
	}
	return nil
}

// GetMasterOrderDetailV2Service fetches a master order detail by `masterOrderId`.
type GetMasterOrderDetailV2Service struct {
	c             *Client
	masterOrderId string
}

// MasterOrderId sets the path parameter.
func (s *GetMasterOrderDetailV2Service) MasterOrderId(id string) *GetMasterOrderDetailV2Service {
	s.masterOrderId = id
	return s
}

// Do sends the request.
func (s *GetMasterOrderDetailV2Service) Do(ctx context.Context, opts ...RequestOption) (res *GetMasterOrderDetailV2Reply, err error) {
	if s.masterOrderId == "" {
		return nil, errors.New("masterOrderId is required")
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("%s/%s", v2MasterOrdersEndpoint, s.masterOrderId),
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(GetMasterOrderDetailV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetMasterOrderDetailV2Reply is the response wrapper for V2 master order detail.
type GetMasterOrderDetailV2Reply struct {
	MasterOrder MasterOrderV2Info `json:"masterOrder"`
}

// GetMasterOrderDetailByClientOrderIdV2Service fetches a master order by
// `clientOrderId` for idempotent lookups.
type GetMasterOrderDetailByClientOrderIdV2Service struct {
	c             *Client
	clientOrderId string
}

// ClientOrderId sets the path parameter.
func (s *GetMasterOrderDetailByClientOrderIdV2Service) ClientOrderId(id string) *GetMasterOrderDetailByClientOrderIdV2Service {
	s.clientOrderId = id
	return s
}

// Do sends the request.
func (s *GetMasterOrderDetailByClientOrderIdV2Service) Do(ctx context.Context, opts ...RequestOption) (res *GetMasterOrderDetailV2Reply, err error) {
	if s.clientOrderId == "" {
		return nil, errors.New("clientOrderId is required")
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("%s/%s", v2MasterOrdersByClientId, s.clientOrderId),
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(GetMasterOrderDetailV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// =============================================================================
//  /user/trading/v2/order-fills (GET)
// =============================================================================

// GetOrderFillsV2Service queries the V2 order fills (sub-orders) list.
type GetOrderFillsV2Service struct {
	c             *Client
	page          *int32
	pageSize      *int32
	masterOrderId *string
	orderId       *string
	clientOrderId *string
	symbol        *string
	status        *string
	startTime     *string
	endTime       *string
}

// Page sets the 1-based page number.
func (s *GetOrderFillsV2Service) Page(page int32) *GetOrderFillsV2Service {
	s.page = &page
	return s
}

// PageSize sets the per-page count. Values above 100 are rejected.
func (s *GetOrderFillsV2Service) PageSize(pageSize int32) *GetOrderFillsV2Service {
	s.pageSize = &pageSize
	return s
}

// MasterOrderId filters by parent master order.
func (s *GetOrderFillsV2Service) MasterOrderId(id string) *GetOrderFillsV2Service {
	s.masterOrderId = &id
	return s
}

// OrderId filters by exchange order ID. V2 replaces V1's `subOrderId`.
func (s *GetOrderFillsV2Service) OrderId(id string) *GetOrderFillsV2Service {
	s.orderId = &id
	return s
}

// ClientOrderId filters by user-defined client order ID.
func (s *GetOrderFillsV2Service) ClientOrderId(id string) *GetOrderFillsV2Service {
	s.clientOrderId = &id
	return s
}

// Symbol filters by trading pair.
func (s *GetOrderFillsV2Service) Symbol(symbol string) *GetOrderFillsV2Service {
	s.symbol = &symbol
	return s
}

// Status filters by sub-order status.
func (s *GetOrderFillsV2Service) Status(status string) *GetOrderFillsV2Service {
	s.status = &status
	return s
}

// StartTime sets the lower bound (RFC3339 / ISO 8601).
func (s *GetOrderFillsV2Service) StartTime(startTime string) *GetOrderFillsV2Service {
	s.startTime = &startTime
	return s
}

// EndTime sets the upper bound (RFC3339 / ISO 8601).
func (s *GetOrderFillsV2Service) EndTime(endTime string) *GetOrderFillsV2Service {
	s.endTime = &endTime
	return s
}

// Do sends the request.
func (s *GetOrderFillsV2Service) Do(ctx context.Context, opts ...RequestOption) (res *GetOrderFillsV2Reply, err error) {
	if err := validatePageSizeV2(s.pageSize); err != nil {
		return nil, err
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: v2OrderFillsEndpoint,
		secType:  secTypeSigned,
	}
	m := params{}
	if s.page != nil {
		m["page"] = *s.page
	}
	if s.pageSize != nil {
		m["pageSize"] = *s.pageSize
	}
	if s.masterOrderId != nil {
		m["masterOrderId"] = *s.masterOrderId
	}
	if s.orderId != nil {
		m["orderId"] = *s.orderId
	}
	if s.clientOrderId != nil {
		m["clientOrderId"] = *s.clientOrderId
	}
	if s.symbol != nil {
		m["symbol"] = *s.symbol
	}
	if s.status != nil {
		m["status"] = *s.status
	}
	if s.startTime != nil {
		m["startTime"] = *s.startTime
	}
	if s.endTime != nil {
		m["endTime"] = *s.endTime
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(GetOrderFillsV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetOrderFillsV2Reply is the response of `GET /user/trading/v2/order-fills`.
type GetOrderFillsV2Reply struct {
	Items    []OrderFillV2Info `json:"items"`
	Total    int32             `json:"total"`
	Page     int32             `json:"page"`
	PageSize int32             `json:"pageSize"`
}

// OrderFillV2Info is the V2 sub-order DTO. V2 hides `fee` / `tradingAccount`,
// renames `filledValue` → `filledNotional` and `subOrderId` → `orderId`,
// `type` → `orderType`.
//
// 关于 decimal-shaped 字段：V2 后端契约是 *string*（避免 JS 精度丢失），但为了
// 兼容历史/异常返回 number 的情况，这里统一使用 FlexDecimalString —— 它会接受
// JSON string 或 JSON number，对外统一成 string。
type OrderFillV2Info struct {
	Id               string            `json:"id"`
	OrderCreatedTime string            `json:"orderCreatedTime"`
	MasterOrderId    string            `json:"masterOrderId"`
	Exchange         string            `json:"exchange"`
	Category         string            `json:"category"`
	Symbol           string            `json:"symbol"`
	Side             string            `json:"side"`
	FilledNotional   FlexDecimalString `json:"filledNotional"`
	FilledQuantity   FlexDecimalString `json:"filledQuantity"`
	AveragePrice     FlexDecimalString `json:"averagePrice"`
	Price            FlexDecimalString `json:"price"`
	Status           string            `json:"status"`
	RejectReason     string            `json:"rejectReason"`
	BaseCurrency     string            `json:"baseCurrency"`
	QuoteCurrency    string            `json:"quoteCurrency"`
	OrderType        string            `json:"orderType"`
	OrderId          string            `json:"orderId"`
	Quantity         FlexDecimalString `json:"quantity"`
	CreatedAt        string            `json:"createdAt"`
	UpdatedAt        string            `json:"updatedAt"`
}

// =============================================================================
//  /user/trading/v2/tca-analysis (GET)
// =============================================================================

// GetTCAAnalysisV2Service queries post-trade TCA analysis results.
type GetTCAAnalysisV2Service struct {
	c         *Client
	symbol    *string
	category  *string
	strategy  *string
	apiKeyId  *string
	startTime *int64
	endTime   *int64
}

// Symbol filters by trading pair.
func (s *GetTCAAnalysisV2Service) Symbol(symbol string) *GetTCAAnalysisV2Service {
	s.symbol = &symbol
	return s
}

// Category filters by trading category: spot, perp, or perp_cm.
func (s *GetTCAAnalysisV2Service) Category(category string) *GetTCAAnalysisV2Service {
	s.category = &category
	return s
}

// Strategy filters by execution algorithm: TWAP, VWAP, or POV.
func (s *GetTCAAnalysisV2Service) Strategy(strategy string) *GetTCAAnalysisV2Service {
	s.strategy = &strategy
	return s
}

// ApiKeyId filters by exchange API key binding ID. Comma-separated IDs are
// supported by the backend.
func (s *GetTCAAnalysisV2Service) ApiKeyId(apiKeyId string) *GetTCAAnalysisV2Service {
	s.apiKeyId = &apiKeyId
	return s
}

// ApiKeyUuid is kept for source compatibility. New code should use ApiKeyId.
func (s *GetTCAAnalysisV2Service) ApiKeyUuid(uuid string) *GetTCAAnalysisV2Service {
	return s.ApiKeyId(uuid)
}

// StartTime sets the lower bound in epoch milliseconds.
func (s *GetTCAAnalysisV2Service) StartTime(startTime int64) *GetTCAAnalysisV2Service {
	s.startTime = &startTime
	return s
}

// EndTime sets the upper bound in epoch milliseconds.
func (s *GetTCAAnalysisV2Service) EndTime(endTime int64) *GetTCAAnalysisV2Service {
	s.endTime = &endTime
	return s
}

// Do sends the request. Successful strategy-api responses return the TCA
// result array directly from `message`.
func (s *GetTCAAnalysisV2Service) Do(ctx context.Context, opts ...RequestOption) (res []*TCAAnalysisV2Info, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: v2TCAAnalysisEndpoint,
		secType:  secTypeSigned,
	}
	m := params{}
	if s.symbol != nil {
		m["symbol"] = *s.symbol
	}
	if s.category != nil {
		m["category"] = *s.category
	}
	if s.strategy != nil {
		m["strategy"] = *s.strategy
	}
	if s.apiKeyId != nil {
		m["apiKeyId"] = *s.apiKeyId
	}
	if s.startTime != nil {
		m["startTime"] = *s.startTime
	}
	if s.endTime != nil {
		m["endTime"] = *s.endTime
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*TCAAnalysisV2Info, 0)
	if len(data) == 0 {
		return res, nil
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// TCAAnalysisV2Info is a single V2 TCA analysis row.
type TCAAnalysisV2Info struct {
	MasterOrderId           string  `json:"masterOrderId"`
	StartTime               string  `json:"startTime"`
	EndTime                 string  `json:"endTime"`
	FinishedTime            string  `json:"finishedTime"`
	Strategy                string  `json:"strategy"`
	Category                string  `json:"category"`
	OrderQuantity           float64 `json:"orderQuantity"`
	OrderNotional           float64 `json:"orderNotional"`
	ArrivalPrice            float64 `json:"arrivalPrice"`
	ExecutionRate           float64 `json:"executionRate"`
	FilledQuantity          float64 `json:"filledQuantity"`
	TakerFilledNotional     float64 `json:"takerFilledNotional"`
	MakerFilledNotional     float64 `json:"makerFilledNotional"`
	FilledNotional          float64 `json:"filledNotional"`
	MakerRate               float64 `json:"makerRate"`
	ChildOrderCount         int32   `json:"childOrderCount"`
	AverageFillPrice        float64 `json:"averageFillPrice"`
	Slippage                float64 `json:"Slippage"`
	SlippagePct             float64 `json:"Slippage_pct"`
	TwapSlippagePct         float64 `json:"TWAP_Slippage_pct"`
	VwapSlippagePct         float64 `json:"VWAP_Slippage_pct"`
	Spread                  float64 `json:"Spread"`
	SlippagePctFartouch     float64 `json:"Slippage_pct_Fartouch"`
	TwapSlippagePctFartouch float64 `json:"TWAP_Slippage_pct_Fartouch"`
	VwapSlippagePctFartouch float64 `json:"VWAP_Slippage_pct_Fartouch"`
	IntervalReturn          float64 `json:"IntervalReturn"`
	ParticipationRate       float64 `json:"ParticipationRate"`
	FeeSavingPct            float64 `json:"FeeSaving_pct"`
	Date                    string  `json:"Date"`
}

// =============================================================================
//  Master order action endpoints (cancel / pause / resume / update / batch)
// =============================================================================

// MasterOrderActionV2Reply is the standard reply for cancel/pause/resume/update.
type MasterOrderActionV2Reply struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CancelMasterOrderV2Service cancels a V2 master order.
type CancelMasterOrderV2Service struct {
	c             *Client
	masterOrderId string
	reason        *string
}

// MasterOrderId sets the path parameter.
func (s *CancelMasterOrderV2Service) MasterOrderId(id string) *CancelMasterOrderV2Service {
	s.masterOrderId = id
	return s
}

// Reason sets the optional cancellation reason.
func (s *CancelMasterOrderV2Service) Reason(reason string) *CancelMasterOrderV2Service {
	s.reason = &reason
	return s
}

// Do sends the request.
func (s *CancelMasterOrderV2Service) Do(ctx context.Context, opts ...RequestOption) (res *MasterOrderActionV2Reply, err error) {
	if s.masterOrderId == "" {
		return nil, errors.New("masterOrderId is required")
	}
	endpoint := fmt.Sprintf("%s/%s/cancel", v2MasterOrdersEndpoint, s.masterOrderId)
	body := params{}
	if s.reason != nil {
		body["reason"] = *s.reason
	}
	data, err := s.c.callAPIV2WithJSONBody(ctx, http.MethodPut, endpoint, body, opts...)
	if err != nil {
		return nil, err
	}
	res = new(MasterOrderActionV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// PauseMasterOrderV2Service pauses a running V2 master order.
type PauseMasterOrderV2Service struct {
	c             *Client
	masterOrderId string
	reason        *string
}

// MasterOrderId sets the path parameter.
func (s *PauseMasterOrderV2Service) MasterOrderId(id string) *PauseMasterOrderV2Service {
	s.masterOrderId = id
	return s
}

// Reason sets the optional pause reason.
func (s *PauseMasterOrderV2Service) Reason(reason string) *PauseMasterOrderV2Service {
	s.reason = &reason
	return s
}

// Do sends the request.
func (s *PauseMasterOrderV2Service) Do(ctx context.Context, opts ...RequestOption) (res *MasterOrderActionV2Reply, err error) {
	if s.masterOrderId == "" {
		return nil, errors.New("masterOrderId is required")
	}
	endpoint := fmt.Sprintf("%s/%s/pause", v2MasterOrdersEndpoint, s.masterOrderId)
	body := params{}
	if s.reason != nil {
		body["reason"] = *s.reason
	}
	data, err := s.c.callAPIV2WithJSONBody(ctx, http.MethodPut, endpoint, body, opts...)
	if err != nil {
		return nil, err
	}
	res = new(MasterOrderActionV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ResumeMasterOrderV2Service resumes a paused V2 master order.
type ResumeMasterOrderV2Service struct {
	c             *Client
	masterOrderId string
	reason        *string
}

// MasterOrderId sets the path parameter.
func (s *ResumeMasterOrderV2Service) MasterOrderId(id string) *ResumeMasterOrderV2Service {
	s.masterOrderId = id
	return s
}

// Reason sets the optional resume reason.
func (s *ResumeMasterOrderV2Service) Reason(reason string) *ResumeMasterOrderV2Service {
	s.reason = &reason
	return s
}

// Do sends the request.
func (s *ResumeMasterOrderV2Service) Do(ctx context.Context, opts ...RequestOption) (res *MasterOrderActionV2Reply, err error) {
	if s.masterOrderId == "" {
		return nil, errors.New("masterOrderId is required")
	}
	endpoint := fmt.Sprintf("%s/%s/resume", v2MasterOrdersEndpoint, s.masterOrderId)
	body := params{}
	if s.reason != nil {
		body["reason"] = *s.reason
	}
	data, err := s.c.callAPIV2WithJSONBody(ctx, http.MethodPut, endpoint, body, opts...)
	if err != nil {
		return nil, err
	}
	res = new(MasterOrderActionV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateMasterOrderParamsV2Service updates parameters of a running V2 master order.
type UpdateMasterOrderParamsV2Service struct {
	c                        *Client
	masterOrderId            string
	totalQuantity            *string
	orderNotional            *string
	upTolerance              *string
	lowTolerance             *string
	enableMake               *bool
	makerRateLimit           *string
	strictUpBound            *bool
	povLimit                 *string
	povMinLimit              *string
	worstPrice               *string
	tailOrderProtection      *bool
	mustComplete             *bool
	executionDurationSeconds *int64
}

// MasterOrderId sets the required path parameter.
func (s *UpdateMasterOrderParamsV2Service) MasterOrderId(id string) *UpdateMasterOrderParamsV2Service {
	s.masterOrderId = id
	return s
}

// TotalQuantity updates the total trade quantity (decimal string).
func (s *UpdateMasterOrderParamsV2Service) TotalQuantity(qty string) *UpdateMasterOrderParamsV2Service {
	s.totalQuantity = &qty
	return s
}

// OrderNotional updates the order notional (decimal string).
func (s *UpdateMasterOrderParamsV2Service) OrderNotional(n string) *UpdateMasterOrderParamsV2Service {
	s.orderNotional = &n
	return s
}

// UpTolerance updates the upper progress tolerance (decimal string).
func (s *UpdateMasterOrderParamsV2Service) UpTolerance(tol string) *UpdateMasterOrderParamsV2Service {
	s.upTolerance = &tol
	return s
}

// LowTolerance updates the lower progress tolerance (decimal string).
func (s *UpdateMasterOrderParamsV2Service) LowTolerance(tol string) *UpdateMasterOrderParamsV2Service {
	s.lowTolerance = &tol
	return s
}

// EnableMake toggles maker orders.
func (s *UpdateMasterOrderParamsV2Service) EnableMake(enabled bool) *UpdateMasterOrderParamsV2Service {
	s.enableMake = &enabled
	return s
}

// MakerRateLimit updates the minimum maker fill ratio (decimal string).
func (s *UpdateMasterOrderParamsV2Service) MakerRateLimit(rate string) *UpdateMasterOrderParamsV2Service {
	s.makerRateLimit = &rate
	return s
}

// StrictUpBound toggles strict upper-bound enforcement.
func (s *UpdateMasterOrderParamsV2Service) StrictUpBound(strict bool) *UpdateMasterOrderParamsV2Service {
	s.strictUpBound = &strict
	return s
}

// PovLimit updates the participation rate cap (decimal string).
func (s *UpdateMasterOrderParamsV2Service) PovLimit(rate string) *UpdateMasterOrderParamsV2Service {
	s.povLimit = &rate
	return s
}

// PovMinLimit updates the participation rate floor (decimal string).
func (s *UpdateMasterOrderParamsV2Service) PovMinLimit(rate string) *UpdateMasterOrderParamsV2Service {
	s.povMinLimit = &rate
	return s
}

// WorstPrice updates the worst acceptable price (decimal string).
func (s *UpdateMasterOrderParamsV2Service) WorstPrice(price string) *UpdateMasterOrderParamsV2Service {
	s.worstPrice = &price
	return s
}

// TailOrderProtection toggles tail-order protection.
func (s *UpdateMasterOrderParamsV2Service) TailOrderProtection(enabled bool) *UpdateMasterOrderParamsV2Service {
	s.tailOrderProtection = &enabled
	return s
}

// MustComplete toggles must-complete behaviour.
func (s *UpdateMasterOrderParamsV2Service) MustComplete(must bool) *UpdateMasterOrderParamsV2Service {
	s.mustComplete = &must
	return s
}

// ExecutionDurationSeconds updates the execution duration in seconds (must be > 10).
func (s *UpdateMasterOrderParamsV2Service) ExecutionDurationSeconds(seconds int64) *UpdateMasterOrderParamsV2Service {
	s.executionDurationSeconds = &seconds
	return s
}

// Do sends the request.
func (s *UpdateMasterOrderParamsV2Service) Do(ctx context.Context, opts ...RequestOption) (res *MasterOrderActionV2Reply, err error) {
	if s.masterOrderId == "" {
		return nil, errors.New("masterOrderId is required")
	}
	if s.executionDurationSeconds != nil && *s.executionDurationSeconds <= 10 {
		return nil, errors.New("executionDurationSeconds must be greater than 10")
	}
	if s.povLimit != nil {
		if err := validatePovLimitV2(*s.povLimit); err != nil {
			return nil, err
		}
	}
	endpoint := fmt.Sprintf("%s/%s/update", v2MasterOrdersEndpoint, s.masterOrderId)
	body := params{}
	if s.totalQuantity != nil {
		body["totalQuantity"] = *s.totalQuantity
	}
	if s.orderNotional != nil {
		body["orderNotional"] = *s.orderNotional
	}
	if s.upTolerance != nil {
		body["upTolerance"] = *s.upTolerance
	}
	if s.lowTolerance != nil {
		body["lowTolerance"] = *s.lowTolerance
	}
	if s.enableMake != nil {
		body["enableMake"] = *s.enableMake
	}
	if s.makerRateLimit != nil {
		body["makerRateLimit"] = *s.makerRateLimit
	}
	if s.strictUpBound != nil {
		body["strictUpBound"] = *s.strictUpBound
	}
	if s.povLimit != nil {
		body["povLimit"] = *s.povLimit
	}
	if s.povMinLimit != nil {
		body["povMinLimit"] = *s.povMinLimit
	}
	if s.worstPrice != nil {
		body["worstPrice"] = *s.worstPrice
	}
	if s.tailOrderProtection != nil {
		body["tailOrderProtection"] = *s.tailOrderProtection
	}
	if s.mustComplete != nil {
		body["mustComplete"] = *s.mustComplete
	}
	if s.executionDurationSeconds != nil {
		body["executionDurationSeconds"] = *s.executionDurationSeconds
	}
	data, err := s.c.callAPIV2WithJSONBody(ctx, http.MethodPut, endpoint, body, opts...)
	if err != nil {
		return nil, err
	}
	res = new(MasterOrderActionV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// BatchCancelMasterOrdersV2Service cancels multiple master orders in one call.
type BatchCancelMasterOrdersV2Service struct {
	c              *Client
	masterOrderIds []string
	reason         *string
}

// MasterOrderIds sets the list of master order IDs to cancel.
func (s *BatchCancelMasterOrdersV2Service) MasterOrderIds(ids []string) *BatchCancelMasterOrdersV2Service {
	s.masterOrderIds = ids
	return s
}

// AddMasterOrderId appends a master order ID to the batch.
func (s *BatchCancelMasterOrdersV2Service) AddMasterOrderId(id string) *BatchCancelMasterOrdersV2Service {
	s.masterOrderIds = append(s.masterOrderIds, id)
	return s
}

// Reason sets the optional batch cancellation reason.
func (s *BatchCancelMasterOrdersV2Service) Reason(reason string) *BatchCancelMasterOrdersV2Service {
	s.reason = &reason
	return s
}

// Do sends the request.
func (s *BatchCancelMasterOrdersV2Service) Do(ctx context.Context, opts ...RequestOption) (res *BatchCancelMasterOrdersV2Reply, err error) {
	if len(s.masterOrderIds) == 0 {
		return nil, errors.New("masterOrderIds must not be empty")
	}
	body := params{
		"masterOrderIds": s.masterOrderIds,
	}
	if s.reason != nil {
		body["reason"] = *s.reason
	}
	data, err := s.c.callAPIV2WithJSONBody(ctx, http.MethodPut, v2BatchCancelEndpoint, body, opts...)
	if err != nil {
		return nil, err
	}
	res = new(BatchCancelMasterOrdersV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// BatchCancelMasterOrdersV2Reply is the response of batch-cancel.
type BatchCancelMasterOrdersV2Reply struct {
	SuccessCount int32                          `json:"successCount"`
	FailedOrders []BatchCancelV2FailedOrderInfo `json:"failedOrders"`
}

// BatchCancelV2FailedOrderInfo describes one failed cancellation attempt.
type BatchCancelV2FailedOrderInfo struct {
	MasterOrderId string `json:"masterOrderId"`
	Reason        string `json:"reason"`
}
