package qe_connector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV2TimeoutActionTransmissionAndValidation(t *testing.T) {
	for _, action := range []int32{-2, -1, 0, 1, 10, 11} {
		calls := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			got, present := body["maxTriggerWaitTimeoutAction"]
			if present != (action != -2) || (present && got != float64(action)) {
				t.Errorf("action %d serialized as %#v (present=%v)", action, got, present)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":200,"message":{"masterOrderId":"test","status":"NEW"}}`))
		}))
		service := NewClient("key", "secret", srv.URL).NewCreateMasterOrderV2Service().
			ApiKeyId("binding").Exchange("OKX").MarketType("SPOT").Symbol("BTCUSDT").
			Side("buy").Algorithm("TWAP").ExecutionDurationSeconds(60).TotalQuantity("1")
		if action != -2 {
			service.MaxTriggerWaitTimeoutAction(action)
		}
		_, err := service.Do(context.Background())
		srv.Close()
		valid := action == -2 || action >= 0 && action <= 10
		if valid && (err != nil || calls != 1) {
			t.Fatalf("action %d: %v, calls=%d", action, err, calls)
		}
		if !valid && (err == nil || calls != 0) {
			t.Fatalf("invalid action %d sent: %v, calls=%d", action, err, calls)
		}
	}
}

func TestV2TimeoutActionResponse(t *testing.T) {
	var order MasterOrderV2Info
	if err := json.Unmarshal([]byte(`{"maxTriggerWaitTimeoutAction":10}`), &order); err != nil {
		t.Fatal(err)
	}
	if order.MaxTriggerWaitTimeoutAction == nil || *order.MaxTriggerWaitTimeoutAction != 10 {
		t.Fatal("response lost action")
	}
	var old MasterOrderV2Info
	if err := json.Unmarshal([]byte(`{}`), &old); err != nil || old.MaxTriggerWaitTimeoutAction != nil {
		t.Fatal("old response incompatible")
	}
}
