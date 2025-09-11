package qe_connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

// ListExchangeApisService list exchange APIs
type ListExchangeApisService struct {
	c        *Client
	page     *int32
	pageSize *int32
	exchange *string
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
func (s *ListExchangeApisService) Exchange(exchange string) *ListExchangeApisService {
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
	Id                 string  `json:"id"`
	CreatedAt          string  `json:"createdAt"`
	AccountName        string  `json:"accountName"`
	Exchange           string  `json:"exchange"`
	ApiKey             string  `json:"apiKey"`
	VerificationMethod string  `json:"verificationMethod"`
	Balance            float64 `json:"balance"`
	Status             string  `json:"status"`
	IsValid            bool    `json:"isValid"`
	IsTradingEnabled   bool    `json:"isTradingEnabled"`
	IsDefault          bool    `json:"isDefault"`
	IsPm               bool    `json:"isPm"`
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

// MasterOrderInfo master order info
type MasterOrderInfo struct {
	MasterOrderId      string  `json:"masterOrderId"`
	Algorithm          string  `json:"algorithm"`
	AlgorithmType      string  `json:"algorithmType"`
	Exchange           string  `json:"exchange"`
	Symbol             string  `json:"symbol"`
	MarketType         string  `json:"marketType"`
	Side               string  `json:"side"`
	TotalQuantity      float64 `json:"totalQuantity"`
	FilledQuantity     float64 `json:"filledQuantity"`
	AveragePrice       float64 `json:"averagePrice"`
	Status             string  `json:"status"`
	ExecutionDuration  int32   `json:"executionDuration"`
	PriceLimit         float64 `json:"priceLimit"`
	StartTime          string  `json:"startTime"`
	EndTime            string  `json:"endTime"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
	Notes              string  `json:"notes"`
	MarginType         string  `json:"marginType"`
	ReduceOnly         bool    `json:"reduceOnly"`
	StrategyType       string  `json:"strategyType"`
	OrderNotional      float64 `json:"orderNotional"`
	MustComplete       bool    `json:"mustComplete"`
	MakerRateLimit     float64 `json:"makerRateLimit"`
	PovLimit           float64 `json:"povLimit"`
	ClientId           string  `json:"clientId"`
	Date               string  `json:"date"`
	TicktimeInt        string  `json:"ticktimeInt"`
	LimitPriceString   string  `json:"limitPriceString"`
	UpTolerance        string  `json:"upTolerance"`
	LowTolerance       string  `json:"lowTolerance"`
	StrictUpBound      bool    `json:"strictUpBound"`
	TicktimeMs         string  `json:"ticktimeMs"`
	Category           string  `json:"category"`
	FilledAmount       float64 `json:"filledAmount"`
	TotalValue         float64 `json:"totalValue"`
	Base               string  `json:"base"`
	Quote              string  `json:"quote"`
	CompletionProgress float64 `json:"completionProgress"`
	Reason             string  `json:"reason"`
}

// GetOrderFillsService get order fills
type GetOrderFillsService struct {
	c             *Client
	page          *int32
	pageSize      *int32
	masterOrderId *string
	subOrderId    *string
	symbol        *string
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

// Symbol set symbol
func (s *GetOrderFillsService) Symbol(symbol string) *GetOrderFillsService {
	s.symbol = &symbol
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
}

// CreateMasterOrderService create master order
type CreateMasterOrderService struct {
	c                   *Client
	algorithm           trading_enums.Algorithm
	algorithmType       string
	exchange            trading_enums.Exchange
	symbol              string
	marketType          trading_enums.MarketType
	side                trading_enums.OrderSide
	totalQuantity       *float64
	orderNotional       *float64
	apiKeyId            string
	strategyType        *trading_enums.StrategyType
	startTime           *string
	executionDuration   *string
	endTime             *string
	limitPrice          *float64
	mustComplete        *bool
	makerRateLimit      *float64
	povLimit            *float64
	marginType          *trading_enums.MarginType
	reduceOnly          *bool
	notes               *string
	worstPrice          *float64
	limitPriceString    *string
	upTolerance         *string
	lowTolerance        *string
	strictUpBound       *bool
	povMinLimit         *float64
	tailOrderProtection *bool
	isTargetPosition    *bool
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
func (s *CreateMasterOrderService) ExecutionDuration(executionDuration string) *CreateMasterOrderService {
	s.executionDuration = &executionDuration
	return s
}

// EndTime set endTime
func (s *CreateMasterOrderService) EndTime(endTime string) *CreateMasterOrderService {
	s.endTime = &endTime
	return s
}

// LimitPrice set limitPrice
func (s *CreateMasterOrderService) LimitPrice(limitPrice float64) *CreateMasterOrderService {
	s.limitPrice = &limitPrice
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

// WorstPrice set worstPrice
func (s *CreateMasterOrderService) WorstPrice(worstPrice float64) *CreateMasterOrderService {
	s.worstPrice = &worstPrice
	return s
}

// LimitPriceString set limitPriceString
func (s *CreateMasterOrderService) LimitPriceString(limitPriceString string) *CreateMasterOrderService {
	s.limitPriceString = &limitPriceString
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

// Do send request
func (s *CreateMasterOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateMasterOrderReply, err error) {
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
	if s.endTime != nil {
		m["endTime"] = *s.endTime
	}
	if s.limitPrice != nil {
		m["limitPrice"] = *s.limitPrice
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
	if s.worstPrice != nil {
		m["worstPrice"] = *s.worstPrice
	}
	if s.limitPriceString != nil {
		m["limitPriceString"] = *s.limitPriceString
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
