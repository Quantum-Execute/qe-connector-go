package qe_connector

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type CancelReplaceService struct {
	c                       *Client
	symbol                  string
	side                    string
	orderType               string
	cancelReplaceMode       string
	timeInForce             *string
	quantity                *float64
	quoteOrderQty           *float64
	price                   *float64
	cancelNewClientOrderId  *string
	cancelOrigClientOrderId *string
	cancelOrderId           *int64
	newClientOrderId        *string
	strategyId              *int32
	strategyType            *int32
	stopPrice               *float64
	trailingDelta           *int64
	icebergQty              *float64
	newOrderRespType        *string
	selfTradePreventionMode *string
	cancelRestrictions      *string
}

// Symbol set symbol
func (s *CancelReplaceService) Symbol(symbol string) *CancelReplaceService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CancelReplaceService) Side(side string) *CancelReplaceService {
	s.side = side
	return s
}

// OrderType set orderType
func (s *CancelReplaceService) OrderType(orderType string) *CancelReplaceService {
	s.orderType = orderType
	return s
}

// CancelReplaceMode set cancelReplaceMode
func (s *CancelReplaceService) CancelReplaceMode(cancelReplaceMode string) *CancelReplaceService {
	s.cancelReplaceMode = cancelReplaceMode
	return s
}

// TimeInForce set timeInForce
func (s *CancelReplaceService) TimeInForce(timeInForce string) *CancelReplaceService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CancelReplaceService) Quantity(quantity float64) *CancelReplaceService {
	s.quantity = &quantity
	return s
}

// QuoteOrderQty set quoteOrderQty
func (s *CancelReplaceService) QuoteOrderQty(quoteOrderQty float64) *CancelReplaceService {
	s.quoteOrderQty = &quoteOrderQty
	return s
}

// Price set price
func (s *CancelReplaceService) Price(price float64) *CancelReplaceService {
	s.price = &price
	return s
}

// CancelNewClientOrderId set cancelNewClientOrderId
func (s *CancelReplaceService) CancelNewClientOrderId(cancelNewClientOrderId string) *CancelReplaceService {
	s.cancelNewClientOrderId = &cancelNewClientOrderId
	return s
}

// CancelOrigClientOrderId set cancelOrigClientOrderId
func (s *CancelReplaceService) CancelOrigClientOrderId(cancelOrigClientOrderId string) *CancelReplaceService {
	s.cancelOrigClientOrderId = &cancelOrigClientOrderId
	return s
}

// CancelOrderId set cancelOrderId
func (s *CancelReplaceService) CancelOrderId(cancelOrderId int64) *CancelReplaceService {
	s.cancelOrderId = &cancelOrderId
	return s
}

// NewClientOrderId set newClientOrderId
func (s *CancelReplaceService) NewClientOrderId(newClientOrderId string) *CancelReplaceService {
	s.newClientOrderId = &newClientOrderId
	return s
}

// StrategyId set strategyId
func (s *CancelReplaceService) StrategyId(strategyId int32) *CancelReplaceService {
	s.strategyId = &strategyId
	return s
}

// StrategyType set strategyType
func (s *CancelReplaceService) StrategyType(strategyType int32) *CancelReplaceService {
	s.strategyType = &strategyType
	return s
}

// StopPrice set stopPrice
func (s *CancelReplaceService) StopPrice(stopPrice float64) *CancelReplaceService {
	s.stopPrice = &stopPrice
	return s
}

// TrailingDelta set trailingDelta
func (s *CancelReplaceService) TrailingDelta(trailingDelta int64) *CancelReplaceService {
	s.trailingDelta = &trailingDelta
	return s
}

// IcebergQty set icebergQty
func (s *CancelReplaceService) IcebergQty(icebergQty float64) *CancelReplaceService {
	s.icebergQty = &icebergQty
	return s
}

// NewOrderRespType set newOrderRespType
func (s *CancelReplaceService) NewOrderRespType(newOrderRespType string) *CancelReplaceService {
	s.newOrderRespType = &newOrderRespType
	return s
}

// SelfTradePreventionMode set selfTradePreventionMode
func (s *CancelReplaceService) SelfTradePreventionMode(selfTradePreventionMode string) *CancelReplaceService {
	s.selfTradePreventionMode = &selfTradePreventionMode
	return s
}

// CancelRestrictions set cancelRestrictions
func (s *CancelReplaceService) CancelRestrictions(cancelRestrictions string) *CancelReplaceService {
	s.cancelRestrictions = &cancelRestrictions
	return s
}

// Do send request
func (s *CancelReplaceService) Do(ctx context.Context, opts ...RequestOption) (res *CancelReplaceResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/api/v3/order/cancelReplace",
		secType:  secTypeSigned,
	}
	m := params{
		"symbol":            s.symbol,
		"side":              s.side,
		"type":              s.orderType,
		"cancelReplaceMode": s.cancelReplaceMode,
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.quantity != nil {
		m["quantity"] = strconv.FormatFloat(*s.quantity, 'f', -1, 64)
	}
	if s.quoteOrderQty != nil {
		m["quoteOrderQty"] = *s.quoteOrderQty
	}
	if s.price != nil {
		m["price"] = *s.price
	}
	if s.cancelNewClientOrderId != nil {
		m["cancelNewClientOrderId"] = *s.cancelNewClientOrderId
	}
	if s.cancelOrigClientOrderId != nil {
		m["cancelOrigClientOrderId"] = *s.cancelOrigClientOrderId
	}
	if s.cancelOrderId != nil {
		m["cancelOrderId"] = *s.cancelOrderId
	}
	if s.newClientOrderId != nil {
		m["newClientOrderId"] = *s.newClientOrderId
	}
	if s.strategyId != nil {
		m["strategyId"] = *s.strategyId
	}
	if s.strategyType != nil {
		m["strategyType"] = *s.strategyType
	}
	if s.stopPrice != nil {
		m["stopPrice"] = *s.stopPrice
	}
	if s.trailingDelta != nil {
		m["trailingDelta"] = *s.trailingDelta
	}
	if s.icebergQty != nil {
		m["icebergQty"] = *s.icebergQty
	}
	if s.newOrderRespType != nil {
		m["newOrderRespType"] = *s.newOrderRespType
	}
	if s.selfTradePreventionMode != nil {
		m["selfTradePreventionMode"] = *s.selfTradePreventionMode
	}
	if s.cancelRestrictions != nil {
		m["cancelRestrictions"] = *s.cancelRestrictions
	}
	r.setParams(m)
	data, _ := s.c.callAPI(ctx, r, opts...)
	res = new(CancelReplaceResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

type CancelReplaceResponse struct {
	Code           int64  `json:"code,omitempty"`
	Msg            string `json:"msg,omitempty"`
	CancelResult   string `json:"cancelResult,omitempty"`
	NewOrderResult string `json:"newOrderResult,omitempty"`
	CancelResponse *struct {
		Code                    int    `json:"code,omitempty"`
		Msg                     string `json:"msg,omitempty"`
		Symbol                  string `json:"symbol,omitempty"`
		OrigClientOrderId       string `json:"origClientOrderId,omitempty"`
		OrderId                 int64  `json:"orderId,omitempty"`
		OrderListId             int64  `json:"orderListId,omitempty"`
		ClientOrderId           string `json:"clientOrderId,omitempty"`
		Price                   string `json:"price,omitempty"`
		OrigQty                 string `json:"origQty,omitempty"`
		ExecutedQty             string `json:"executedQty,omitempty"`
		CummulativeQuoteQty     string `json:"cummulativeQuoteQty,omitempty"`
		Status                  string `json:"status,omitempty"`
		TimeInForce             string `json:"timeInForce,omitempty"`
		Type                    string `json:"type,omitempty"`
		Side                    string `json:"side,omitempty"`
		SelfTradePreventionMode string `json:"selfTradePreventionMode,omitempty"`
	} `json:"cancelResponse,omitempty"`
	NewOrderResponse *struct {
		Code                int64  `json:"code,omitempty"`
		Msg                 string `json:"msg,omitempty"`
		Symbol              string `json:"symbol,omitempty"`
		OrderId             int64  `json:"orderId,omitempty"`
		OrderListId         int64  `json:"orderListId,omitempty"`
		ClientOrderId       string `json:"clientOrderId,omitempty"`
		TransactTime        uint64 `json:"transactTime,omitempty"`
		Price               string `json:"price,omitempty"`
		OrigQty             string `json:"origQty,omitempty"`
		ExecutedQty         string `json:"executedQty,omitempty"`
		CummulativeQuoteQty string `json:"cummulativeQuoteQty,omitempty"`
		Status              string `json:"status,omitempty"`
		TimeInForce         string `json:"timeInForce,omitempty"`
		Type                string `json:"type,omitempty"`
		Side                string `json:"side,omitempty"`
		Fills               []struct {
			Price           string `json:"price"`
			Qty             string `json:"qty"`
			Commission      string `json:"commission"`
			CommissionAsset string `json:"commissionAsset"`
			TradeId         int64  `json:"tradeId"`
		} `json:"fills,omitempty"`
		SelfTradePreventionMode string `json:"selfTradePreventionMode,omitempty"`
	} `json:"newOrderResponse,omitempty"`
	Data *struct {
		CancelResult   string `json:"cancelResult,omitempty"`
		NewOrderResult string `json:"newOrderResult,omitempty"`
		CancelResponse *struct {
			Code                    int64  `json:"code,omitempty"`
			Msg                     string `json:"msg,omitempty"`
			Symbol                  string `json:"symbol,omitempty"`
			OrigClientOrderId       string `json:"origClientOrderId,omitempty"`
			OrderId                 int64  `json:"orderId,omitempty"`
			OrderListId             int64  `json:"orderListId,omitempty"`
			ClientOrderId           string `json:"clientOrderId,omitempty"`
			Price                   string `json:"price,omitempty"`
			OrigQty                 string `json:"origQty,omitempty"`
			ExecutedQty             string `json:"executedQty,omitempty"`
			CummulativeQuoteQty     string `json:"cummulativeQuoteQty,omitempty"`
			Status                  string `json:"status,omitempty"`
			TimeInForce             string `json:"timeInForce,omitempty"`
			Type                    string `json:"type,omitempty"`
			Side                    string `json:"side,omitempty"`
			SelfTradePreventionMode string `json:"selfTradePreventionMode,omitempty"`
		} `json:"cancelResponse,omitempty"`
		NewOrderResponse struct {
			Code                    int64    `json:"code,omitempty"`
			Msg                     string   `json:"msg,omitempty"`
			Symbol                  string   `json:"symbol,omitempty"`
			OrderId                 int64    `json:"orderId,omitempty"`
			OrderListId             int64    `json:"orderListId,omitempty"`
			ClientOrderId           string   `json:"clientOrderId,omitempty"`
			TransactTime            uint64   `json:"transactTime,omitempty"`
			Price                   string   `json:"price,omitempty"`
			OrigQty                 string   `json:"origQty,omitempty"`
			ExecutedQty             string   `json:"executedQty,omitempty"`
			CummulativeQuoteQty     string   `json:"cummulativeQuoteQty,omitempty"`
			Status                  string   `json:"status,omitempty"`
			TimeInForce             string   `json:"timeInForce,omitempty"`
			Type                    string   `json:"type,omitempty"`
			Side                    string   `json:"side,omitempty"`
			Fills                   []string `json:"fills,omitempty"`
			SelfTradePreventionMode string   `json:"selfTradePreventionMode,omitempty"`
		} `json:"newOrderResponse"`
	} `json:"data,omitempty"`
}

// Query Order (USER_DATA)
// Binance Query Order (USER_DATA) (GET /api/v3/order)
// GetOrderService get order
type GetOrderService struct {
	c                 *Client
	symbol            string
	orderId           *int64
	origClientOrderId *string
}

// Symbol set symbol
func (s *GetOrderService) Symbol(symbol string) *GetOrderService {
	s.symbol = symbol
	return s
}

// OrderId set orderId
func (s *GetOrderService) OrderId(orderId int64) *GetOrderService {
	s.orderId = &orderId
	return s
}

// OrigClientOrderId set origClientOrderId
func (s *GetOrderService) OrigClientOrderId(origClientOrderId string) *GetOrderService {
	s.origClientOrderId = &origClientOrderId
	return s
}

// Do send request
func (s *GetOrderService) Do(ctx context.Context, opts ...RequestOption) (res *GetOrderResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/order",
		secType:  secTypeSigned,
	}
	m := params{
		"symbol": s.symbol,
	}
	if s.orderId != nil {
		m["orderId"] = *s.orderId
	}
	if s.origClientOrderId != nil {
		m["origClientOrderId"] = *s.origClientOrderId
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(GetOrderResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Create GetOrderResponse
type GetOrderResponse struct {
	Symbol                  string `json:"symbol"`
	OrderId                 int64  `json:"orderId"`
	OrderListId             int64  `json:"orderListId"`
	ClientOrderId           string `json:"clientOrderId"`
	Price                   string `json:"price"`
	OrigQty                 string `json:"origQty"`
	ExecutedQty             string `json:"executedQty"`
	CummulativeQuoteQty     string `json:"cummulativeQuoteQty"`
	Status                  string `json:"status"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	Side                    string `json:"side"`
	StopPrice               string `json:"stopPrice"`
	IcebergQty              string `json:"icebergQty,omitempty"`
	Time                    uint64 `json:"time"`
	UpdateTime              uint64 `json:"updateTime"`
	IsWorking               bool   `json:"isWorking"`
	WorkingTime             uint64 `json:"workingTime"`
	OrigQuoteOrderQty       string `json:"origQuoteOrderQty"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	PreventedMatchId        int64  `json:"preventedMatchId,omitempty"`
	PreventedQuantity       string `json:"preventedQuantity,omitempty"`
	StrategyId              int64  `json:"strategyId,omitempty"`
	StrategyType            int64  `json:"strategyType,omitempty"`
	TrailingDelta           string `json:"trailingDelta,omitempty"`
	TrailingTime            int64  `json:"trailingTime,omitempty"`
}
