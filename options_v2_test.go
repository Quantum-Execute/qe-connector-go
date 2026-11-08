package qe_connector

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
	"github.com/Quantum-Execute/qe-connector-go/handlers"
)

func optionService(c *Client) *CreateMasterOrderV2Service {
	return c.NewCreateMasterOrderV2Service().ApiKeyId("binding-id").
		Exchange(trading_enums.ExchangeDeribit).MarketType(trading_enums.MarketTypeOption).
		Symbol("SOL_USDC-25DEC26-100-C").Side(trading_enums.OrderSideBuy).
		Algorithm(trading_enums.AlgorithmTWAP).TotalQuantity("10").WorstPrice("0.125")
}

func TestOptionCreateSignsDirectionAndPreservesUnits(t *testing.T) {
	for _, tc := range []struct{ symbol, quantity, premium string }{
		{"SOL_USDC-25DEC26-100-C", "10", "0.125"},
		{"SOL_USDC-25DEC26-100-C", "10", " 0.125 "},
		{"BTC-25DEC26-100000-P", "0.1", "0.0005"},
		{"ETH-25DEC26-3000-C", "0.1", "0.0001"},
	} {
		t.Run(tc.symbol, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				var fields map[string]interface{}
				if err := json.Unmarshal(body, &fields); err != nil {
					t.Fatal(err)
				}
				if fields["symbol"] != tc.symbol || fields["totalQuantity"] != tc.quantity || fields["worstPrice"] != tc.premium {
					t.Errorf("option units changed: %s", body)
				}
				if fields["side"] != "buy" || fields["expectedExecutionSide"] != "sell" || fields["executionDurationSeconds"] != float64(1800) {
					t.Errorf("target/duration fields: %s", body)
				}
				if fields["mustComplete"] != false || fields["tailOrderProtection"] != false {
					t.Errorf("false values lost: %s", body)
				}
				if r.URL.Query().Get("signature") != signLikeBackend(t, "secret", r.URL.RawQuery, body) {
					t.Error("signature mismatch")
				}
				if r.Method != "POST" || r.URL.Path != v2MasterOrdersEndpoint {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				io.WriteString(w, `{"code":200,"message":{"masterOrderId":"option-order","status":"NEW"}}`)
			}))
			defer srv.Close()
			client := NewClient("key", "secret")
			client.BaseURL = srv.URL
			reply, err := optionService(client).Symbol(tc.symbol).TotalQuantity(tc.quantity).WorstPrice(tc.premium).
				IsTargetPosition(true).ExpectedExecutionSide(trading_enums.OrderSideSell).
				MustComplete(false).TailOrderProtection(false).Do(context.Background())
			if err != nil || reply.MasterOrderId != "option-order" {
				t.Fatalf("create: %v / %#v", err, reply)
			}
		})
	}
}

func TestInvalidOptionsDoNotReachHTTP(t *testing.T) {
	for name, change := range map[string]func(*CreateMasterOrderV2Service){
		"exchange":            func(s *CreateMasterOrderV2Service) { s.Exchange(trading_enums.ExchangeBinance) },
		"notional":            func(s *CreateMasterOrderV2Service) { s.totalQuantity = nil; s.OrderNotional("10") },
		"margin":              func(s *CreateMasterOrderV2Service) { s.IsMargin(true) },
		"missing premium":     func(s *CreateMasterOrderV2Service) { s.worstPrice = nil },
		"zero premium":        func(s *CreateMasterOrderV2Service) { s.WorstPrice("0") },
		"unlimited premium":   func(s *CreateMasterOrderV2Service) { s.WorstPrice("-1") },
		"nan premium":         func(s *CreateMasterOrderV2Service) { s.WorstPrice("NaN") },
		"zero quantity":       func(s *CreateMasterOrderV2Service) { s.TotalQuantity("0") },
		"target preview":      func(s *CreateMasterOrderV2Service) { s.IsTargetPosition(true) },
		"target side":         func(s *CreateMasterOrderV2Service) { s.IsTargetPosition(true).ExpectedExecutionSide("long") },
		"trigger wait":        func(s *CreateMasterOrderV2Service) { s.TriggerPrice("0.1") },
		"second trigger wait": func(s *CreateMasterOrderV2Service) { s.TriggerPrice1("0.1").MaxTriggerWaitSecs(0) },
	} {
		t.Run(name, func(t *testing.T) {
			client := NewClient("key", "secret")
			client.do = func(*http.Request) (*http.Response, error) { t.Fatal("invalid option sent"); return nil, nil }
			svc := optionService(client)
			change(svc)
			if _, err := svc.Do(context.Background()); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestOptionZeroTargetPreservesDirectionConflict(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["totalQuantity"] != "0" {
			t.Errorf("zero target lost: %#v", body)
		}
		io.WriteString(w, `{"code":409,"reason":"OPTION_DIRECTION_CHANGED","message":"refresh position"}`)
	}))
	defer srv.Close()
	client := NewClient("key", "secret")
	client.BaseURL = srv.URL
	_, err := optionService(client).TotalQuantity("0").IsTargetPosition(true).
		ExpectedExecutionSide(trading_enums.OrderSideSell).Do(context.Background())
	var apiErr *handlers.APIError
	if !errors.As(err, &apiErr) || apiErr.Reason != "OPTION_DIRECTION_CHANGED" {
		t.Fatalf("conflict lost: %v", err)
	}
}

func TestOptionTriggerPricesValidateBeforeHTTP(t *testing.T) {
	for _, field := range []string{"triggerPrice", "triggerPrice1"} {
		for _, tc := range []struct {
			name  string
			price string
			wait  float64
			valid bool
		}{
			{name: "fraction", price: "1/2", wait: 120},
			{name: "zero fraction", price: "0/2"},
			{name: "hexadecimal", price: "0x1p-1", wait: 120},
			{name: "digit separator", price: "1_000", wait: 120},
			{name: "negative", price: "-0.125", wait: 120},
			{name: "not a number", price: "NaN", wait: 120},
			{name: "infinity", price: "Inf", wait: 120},
			{name: "malformed exponent", price: "1e", wait: 120},
			{name: "negative zero", price: "-0", wait: 120},
			{name: "integer overflow", price: "1e20", wait: 120},
			{name: "precision overflow", price: "1e-11", wait: 120},
			{name: "exponent overflow", price: "1e9999999999", wait: 120},
			{name: "excess decimals", price: "65000.12345678901", wait: 120},
			{name: "excess precision near maximum", price: "99999999999999999999.99999999991", wait: 120},
			{name: "disabled empty", price: "", valid: true},
			{name: "disabled whitespace", price: " \t ", valid: true},
			{name: "disabled zero", price: "0", valid: true},
			{name: "disabled decimal zero", price: "0.0", valid: true},
			{name: "disabled exponent zero", price: "0e3", valid: true},
			{name: "disabled zero with large exponent", price: "0e9999999999", valid: true},
			{name: "decimal", price: "0.125", wait: 120, valid: true},
			{name: "trimmed decimal", price: " 0.125 ", wait: 120, valid: true},
			{name: "leading decimal point", price: ".125", wait: 120, valid: true},
			{name: "trailing decimal point", price: "1.", wait: 120, valid: true},
			{name: "positive sign", price: "+0.125", wait: 120, valid: true},
			{name: "negative exponent", price: "1.25e-1", wait: 120, valid: true},
			{name: "positive exponent", price: "+1E+1", wait: 120, valid: true},
			{name: "maximum price", price: "99999999999999999999.1234567891", wait: 120, valid: true},
			{name: "minimum positive price", price: "1e-10", wait: 120, valid: true},
			{name: "maximum exponent", price: "1e19", wait: 120, valid: true},
			{name: "insignificant zeros", price: "0001.1234567891000", wait: 120, valid: true},
			{name: "shifted decimal", price: "1.234567891234e3", wait: 120, valid: true},
		} {
			t.Run(field+"/"+tc.name, func(t *testing.T) {
				client := NewClient("key", "secret")
				sent := false
				client.do = func(r *http.Request) (*http.Response, error) {
					sent = true
					var body map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Fatal(err)
					}
					if body[field] != tc.price {
						t.Errorf("trigger price changed: got %v, want %q", body[field], tc.price)
					}
					return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"code":200,"message":{"masterOrderId":"test","status":"NEW"}}`))}, nil
				}
				svc := optionService(client).MaxTriggerWaitSecs(tc.wait)
				if field == "triggerPrice" {
					svc.TriggerPrice(tc.price)
				} else {
					svc.TriggerPrice1(tc.price)
				}
				_, err := svc.Do(context.Background())
				if tc.valid && err != nil {
					t.Fatalf("valid trigger price rejected: %v", err)
				}
				if !tc.valid && err == nil {
					t.Error("invalid trigger price accepted")
				}
				if sent != tc.valid {
					t.Errorf("HTTP reached = %v, want %v", sent, tc.valid)
				}
			})
		}
	}
}

func TestOptionInstrumentQueryAndMetadata(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("symbol") != "AVAX_USDC-25DEC26-10d3-C" || r.URL.Query().Get("marketType") != "OPTION" || r.URL.Query().Get("exchange") != "Deribit" {
			t.Errorf("bad query: %s", r.URL.RawQuery)
		}
		io.WriteString(w, `{"code":200,"message":{"items":[{"symbol":"AVAX_USDC-25DEC26-10d3-C","marketType":"OPTION","contractSize":"100","strike":"10.3","optionType":"call","expirationTimestamp":"1798185600000","tickSize":"0.0001","minOrderQty":"100","stepSize":"100","rulesReady":true,"tickSizeSteps":[{"abovePrice":"0.005","tickSize":"0.0005"}]}]}}`)
	}))
	defer srv.Close()
	client := NewClient("", "")
	client.BaseURL = srv.URL
	reply, err := client.NewTradingPairsV2Service().Exchange(trading_enums.ExchangeDeribit).
		MarketType(trading_enums.MarketTypeOption).Symbol("AVAX_USDC-25DEC26-10d3-C").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	item := reply.Items[0]
	if !item.RulesReady || item.ContractSize != "100" || item.ExpirationTimestamp.Int64() != 1798185600000 || item.TickSizeSteps[0].TickSize != "0.0005" {
		t.Fatalf("rules lost: %#v", item)
	}
	for _, timestamp := range []string{`1798185600000`, `"1798185600000"`} {
		payload := []byte(`{"symbol":"BTC-25DEC26-100000-P","quoteCurrency":"BTC","optionMetadata":{"optionType":"put","strike":"100000","expirationTimestamp":` + timestamp + `,"contractSize":"1","contractType":"reversed","underlyingIndex":"btc_usd"}}`)
		var master MasterOrderV2Info
		var fill OrderFillV2Info
		var wsMaster WsMasterOrderDetail
		var wsFill WsOrderFillDetail
		for _, target := range []interface{}{&master, &fill, &wsMaster, &wsFill} {
			if err := json.Unmarshal(payload, target); err != nil {
				t.Fatal(err)
			}
		}
		for _, meta := range []*OptionMetadata{master.OptionMetadata, fill.OptionMetadata, wsMaster.OptionMetadata, wsFill.OptionMetadata} {
			if meta == nil || meta.ExpirationTimestamp.Int64() != 1798185600000 || meta.ContractType != "reversed" {
				t.Fatalf("metadata lost: %#v", meta)
			}
		}
		if master.QuoteCurrency != "BTC" || fill.QuoteCurrency != "BTC" {
			t.Fatal("premium currency lost")
		}
	}
}

func TestOptionDefaultsAndTriggersStayScoped(t *testing.T) {
	for _, market := range []trading_enums.MarketType{trading_enums.MarketTypeSpot, trading_enums.MarketTypePerp, trading_enums.MarketTypeOption} {
		t.Run(string(market), func(t *testing.T) {
			client := NewClient("key", "secret")
			client.do = func(r *http.Request) (*http.Response, error) {
				var body map[string]interface{}
				json.NewDecoder(r.Body).Decode(&body)
				wantDuration := float64(300)
				if market == trading_enums.MarketTypeOption {
					wantDuration = 1800
				}
				if body["executionDurationSeconds"] != wantDuration {
					t.Errorf("unexpected duration: %#v", body)
				}
				for _, key := range []string{"expectedExecutionSide", "mustComplete", "tailOrderProtection"} {
					if _, ok := body[key]; ok {
						t.Errorf("unset field %s was injected", key)
					}
				}
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"code":200,"message":{"masterOrderId":"test","status":"NEW"}}`))}, nil
			}
			svc := optionService(client).MarketType(market).ExecutionDurationSeconds(0)
			if market != trading_enums.MarketTypeOption {
				if err := svc.validate(); err == nil {
					t.Fatal("SPOT/PERP zero duration must still fail")
				}
				svc.ExecutionDurationSeconds(300).IsTargetPosition(true).WorstPrice("-1")
			}
			if _, err := svc.TriggerPrice("0").TriggerPrice1("").MaxTriggerWaitSecs(0).Do(context.Background()); err != nil {
				t.Fatal(err)
			}
		})
	}
	client := NewClient("", "")
	if err := optionService(client).TriggerPrice1("0.001").MaxTriggerWaitSecs(120).validate(); err != nil {
		t.Fatalf("valid option trigger rejected: %v", err)
	}
}

func TestAbsentOptionMetadataRemainsOptional(t *testing.T) {
	for _, payload := range []string{`{"symbol":"BTCUSDT"}`, `{"symbol":"BTCUSDT","optionMetadata":null}`} {
		var master MasterOrderV2Info
		var fill OrderFillV2Info
		var wsMaster WsMasterOrderDetail
		var wsFill WsOrderFillDetail
		for _, target := range []interface{}{&master, &fill, &wsMaster, &wsFill} {
			if err := json.Unmarshal([]byte(payload), target); err != nil {
				t.Fatal(err)
			}
		}
		if master.OptionMetadata != nil || fill.OptionMetadata != nil || wsMaster.OptionMetadata != nil || wsFill.OptionMetadata != nil {
			t.Fatal("absent metadata must not create an option object")
		}
	}
}
