package qe_connector

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

const tradingPairsV2Endpoint = "/pub/v2/trading-pairs"

// TradingPairsV2Service lists public trading pairs using the V2 contract.
// V2 uses SPOT/PERP/OPTION market types and returns the complete result without V1
// pagination or internal record metadata.
type TradingPairsV2Service struct {
	c          *Client
	exchange   *trading_enums.Exchange
	marketType *trading_enums.MarketType
	isCoin     *bool
	symbol     *string
}

// Exchange filters trading pairs by exchange.
func (s *TradingPairsV2Service) Exchange(exchange trading_enums.Exchange) *TradingPairsV2Service {
	s.exchange = &exchange
	return s
}

// MarketType filters trading pairs by SPOT, PERP, or OPTION (Deribit).
func (s *TradingPairsV2Service) MarketType(marketType trading_enums.MarketType) *TradingPairsV2Service {
	s.marketType = &marketType
	return s
}

// IsCoin selects coin-margined contracts when querying PERP pairs.
func (s *TradingPairsV2Service) IsCoin(isCoin bool) *TradingPairsV2Service {
	s.isCoin = &isCoin
	return s
}

// Symbol selects an exact published instrument. For an active option this also
// resolves price tiers; check RulesReady before using the returned rules.
func (s *TradingPairsV2Service) Symbol(symbol string) *TradingPairsV2Service {
	s.symbol = &symbol
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
	if s.symbol != nil {
		m["symbol"] = *s.symbol
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
	ContractSize        string            `json:"contractSize,omitempty"`
	Strike              string            `json:"strike,omitempty"`
	OptionType          string            `json:"optionType,omitempty"`
	ExpirationTimestamp *FlexInt64        `json:"expirationTimestamp,omitempty"`
	UnderlyingIndex     string            `json:"underlyingIndex,omitempty"`
	TickSize            string            `json:"tickSize,omitempty"`
	MinOrderQty         string            `json:"minOrderQty,omitempty"`
	StepSize            string            `json:"stepSize,omitempty"`
	TickSizeSteps       []OptionPriceStep `json:"tickSizeSteps,omitempty"`
	RulesReady          bool              `json:"rulesReady,omitempty"`
	Exchange            string            `json:"exchange"`
	Symbol              string            `json:"symbol"`
	BaseAsset           string            `json:"baseAsset"`
	QuoteAsset          string            `json:"quoteAsset"`
	Status              string            `json:"status"`
	MarketType          string            `json:"marketType"`
	ContractType        *string           `json:"contractType,omitempty"`
	DeliveryDate        *FlexInt64        `json:"deliveryDate,omitempty"`
}
