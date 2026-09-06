package qe_connector

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

const tradingPairsV2Endpoint = "/pub/v2/trading-pairs"

// TradingPairsV2Service lists public trading pairs using the V2 contract.
// V2 uses SPOT/PERP market types and returns the complete result without V1
// pagination or internal record metadata.
type TradingPairsV2Service struct {
	c          *Client
	exchange   *trading_enums.Exchange
	marketType *trading_enums.MarketType
	isCoin     *bool
}

// Exchange filters trading pairs by exchange.
func (s *TradingPairsV2Service) Exchange(exchange trading_enums.Exchange) *TradingPairsV2Service {
	s.exchange = &exchange
	return s
}

// MarketType filters trading pairs by SPOT or PERP.
func (s *TradingPairsV2Service) MarketType(marketType trading_enums.MarketType) *TradingPairsV2Service {
	s.marketType = &marketType
	return s
}

// IsCoin selects coin-margined contracts when querying PERP pairs.
func (s *TradingPairsV2Service) IsCoin(isCoin bool) *TradingPairsV2Service {
	s.isCoin = &isCoin
	return s
}

// Do sends GET /pub/v2/trading-pairs.
func (s *TradingPairsV2Service) Do(ctx context.Context, opts ...RequestOption) (*TradingPairsV2Reply, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: tradingPairsV2Endpoint,
		secType:  secTypeNone,
	}
	m := params{}
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
	res := new(TradingPairsV2Reply)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

// TradingPairsV2Reply is the response of GET /pub/v2/trading-pairs.
type TradingPairsV2Reply struct {
	Items []*TradingPairV2 `json:"items"`
}

// TradingPairV2 is a public V2 trading-pair record. ContractType and
// DeliveryDate are optional and may be absent when the uploaded product data
// does not provide them, including for perpetual contracts.
type TradingPairV2 struct {
	Exchange     string     `json:"exchange"`
	Symbol       string     `json:"symbol"`
	BaseAsset    string     `json:"baseAsset"`
	QuoteAsset   string     `json:"quoteAsset"`
	Status       string     `json:"status"`
	MarketType   string     `json:"marketType"`
	ContractType *string    `json:"contractType,omitempty"`
	DeliveryDate *FlexInt64 `json:"deliveryDate,omitempty"`
}
