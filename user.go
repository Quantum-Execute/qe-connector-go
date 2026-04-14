package qe_connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
	"github.com/Quantum-Execute/qe-connector-go/dto/algorithm_dto"
)

// ListExchangeApisService list exchange APIs
type ListExchangeApisService struct {
	c        *Client
	page     *int32
	pageSize *int32
	exchange *trading_enums.Exchange
}

// Page set page
func (s *ListExchangeApisService) Page(page int32) *ListExchangeApisService {
	s.page = &page
	return s
}

// PageSize set pageSize
func (s *ListExchangeApisService) PageSize(pageSize int32) *ListExchangeApisService {
	s.pageSize = &pageSize
	return s
}

// Exchange set exchange
func (s *ListExchangeApisService) Exchange(exchange trading_enums.Exchange) *ListExchangeApisService {
	s.exchange = &exchange
	return s
}

// Do send request
func (s *ListExchangeApisService) Do(ctx context.Context, opts ...RequestOption) (res *ListExchangeApisReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis",
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
	res = new(ListExchangeApisReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListExchangeApisReply list exchange APIs response
type ListExchangeApisReply struct {
	Items    []ExchangeApiInfo `json:"items"`
	Total    int32             `json:"total"`
	Page     int32             `json:"page"`
	PageSize int32             `json:"pageSize"`
}

// ExchangeApiInfo exchange API info
type ExchangeApiInfo struct {
	Id                 string `json:"id"`
	CreatedAt          string `json:"createdAt"`
	AccountName        string `json:"accountName"`
	Exchange           string `json:"exchange"`
	ApiKey             string `json:"apiKey"`
	VerificationMethod string `json:"verificationMethod"`
	Status             string `json:"status"`
	IsValid            bool   `json:"isValid"`
	IsTradingEnabled   bool   `json:"isTradingEnabled"`
	IsDefault          bool   `json:"isDefault"`
	IsPm               bool   `json:"isPm"`
}

// GetMasterOrdersService get master orders
type GetMasterOrdersService struct {
	c         *Client
	page      *int32
	pageSize  *int32
	status    *trading_enums.MasterOrderStatus
	exchange  *string
	symbol    *string
	startTime *string
	endTime   *string
}

// Page set page
func (s *GetMasterOrdersService) Page(page int32) *GetMasterOrdersService {
	s.page = &page
	return s
}

// PageSize set pageSize
func (s *GetMasterOrdersService) PageSize(pageSize int32) *GetMasterOrdersService {
	s.pageSize = &pageSize
	return s
}

// Status set status
func (s *GetMasterOrdersService) Status(status trading_enums.MasterOrderStatus) *GetMasterOrdersService {
	s.status = &status
	return s
}

// Exchange set exchange
func (s *GetMasterOrdersService) Exchange(exchange string) *GetMasterOrdersService {
	s.exchange = &exchange
	return s
}

// Symbol set symbol
func (s *GetMasterOrdersService) Symbol(symbol string) *GetMasterOrdersService {
	s.symbol = &symbol
	return s
}

// StartTime set startTime
func (s *GetMasterOrdersService) StartTime(startTime string) *GetMasterOrdersService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *GetMasterOrdersService) EndTime(endTime string) *GetMasterOrdersService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *GetMasterOrdersService) Do(ctx context.Context, opts ...RequestOption) (res *GetMasterOrdersReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/trading/master-orders",
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
		m["status"] = *s.status
	}
	if s.exchange != nil {
		m["exchange"] = *s.exchange
	}
	if s.symbol != nil {
		m["symbol"] = *s.symbol
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
	res = new(GetMasterOrdersReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetMasterOrdersReply get master orders response
type GetMasterOrdersReply struct {
	Items    []MasterOrderInfo `json:"items"`
	Total    string            `json:"total"`
	Page     int32             `json:"page"`
	PageSize int32             `json:"pageSize"`
}

// GetMasterOrderDetailService get master order detail
type GetMasterOrderDetailService struct {
	c             *Client
	masterOrderId string
}

// MasterOrderId set masterOrderId
func (s *GetMasterOrderDetailService) MasterOrderId(masterOrderId string) *GetMasterOrderDetailService {
	s.masterOrderId = masterOrderId
	return s
}

// Do send request
func (s *GetMasterOrderDetailService) Do(ctx context.Context, opts ...RequestOption) (res *GetMasterOrderDetailReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/user/trading/master-orders/%s", s.masterOrderId),
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(GetMasterOrderDetailReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetMasterOrderDetailReply get master order detail response
type GetMasterOrderDetailReply struct {
	MasterOrder MasterOrderInfo `json:"masterOrder"`
}

// GetMasterOrderDetailByClientOrderIdService get master order detail by client order id
type GetMasterOrderDetailByClientOrderIdService struct {
	c             *Client
	clientOrderId string
}

// ClientOrderId set clientOrderId
func (s *GetMasterOrderDetailByClientOrderIdService) ClientOrderId(clientOrderId string) *GetMasterOrderDetailByClientOrderIdService {
	s.clientOrderId = clientOrderId
	return s
}

// Do send request
func (s *GetMasterOrderDetailByClientOrderIdService) Do(ctx context.Context, opts ...RequestOption) (res *GetMasterOrderDetailReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/user/trading/master-orders/by-client-order-id/%s", s.clientOrderId),
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(GetMasterOrderDetailReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MasterOrderInfo master order info
type MasterOrderInfo struct {
	MasterOrderId            string    `json:"masterOrderId"`
	Algorithm                string    `json:"algorithm"`
	AlgorithmType            string    `json:"algorithmType"`
	Exchange                 string    `json:"exchange"`
	Symbol                   string    `json:"symbol"`
	MarketType               string    `json:"marketType"`
	Side                     string    `json:"side"`
	TotalQuantity            float64   `json:"totalQuantity"`
	FilledQuantity           float64   `json:"filledQuantity"`
	AveragePrice             float64   `json:"averagePrice"`
	Status                   string    `json:"status"`
	ExecutionDuration        int32     `json:"executionDuration"`
	ExecutionDurationSeconds *int32    `json:"executionDurationSeconds,omitempty"`
	PriceLimit               float64   `json:"priceLimit"`
	StartTime                string    `json:"startTime"`
	EndTime                  string    `json:"endTime"`
	CreatedAt                string    `json:"createdAt"`
	UpdatedAt                string    `json:"updatedAt"`
	Notes                    string    `json:"notes"`
	MarginType               string    `json:"marginType"`
	ReduceOnly               bool      `json:"reduceOnly"`
	StrategyType             string    `json:"strategyType"`
	OrderNotional            float64   `json:"orderNotional"`
	MustComplete             bool      `json:"mustComplete"`
	MakerRateLimit           float64   `json:"makerRateLimit"`
	PovLimit                 float64   `json:"povLimit"`
	ClientId                 string    `json:"clientId"`
	Date                     string    `json:"date"`
	TicktimeInt              string    `json:"ticktimeInt"`
	LimitPriceString         string    `json:"limitPriceString"`
	UpTolerance              string    `json:"upTolerance"`
	LowTolerance             string    `json:"lowTolerance"`
	StrictUpBound            bool      `json:"strictUpBound"`
	TicktimeMs               string    `json:"ticktimeMs"`
	Category                 string    `json:"category"`
	FilledAmount             float64   `json:"filledAmount"`
	TotalValue               float64   `json:"totalValue"`
	Base                     string    `json:"base"`
	Quote                    string    `json:"quote"`
	CompletionProgress       float64   `json:"completionProgress"`
	Reason                   string    `json:"reason"`
	TakerMakerRate           float64   `json:"takerMakerRate"`
	MakerRate                float64   `json:"makerRate"`
	TailOrderProtection      bool      `json:"tailOrderProtection"`
	TradingAccount           string    `json:"tradingAccount"`
	EnableMake               bool      `json:"enableMake"`
	ClientOrderId            string    `json:"clientOrderId"`
	FinishedMs               FlexInt64 `json:"finishedMs"`
	WorstPrice               float64   `json:"worstPrice"`
	PovMinLimit              float64   `json:"povMinLimit"`
}

// GetOrderFillsService get order fills
type GetOrderFillsService struct {
	c             *Client
	page          *int32
	pageSize      *int32
	masterOrderId *string
	subOrderId    *string
	orderId       *string
	symbol        *string
	status        *string
	startTime     *string
	endTime       *string
}

// Page set page
func (s *GetOrderFillsService) Page(page int32) *GetOrderFillsService {
	s.page = &page
	return s
}

// PageSize set pageSize
func (s *GetOrderFillsService) PageSize(pageSize int32) *GetOrderFillsService {
	s.pageSize = &pageSize
	return s
}

// MasterOrderId set masterOrderId
func (s *GetOrderFillsService) MasterOrderId(masterOrderId string) *GetOrderFillsService {
	s.masterOrderId = &masterOrderId
	return s
}

// SubOrderId set subOrderId
func (s *GetOrderFillsService) SubOrderId(subOrderId string) *GetOrderFillsService {
	s.subOrderId = &subOrderId
	return s
}

// OrderId set orderId (exchange order ID filter)
func (s *GetOrderFillsService) OrderId(orderId string) *GetOrderFillsService {
	s.orderId = &orderId
	return s
}

// Symbol set symbol
func (s *GetOrderFillsService) Symbol(symbol string) *GetOrderFillsService {
	s.symbol = &symbol
	return s
}

// Status set symbol
func (s *GetOrderFillsService) Status(status string) *GetOrderFillsService {
	s.status = &status
	return s
}

// StartTime set startTime
func (s *GetOrderFillsService) StartTime(startTime string) *GetOrderFillsService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *GetOrderFillsService) EndTime(endTime string) *GetOrderFillsService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *GetOrderFillsService) Do(ctx context.Context, opts ...RequestOption) (res *GetOrderFillsReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/trading/order-fills",
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
	if s.subOrderId != nil {
		m["subOrderId"] = *s.subOrderId
	}
	if s.orderId != nil {
		m["orderId"] = *s.orderId
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
	res = new(GetOrderFillsReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetOrderFillsReply get order fills response
type GetOrderFillsReply struct {
	Items    []OrderFillInfo `json:"items"`
	Total    string          `json:"total"`
	Page     int32           `json:"page"`
	PageSize int32           `json:"pageSize"`
}

// OrderFillInfo order fill info
type OrderFillInfo struct {
	Id               string  `json:"id"`
	OrderCreatedTime string  `json:"orderCreatedTime"`
	MasterOrderId    string  `json:"masterOrderId"`
	Exchange         string  `json:"exchange"`
	Category         string  `json:"category"`
	Symbol           string  `json:"symbol"`
	Side             string  `json:"side"`
	FilledValue      float64 `json:"filledValue"`
	FilledQuantity   float64 `json:"filledQuantity"`
	AvgPrice         float64 `json:"avgPrice"`
	Price            float64 `json:"price"`
	Fee              float64 `json:"fee"`
	TradingAccount   string  `json:"tradingAccount"`
	Status           string  `json:"status"`
	RejectReason     string  `json:"rejectReason"`
	Base             string  `json:"base"`
	Quote            string  `json:"quote"`
	Type             string  `json:"type"`
	OrderId          string  `json:"orderId"`
	Quantity         float64 `json:"quantity"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

// CreateMasterOrderService create master order
type CreateMasterOrderService struct {
	c                        *Client
	algorithm                trading_enums.Algorithm
	exchange                 trading_enums.Exchange
	symbol                   string
	marketType               trading_enums.MarketType
	side                     trading_enums.OrderSide
	totalQuantity            *float64
	orderNotional            *float64
	apiKeyId                 string
	strategyType             *trading_enums.StrategyType
	startTime                *string
	executionDuration        *int32
	executionDurationSeconds *int32
	limitPrice               *float64
	mustComplete             *bool
	makerRateLimit           *float64
	povLimit                 *float64
	marginType               *trading_enums.MarginType
	reduceOnly               *bool
	notes                    *string
	upTolerance              *string
	lowTolerance             *string
	strictUpBound            *bool
	povMinLimit              *float64
	tailOrderProtection      *bool
	isTargetPosition         *bool
	isMargin                 *bool
	enableMake               *bool
	clientOrderId            *string
	worstPrice               *float64
}

// Algorithm set algorithm
func (s *CreateMasterOrderService) Algorithm(algorithm trading_enums.Algorithm) *CreateMasterOrderService {
	s.algorithm = algorithm
	return s
}

// Exchange set exchange
func (s *CreateMasterOrderService) Exchange(exchange trading_enums.Exchange) *CreateMasterOrderService {
	s.exchange = exchange
	return s
}

// Symbol set symbol
func (s *CreateMasterOrderService) Symbol(symbol string) *CreateMasterOrderService {
	s.symbol = symbol
	return s
}

// MarketType set marketType
func (s *CreateMasterOrderService) MarketType(marketType trading_enums.MarketType) *CreateMasterOrderService {
	s.marketType = marketType
	return s
}

// Side set side
func (s *CreateMasterOrderService) Side(side trading_enums.OrderSide) *CreateMasterOrderService {
	s.side = side
	return s
}

// TotalQuantity set totalQuantity
func (s *CreateMasterOrderService) TotalQuantity(totalQuantity float64) *CreateMasterOrderService {
	s.totalQuantity = &totalQuantity
	return s
}

// OrderNotional set orderNotional
func (s *CreateMasterOrderService) OrderNotional(orderNotional float64) *CreateMasterOrderService {
	s.orderNotional = &orderNotional
	return s
}

// ApiKeyId set apiKeyId
func (s *CreateMasterOrderService) ApiKeyId(apiKeyId string) *CreateMasterOrderService {
	s.apiKeyId = apiKeyId
	return s
}

// StrategyType set strategyType
func (s *CreateMasterOrderService) StrategyType(strategyType trading_enums.StrategyType) *CreateMasterOrderService {
	s.strategyType = &strategyType
	return s
}

// StartTime set startTime
func (s *CreateMasterOrderService) StartTime(startTime string) *CreateMasterOrderService {
	s.startTime = &startTime
	return s
}

// ExecutionDuration set executionDuration
func (s *CreateMasterOrderService) ExecutionDuration(executionDuration int32) *CreateMasterOrderService {
	s.executionDuration = &executionDuration
	return s
}

// ExecutionDurationSeconds set executionDurationSeconds
//
// Note: Only used for TWAP-1. When provided and > 0, it takes precedence over executionDuration (minutes).
// It must be greater than 10 seconds.
func (s *CreateMasterOrderService) ExecutionDurationSeconds(executionDurationSeconds int32) *CreateMasterOrderService {
	s.executionDurationSeconds = &executionDurationSeconds
	return s
}

// EndTime set endTime
//
// Deprecated: EndTime is deprecated and no longer used.
// This method is kept for backward compatibility but does nothing.
// The endTime field has been removed from the API.
func (s *CreateMasterOrderService) EndTime(endTime string) *CreateMasterOrderService {
	// No-op: endTime is deprecated and no longer sent to the API
	return s
}

// LimitPrice set limitPrice
//
// Deprecated: Use WorstPrice instead. LimitPrice is kept for backward compatibility.
func (s *CreateMasterOrderService) LimitPrice(limitPrice float64) *CreateMasterOrderService {
	s.limitPrice = &limitPrice
	return s
}

// WorstPrice set worstPrice (worst acceptable price)
//
// The worst acceptable trading price. For buy orders this is the maximum buy price;
// for sell orders this is the minimum sell price. Use -1 for no limit.
func (s *CreateMasterOrderService) WorstPrice(worstPrice float64) *CreateMasterOrderService {
	s.worstPrice = &worstPrice
	return s
}

// MustComplete set mustComplete
func (s *CreateMasterOrderService) MustComplete(mustComplete bool) *CreateMasterOrderService {
	s.mustComplete = &mustComplete
	return s
}

// MakerRateLimit set makerRateLimit
func (s *CreateMasterOrderService) MakerRateLimit(makerRateLimit float64) *CreateMasterOrderService {
	s.makerRateLimit = &makerRateLimit
	return s
}

// PovLimit set povLimit
func (s *CreateMasterOrderService) PovLimit(povLimit float64) *CreateMasterOrderService {
	s.povLimit = &povLimit
	return s
}

// MarginType set marginType
func (s *CreateMasterOrderService) MarginType(marginType trading_enums.MarginType) *CreateMasterOrderService {
	s.marginType = &marginType
	return s
}

// ReduceOnly set reduceOnly
func (s *CreateMasterOrderService) ReduceOnly(reduceOnly bool) *CreateMasterOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// Notes set notes
func (s *CreateMasterOrderService) Notes(notes string) *CreateMasterOrderService {
	s.notes = &notes
	return s
}

// UpTolerance set upTolerance
func (s *CreateMasterOrderService) UpTolerance(upTolerance string) *CreateMasterOrderService {
	s.upTolerance = &upTolerance
	return s
}

// LowTolerance set lowTolerance
func (s *CreateMasterOrderService) LowTolerance(lowTolerance string) *CreateMasterOrderService {
	s.lowTolerance = &lowTolerance
	return s
}

// StrictUpBound set strictUpBound
func (s *CreateMasterOrderService) StrictUpBound(strictUpBound bool) *CreateMasterOrderService {
	s.strictUpBound = &strictUpBound
	return s
}

// PovMinLimit set povMinLimit
func (s *CreateMasterOrderService) PovMinLimit(povMinLimit float64) *CreateMasterOrderService {
	s.povMinLimit = &povMinLimit
	return s
}

// TailOrderProtection set tailOrderProtection
func (s *CreateMasterOrderService) TailOrderProtection(tailOrderProtection bool) *CreateMasterOrderService {
	s.tailOrderProtection = &tailOrderProtection
	return s
}

// IsTargetPosition set tailOrderProtection
func (s *CreateMasterOrderService) IsTargetPosition(isTargetPosition bool) *CreateMasterOrderService {
	s.isTargetPosition = &isTargetPosition
	return s
}

// IsMargin set isMargin
func (s *CreateMasterOrderService) IsMargin(isMargin bool) *CreateMasterOrderService {
	s.isMargin = &isMargin
	return s
}

// EnableMake set isMargin
func (s *CreateMasterOrderService) EnableMake(enableMake bool) *CreateMasterOrderService {
	s.enableMake = &enableMake
	return s
}

// ClientOrderId set clientOrderId
func (s *CreateMasterOrderService) ClientOrderId(clientOrderId string) *CreateMasterOrderService {
	s.clientOrderId = &clientOrderId
	return s
}

// Do send request
func (s *CreateMasterOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateMasterOrderReply, err error) {
	// Deribit special rules:
	// - When trading BTCUSD/ETHUSD, only totalQuantity is allowed, and orderNotional is not allowed.
	if s.exchange == trading_enums.ExchangeDeribit &&
		(strings.EqualFold(s.symbol, "BTCUSD") || strings.EqualFold(s.symbol, "ETHUSD")) {
		if s.orderNotional != nil {
			return nil, errors.New("orderNotional is not allowed when exchange is Deribit and symbol is BTCUSD or ETHUSD; use totalQuantity (unit: USD) instead")
		}
		if s.totalQuantity == nil {
			return nil, errors.New("totalQuantity is required when exchange is Deribit and symbol is BTCUSD or ETHUSD (unit: USD)")
		}
	}

	// Binance coin-margined perp special rules:
	// - When trading Binance PERP with marginType=C, only totalQuantity is allowed.
	// - totalQuantity unit is contracts and must be an integer.
	if s.exchange == trading_enums.ExchangeBinance &&
		s.marketType == trading_enums.MarketTypePerp &&
		s.marginType != nil &&
		*s.marginType == trading_enums.MarginTypeC {
		if s.orderNotional != nil {
			return nil, errors.New("orderNotional is not allowed when exchange is Binance and marginType is C for PERP orders; use totalQuantity (unit: contracts) instead")
		}
		if s.totalQuantity == nil {
			return nil, errors.New("totalQuantity is required when exchange is Binance and marginType is C for PERP orders (unit: contracts)")
		}
		if math.Trunc(*s.totalQuantity) != *s.totalQuantity {
			return nil, errors.New("totalQuantity must be an integer when exchange is Binance and marginType is C for PERP orders (unit: contracts)")
		}
	}

	r := &request{
		method:   http.MethodPost,
		endpoint: "/user/trading/master-orders",
		secType:  secTypeSigned,
	}
	m := params{
		"algorithm":     s.algorithm,
		"algorithmType": "TWAP",
		"exchange":      s.exchange,
		"symbol":        s.symbol,
		"marketType":    s.marketType,
		"side":          s.side,
		"apiKeyId":      s.apiKeyId,
	}
	if s.totalQuantity != nil {
		m["totalQuantity"] = *s.totalQuantity
	}
	if s.orderNotional != nil {
		m["orderNotional"] = *s.orderNotional
	}
	if s.strategyType != nil {
		m["strategyType"] = *s.strategyType
	}
	if s.startTime != nil {
		m["startTime"] = *s.startTime
	}
	if s.executionDuration != nil {
		m["executionDuration"] = *s.executionDuration
	}
	if s.executionDurationSeconds != nil {
		m["executionDurationSeconds"] = *s.executionDurationSeconds
	}
	if s.limitPrice != nil {
		m["limitPrice"] = *s.limitPrice
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
	if s.povLimit != nil {
		m["povLimit"] = *s.povLimit
	}
	if s.marginType != nil {
		m["marginType"] = *s.marginType
	}
	if s.reduceOnly != nil {
		m["reduceOnly"] = *s.reduceOnly
	}
	if s.notes != nil {
		m["notes"] = *s.notes
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
	if s.povMinLimit != nil {
		m["povMinLimit"] = *s.povMinLimit
	}
	if s.tailOrderProtection != nil {
		m["tailOrderProtection"] = *s.tailOrderProtection
	} else {
		m["tailOrderProtection"] = true
	}
	if s.enableMake != nil {
		m["enableMake"] = *s.enableMake
	} else {
		m["enableMake"] = true
	}
	if s.isTargetPosition != nil {
		m["isTargetPosition"] = *s.isTargetPosition
		if *s.isTargetPosition {
			if s.totalQuantity == nil || s.orderNotional != nil {
				return nil, errors.New("totalQuantity is required and orderNotional not required when isTargetPosition is true")
			}
		}
	} else {
		m["isTargetPosition"] = false
	}
	if s.isMargin != nil {
		m["isMargin"] = *s.isMargin
	}
	if s.clientOrderId != nil {
		m["clientOrderId"] = *s.clientOrderId
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateMasterOrderReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateMasterOrderReply create master order response
type CreateMasterOrderReply struct {
	MasterOrderId string `json:"masterOrderId"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
}

// CancelMasterOrderService cancel master order
type CancelMasterOrderService struct {
	c             *Client
	masterOrderId string
	reason        *string
}

// MasterOrderId set masterOrderId
func (s *CancelMasterOrderService) MasterOrderId(masterOrderId string) *CancelMasterOrderService {
	s.masterOrderId = masterOrderId
	return s
}

// Reason set reason
func (s *CancelMasterOrderService) Reason(reason string) *CancelMasterOrderService {
	s.reason = &reason
	return s
}

// Do send request
func (s *CancelMasterOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelMasterOrderReply, err error) {
	r := &request{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("/user/trading/master-orders/%s/cancel", s.masterOrderId),
		secType:  secTypeSigned,
	}
	m := params{
		"masterOrderId": s.masterOrderId,
	}
	if s.reason != nil {
		m["reason"] = *s.reason
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CancelMasterOrderReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CancelMasterOrderReply cancel master order response
type CancelMasterOrderReply struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PauseMasterOrderService pause a running master order
type PauseMasterOrderService struct {
	c             *Client
	masterOrderId string
	reason        *string
}

// MasterOrderId set masterOrderId
func (s *PauseMasterOrderService) MasterOrderId(masterOrderId string) *PauseMasterOrderService {
	s.masterOrderId = masterOrderId
	return s
}

// Reason set reason (optional)
func (s *PauseMasterOrderService) Reason(reason string) *PauseMasterOrderService {
	s.reason = &reason
	return s
}

// Do send request
func (s *PauseMasterOrderService) Do(ctx context.Context, opts ...RequestOption) (res *PauseMasterOrderReply, err error) {
	r := &request{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("/user/trading/master-orders/%s/pause", s.masterOrderId),
		secType:  secTypeSigned,
	}
	m := params{
		"masterOrderId": s.masterOrderId,
	}
	if s.reason != nil {
		m["reason"] = *s.reason
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(PauseMasterOrderReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PauseMasterOrderReply pause master order response
type PauseMasterOrderReply struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ResumeMasterOrderService resume a paused master order
type ResumeMasterOrderService struct {
	c             *Client
	masterOrderId string
}

// MasterOrderId set masterOrderId
func (s *ResumeMasterOrderService) MasterOrderId(masterOrderId string) *ResumeMasterOrderService {
	s.masterOrderId = masterOrderId
	return s
}

// Do send request
func (s *ResumeMasterOrderService) Do(ctx context.Context, opts ...RequestOption) (res *ResumeMasterOrderReply, err error) {
	r := &request{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("/user/trading/master-orders/%s/resume", s.masterOrderId),
		secType:  secTypeSigned,
	}
	m := params{
		"masterOrderId": s.masterOrderId,
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(ResumeMasterOrderReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ResumeMasterOrderReply resume master order response
type ResumeMasterOrderReply struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateMasterOrderParamsService update parameters of a running master order
type UpdateMasterOrderParamsService struct {
	c                        *Client
	masterOrderId            string
	orderNotional            *float64
	totalQuantity            *float64
	upTolerance              *string
	lowTolerance             *string
	enableMake               *bool
	makerRateLimit           *float64
	strictUpBound            *bool
	povLimit                 *float64
	povMinLimit              *float64
	limitPrice               *float64
	worstPrice               *float64
	tailOrderProtection      *bool
	mustComplete             *bool
	executionDurationSeconds *int32
	executionDuration        *int32
}

// MasterOrderId set masterOrderId (required)
func (s *UpdateMasterOrderParamsService) MasterOrderId(masterOrderId string) *UpdateMasterOrderParamsService {
	s.masterOrderId = masterOrderId
	return s
}

// OrderNotional set orderNotional
func (s *UpdateMasterOrderParamsService) OrderNotional(orderNotional float64) *UpdateMasterOrderParamsService {
	s.orderNotional = &orderNotional
	return s
}

// TotalQuantity set totalQuantity
func (s *UpdateMasterOrderParamsService) TotalQuantity(totalQuantity float64) *UpdateMasterOrderParamsService {
	s.totalQuantity = &totalQuantity
	return s
}

// UpTolerance set upTolerance
func (s *UpdateMasterOrderParamsService) UpTolerance(upTolerance string) *UpdateMasterOrderParamsService {
	s.upTolerance = &upTolerance
	return s
}

// LowTolerance set lowTolerance
func (s *UpdateMasterOrderParamsService) LowTolerance(lowTolerance string) *UpdateMasterOrderParamsService {
	s.lowTolerance = &lowTolerance
	return s
}

// EnableMake set enableMake
func (s *UpdateMasterOrderParamsService) EnableMake(enableMake bool) *UpdateMasterOrderParamsService {
	s.enableMake = &enableMake
	return s
}

// MakerRateLimit set makerRateLimit
func (s *UpdateMasterOrderParamsService) MakerRateLimit(makerRateLimit float64) *UpdateMasterOrderParamsService {
	s.makerRateLimit = &makerRateLimit
	return s
}

// StrictUpBound set strictUpBound
func (s *UpdateMasterOrderParamsService) StrictUpBound(strictUpBound bool) *UpdateMasterOrderParamsService {
	s.strictUpBound = &strictUpBound
	return s
}

// PovLimit set povLimit
func (s *UpdateMasterOrderParamsService) PovLimit(povLimit float64) *UpdateMasterOrderParamsService {
	s.povLimit = &povLimit
	return s
}

// PovMinLimit set povMinLimit
func (s *UpdateMasterOrderParamsService) PovMinLimit(povMinLimit float64) *UpdateMasterOrderParamsService {
	s.povMinLimit = &povMinLimit
	return s
}

// LimitPrice set limitPrice
func (s *UpdateMasterOrderParamsService) LimitPrice(limitPrice float64) *UpdateMasterOrderParamsService {
	s.limitPrice = &limitPrice
	return s
}

// WorstPrice set worstPrice
func (s *UpdateMasterOrderParamsService) WorstPrice(worstPrice float64) *UpdateMasterOrderParamsService {
	s.worstPrice = &worstPrice
	return s
}

// TailOrderProtection set tailOrderProtection
func (s *UpdateMasterOrderParamsService) TailOrderProtection(tailOrderProtection bool) *UpdateMasterOrderParamsService {
	s.tailOrderProtection = &tailOrderProtection
	return s
}

// MustComplete set mustComplete
func (s *UpdateMasterOrderParamsService) MustComplete(mustComplete bool) *UpdateMasterOrderParamsService {
	s.mustComplete = &mustComplete
	return s
}

// ExecutionDurationSeconds set executionDurationSeconds (in seconds, must be > 10).
// Mutually exclusive with ExecutionDuration — only one of the two may be provided.
func (s *UpdateMasterOrderParamsService) ExecutionDurationSeconds(executionDurationSeconds int32) *UpdateMasterOrderParamsService {
	s.executionDurationSeconds = &executionDurationSeconds
	return s
}

// ExecutionDuration set executionDuration (in minutes, must be >= 1).
// Mutually exclusive with ExecutionDurationSeconds — only one of the two may be provided.
func (s *UpdateMasterOrderParamsService) ExecutionDuration(executionDuration int32) *UpdateMasterOrderParamsService {
	s.executionDuration = &executionDuration
	return s
}

// Do send request
func (s *UpdateMasterOrderParamsService) Do(ctx context.Context, opts ...RequestOption) (res *UpdateMasterOrderParamsReply, err error) {
	r := &request{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("/user/trading/master-orders/%s/update", s.masterOrderId),
		secType:  secTypeSigned,
	}
	m := params{
		"masterOrderId": s.masterOrderId,
	}
	if s.orderNotional != nil {
		m["orderNotional"] = *s.orderNotional
	}
	if s.totalQuantity != nil {
		m["totalQuantity"] = *s.totalQuantity
	}
	if s.upTolerance != nil {
		m["upTolerance"] = *s.upTolerance
	}
	if s.lowTolerance != nil {
		m["lowTolerance"] = *s.lowTolerance
	}
	if s.enableMake != nil {
		m["enableMake"] = *s.enableMake
	}
	if s.makerRateLimit != nil {
		m["makerRateLimit"] = *s.makerRateLimit
	}
	if s.strictUpBound != nil {
		m["strictUpBound"] = *s.strictUpBound
	}
	if s.povLimit != nil {
		m["povLimit"] = *s.povLimit
	}
	if s.povMinLimit != nil {
		m["povMinLimit"] = *s.povMinLimit
	}
	if s.limitPrice != nil {
		m["limitPrice"] = *s.limitPrice
	}
	if s.worstPrice != nil {
		m["worstPrice"] = *s.worstPrice
	}
	if s.tailOrderProtection != nil {
		m["tailOrderProtection"] = *s.tailOrderProtection
	}
	if s.mustComplete != nil {
		m["mustComplete"] = *s.mustComplete
	}
	if s.executionDurationSeconds != nil {
		m["executionDurationSeconds"] = *s.executionDurationSeconds
	}
	if s.executionDuration != nil {
		m["executionDuration"] = *s.executionDuration
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(UpdateMasterOrderParamsReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateMasterOrderParamsReply update master order params response
type UpdateMasterOrderParamsReply struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateListenKeyService create listen key
type CreateListenKeyService struct {
	c *Client
}

// Do send request
func (s *CreateListenKeyService) Do(ctx context.Context, opts ...RequestOption) (res *CreateListenKeyReply, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/user/trading/listen-key",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateListenKeyReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateListenKeyReply create listen key response
type CreateListenKeyReply struct {
	ListenKey string `json:"listenKey"`
	ExpireAt  string `json:"expireAt"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

// GetTcaAnalysisService get TCA analysis full data list
type GetTcaAnalysisService struct {
	c         *Client
	symbol    *string
	category  *string
	apikey    *string
	startTime *int64
	endTime   *int64
}

// Symbol set symbol
func (s *GetTcaAnalysisService) Symbol(symbol string) *GetTcaAnalysisService {
	s.symbol = &symbol
	return s
}

// Category set category
func (s *GetTcaAnalysisService) Category(category string) *GetTcaAnalysisService {
	s.category = &category
	return s
}

// Apikey set apikey (comma-separated supported by server)
func (s *GetTcaAnalysisService) Apikey(apikey string) *GetTcaAnalysisService {
	s.apikey = &apikey
	return s
}

// StartTime set startTime (unix milli)
func (s *GetTcaAnalysisService) StartTime(startTime int64) *GetTcaAnalysisService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime (unix milli)
func (s *GetTcaAnalysisService) EndTime(endTime int64) *GetTcaAnalysisService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *GetTcaAnalysisService) Do(ctx context.Context, opts ...RequestOption) (res []*algorithm_dto.TCAAnalysisResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/trading/tca-analysis",
		secType:  secTypeSigned,
	}
	m := params{}
	if s.symbol != nil {
		m["symbol"] = *s.symbol
	}
	if s.category != nil {
		m["category"] = *s.category
	}
	if s.apikey != nil {
		m["apikey"] = *s.apikey
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
	res = make([]*algorithm_dto.TCAAnalysisResponse, 0)
	if len(data) == 0 {
		return res, nil
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
