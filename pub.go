package qe_connector

import (
	"context"
	"encoding/json"
	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
	"net/http"
)

// TradingPairsService list trading pairs
type TradingPairsService struct {
	c          *Client
	page       *int32
	pageSize   *int32
	exchange   *trading_enums.Exchange
	marketType *trading_enums.TradingPairMarketType
	isCoin     *bool
}

// Page set page
func (s *TradingPairsService) Page(page int32) *TradingPairsService {
	s.page = &page
	return s
}

// PageSize set pageSize
func (s *TradingPairsService) PageSize(pageSize int32) *TradingPairsService {
	s.pageSize = &pageSize
	return s
}

// Exchange set exchange
func (s *TradingPairsService) Exchange(exchange trading_enums.Exchange) *TradingPairsService {
	s.exchange = &exchange
	return s
}

// MarketType set marketType
func (s *TradingPairsService) MarketType(marketType trading_enums.TradingPairMarketType) *TradingPairsService {
	s.marketType = &marketType
	return s
}

// IsCoin set isCoin
func (s *TradingPairsService) IsCoin(isCoin bool) *TradingPairsService {
	s.isCoin = &isCoin
	return s
}

// Do send request
func (s *TradingPairsService) Do(ctx context.Context, opts ...RequestOption) (res *TradingPairMessage, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/pub/trading-pairs",
		secType:  secTypeNone,
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
	if s.marketType != nil {
		m["marketType"] = *s.marketType
	}
	if s.isCoin != nil {
		m["isCoin"] = *s.isCoin
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	resp := new(TradingPairMessage)
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type TradingPairMessage struct {
	Items    []*TradingPairs `json:"items"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Total    string          `json:"total"`
}

type TradingPairs struct {
	BaseAsset    string `json:"baseAsset"`
	ContractType string `json:"contractType"`
	CreatedAt    string `json:"createdAt"`
	DeliveryDate string `json:"deliveryDate"`
	Exchange     string `json:"exchange"`
	Id           int    `json:"id"`
	MarketType   string `json:"marketType"`
	QuoteAsset   string `json:"quoteAsset"`
	Status       string `json:"status"`
	Symbol       string `json:"symbol"`
	UpdatedAt    string `json:"updatedAt"`
}
