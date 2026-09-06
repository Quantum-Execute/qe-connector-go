package qe_connector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

func TestCreateMasterOrderV1SendsBitgetPayload(t *testing.T) {
	var queryValues map[string][]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/user/trading/master-orders" {
			t.Errorf("path = %s, want /user/trading/master-orders", r.URL.Path)
		}
		queryValues = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":{"masterOrderId":"mo_bitget_v1","success":true,"message":"accepted"}}`))
	}))
	defer srv.Close()

	client := NewClient("api-key", "secret", srv.URL)
	reply, err := client.NewCreateMasterOrderService().
		ApiKeyId("bitget-binding-id").
		Exchange(trading_enums.ExchangeBitget).
		MarketType(trading_enums.MarketTypePerp).
		Symbol("BTCUSD_CM").
		Side(trading_enums.OrderSideBuy).
		Algorithm(trading_enums.AlgorithmTWAP).
		ExecutionDurationSeconds(600).
		TotalQuantity(2).
		MarginType(trading_enums.MarginTypeC).
		IsTargetPosition(false).
		IsMargin(false).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if reply.MasterOrderId != "mo_bitget_v1" || !reply.Success {
		t.Fatalf("unexpected reply: %#v", reply)
	}

	want := map[string]string{
		"apiKeyId":                 "bitget-binding-id",
		"exchange":                 "Bitget",
		"marketType":               "PERP",
		"symbol":                   "BTCUSD_CM",
		"side":                     "buy",
		"algorithm":                "TWAP",
		"executionDurationSeconds": "600",
		"totalQuantity":            "2",
		"marginType":               "C",
		"isTargetPosition":         "false",
		"isMargin":                 "false",
	}
	for key, wantValue := range want {
		if got := firstQueryValue(queryValues, key); got != wantValue {
			t.Errorf("query[%q] = %q, want %q", key, got, wantValue)
		}
	}
	if _, ok := queryValues["VenueCategory"]; ok {
		t.Fatal("V1 payload must not contain VenueCategory")
	}
}

func TestCreateMasterOrderV2SendsBitgetMarginAndTargetFields(t *testing.T) {
	tests := []struct {
		name             string
		marketType       trading_enums.MarketType
		symbol           string
		marginType       trading_enums.MarginType
		totalQuantity    string
		orderNotional    string
		isMargin         bool
		isTargetPosition bool
	}{
		{
			name:          "spot margin order",
			marketType:    trading_enums.MarketTypeSpot,
			symbol:        "BTCUSDT",
			orderNotional: "100",
			isMargin:      true,
		},
		{
			name:             "U-margined target position",
			marketType:       trading_enums.MarketTypePerp,
			symbol:           "BTCUSDT",
			marginType:       trading_enums.MarginTypeU,
			totalQuantity:    "0",
			isTargetPosition: true,
		},
		{
			name:             "coin-margined normal order",
			marketType:       trading_enums.MarketTypePerp,
			symbol:           "BTCUSD_CM",
			marginType:       trading_enums.MarginTypeC,
			totalQuantity:    "2",
			isTargetPosition: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body map[string]interface{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %s, want POST", r.Method)
				}
				if r.URL.Path != "/user/trading/v2/master-orders" {
					t.Errorf("path = %s, want /user/trading/v2/master-orders", r.URL.Path)
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode request body: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"code":200,"message":{"masterOrderId":"mo_bitget_v2","status":"NEW","clientOrderId":"client-bitget"}}`))
			}))
			defer srv.Close()

			client := NewClient("api-key", "secret", srv.URL)
			service := client.NewCreateMasterOrderV2Service().
				ApiKeyId("bitget-binding-id").
				Exchange(trading_enums.ExchangeBitget).
				MarketType(tc.marketType).
				Symbol(tc.symbol).
				Side(trading_enums.OrderSideBuy).
				Algorithm(trading_enums.AlgorithmTWAP).
				ExecutionDurationSeconds(600).
				IsTargetPosition(tc.isTargetPosition).
				IsMargin(tc.isMargin).
				PovLimit("0.8").
				ClientOrderId("client-bitget")
			if tc.totalQuantity != "" {
				service.TotalQuantity(tc.totalQuantity)
			}
			if tc.orderNotional != "" {
				service.OrderNotional(tc.orderNotional)
			}
			if tc.marginType != "" {
				service.MarginType(tc.marginType)
			}
			reply, err := service.Do(context.Background())
			if err != nil {
				t.Fatalf("Do() error = %v", err)
			}
			if reply.MasterOrderId != "mo_bitget_v2" {
				t.Fatalf("unexpected reply: %#v", reply)
			}

			if got := body["exchange"]; got != "Bitget" {
				t.Errorf("exchange = %#v, want Bitget", got)
			}
			if got := body["marketType"]; got != string(tc.marketType) {
				t.Errorf("marketType = %#v, want %q", got, tc.marketType)
			}
			if tc.marginType == "" {
				if _, ok := body["marginType"]; ok {
					t.Errorf("marginType must be omitted for SPOT: %#v", body["marginType"])
				}
			} else if got := body["marginType"]; got != string(tc.marginType) {
				t.Errorf("marginType = %#v, want %q", got, tc.marginType)
			}
			if tc.totalQuantity != "" {
				if got := body["totalQuantity"]; got != tc.totalQuantity {
					t.Errorf("totalQuantity = %#v, want %q", got, tc.totalQuantity)
				}
				if _, ok := body["orderNotional"]; ok {
					t.Fatal("orderNotional must remain omitted when totalQuantity is used")
				}
			} else {
				if got := body["orderNotional"]; got != tc.orderNotional {
					t.Errorf("orderNotional = %#v, want %q", got, tc.orderNotional)
				}
				if _, ok := body["totalQuantity"]; ok {
					t.Fatal("totalQuantity must remain omitted when orderNotional is used")
				}
			}
			if got := body["isTargetPosition"]; got != tc.isTargetPosition {
				t.Errorf("isTargetPosition = %#v, want %t", got, tc.isTargetPosition)
			}
			if got := body["isMargin"]; got != tc.isMargin {
				t.Errorf("isMargin = %#v, want %t", got, tc.isMargin)
			}
			if _, ok := body["VenueCategory"]; ok {
				t.Fatal("V2 payload must not contain VenueCategory")
			}
		})
	}
}

func TestTradingPairsV1SupportsBitget(t *testing.T) {
	var queryValues map[string][]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/pub/trading-pairs" {
			t.Errorf("path = %s, want /pub/trading-pairs", r.URL.Path)
		}
		queryValues = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":{"items":[{"id":7,"exchange":"Bitget","symbol":"BTCUSD_CM","baseAsset":"BTC","quoteAsset":"USD","status":"TRADING","marketType":"PERP"}],"total":"1","page":1,"pageSize":20}}`))
	}))
	defer srv.Close()

	client := NewClient("", "", srv.URL)
	reply, err := client.NewTradingPairsService().
		Exchange(trading_enums.ExchangeBitget).
		MarketType(trading_enums.TradingPairPerp).
		IsCoin(true).
		Page(1).
		PageSize(20).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	wantQuery := map[string]string{
		"exchange":   "Bitget",
		"marketType": "PERP",
		"isCoin":     "true",
		"page":       "1",
		"pageSize":   "20",
	}
	for key, wantValue := range wantQuery {
		if got := firstQueryValue(queryValues, key); got != wantValue {
			t.Errorf("query[%q] = %q, want %q", key, got, wantValue)
		}
	}
	if reply.Total != "1" || len(reply.Items) != 1 {
		t.Fatalf("unexpected reply: %#v", reply)
	}
	if pair := reply.Items[0]; pair.Exchange != "Bitget" || pair.Symbol != "BTCUSD_CM" || pair.MarketType != "PERP" {
		t.Fatalf("unexpected pair: %#v", pair)
	}
}

func TestTradingPairsV2SupportsBitgetPathQueryAndResponse(t *testing.T) {
	var queryValues map[string][]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/pub/v2/trading-pairs" {
			t.Errorf("path = %s, want /pub/v2/trading-pairs", r.URL.Path)
		}
		queryValues = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":{"items":[{"exchange":"Bitget","symbol":"BTCUSD_CM","baseAsset":"BTC","quoteAsset":"USD","status":"TRADING","marketType":"PERP"}]}}`))
	}))
	defer srv.Close()

	client := NewClient("", "", srv.URL)
	reply, err := client.NewTradingPairsV2Service().
		Exchange(trading_enums.ExchangeBitget).
		MarketType(trading_enums.MarketTypePerp).
		IsCoin(true).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	wantQuery := map[string]string{
		"exchange":   "Bitget",
		"marketType": "PERP",
		"isCoin":     "true",
	}
	for key, wantValue := range wantQuery {
		if got := firstQueryValue(queryValues, key); got != wantValue {
			t.Errorf("query[%q] = %q, want %q", key, got, wantValue)
		}
	}
	if len(reply.Items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(reply.Items))
	}
	pair := reply.Items[0]
	if pair.Exchange != "Bitget" || pair.Symbol != "BTCUSD_CM" || pair.BaseAsset != "BTC" ||
		pair.QuoteAsset != "USD" || pair.Status != "TRADING" || pair.MarketType != "PERP" {
		t.Fatalf("unexpected pair: %#v", pair)
	}
	if pair.ContractType != nil {
		t.Fatalf("contractType = %#v, want nil when omitted", pair.ContractType)
	}
	if pair.DeliveryDate != nil {
		t.Fatalf("deliveryDate = %#v, want nil when omitted", pair.DeliveryDate)
	}
}

func firstQueryValue(values map[string][]string, key string) string {
	if len(values[key]) == 0 {
		return ""
	}
	return values[key][0]
}
