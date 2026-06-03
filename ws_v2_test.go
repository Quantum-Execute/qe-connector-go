package qe_connector

import (
	"encoding/json"
	"testing"
)

func TestWebSocketServiceDefaultsToV2URL(t *testing.T) {
	client := NewClient("key", "secret")
	ws := client.NewWebSocketService("wss://example.test")
	ws.listenKey = "lk"

	got := ws.getWebSocketURL()
	want := "wss://example.test/api/ws/v2?listen_key=lk"
	if got != want {
		t.Fatalf("getWebSocketURL() = %s, want %s", got, want)
	}
}

func TestWebSocketServiceCanUseLegacyV1URL(t *testing.T) {
	client := NewClient("key", "secret")
	ws := client.NewWebSocketService("wss://example.test").UseV1()
	ws.listenKey = "lk"

	got := ws.getWebSocketURL()
	want := "wss://example.test/api/ws?listen_key=lk"
	if got != want {
		t.Fatalf("getWebSocketURL() = %s, want %s", got, want)
	}
}

func TestCreateListenKeyV2ServiceUsesV2Route(t *testing.T) {
	client := NewClient("key", "secret")
	svc := client.NewCreateListenKeyV2Service()

	if svc.endpoint() != "/user/trading/v2/listen-key" {
		t.Fatalf("endpoint() = %s", svc.endpoint())
	}
}

func TestWsMasterOrderDetailDecodesV2LowerCamelCasePayload(t *testing.T) {
	payload := []byte(`{
		"masterOrderId":"DOGEUSDT-20260601-2061347360656314368",
		"clientOrderId":"go-sdk-qty-1",
		"apiKeyId":"binding-id",
		"marketType":"SPOT",
		"tradingAccount":"190pm",
		"symbol":"DOGEUSDT",
		"side":"buy",
		"totalQuantity":"100",
		"cumFilledQty":"100",
		"cumFilledNotional":"25",
		"avgFilledPrice":"0.25",
		"executionDurationSeconds":"60",
		"status":"COMPLETED",
		"commission":{"USDT":"0.01"}
	}`)

	var msg WsMasterOrderDetail
	if err := json.Unmarshal(payload, &msg); err != nil {
		t.Fatalf("decode master detail: %v", err)
	}

	if msg.MasterOrderID != "DOGEUSDT-20260601-2061347360656314368" {
		t.Fatalf("MasterOrderID = %q", msg.MasterOrderID)
	}
	if msg.ClientOrderID != "go-sdk-qty-1" {
		t.Fatalf("ClientOrderID = %q", msg.ClientOrderID)
	}
	if msg.MarketType != "SPOT" || msg.TradingAccount != "190pm" {
		t.Fatalf("unexpected account fields: %#v", msg)
	}
	if msg.TotalQuantity.String() != "100" || msg.CumFilledQty.String() != "100" {
		t.Fatalf("unexpected quantity fields: %#v", msg)
	}
	if msg.ExecutionDurationSeconds.Int64() != 60 {
		t.Fatalf("ExecutionDurationSeconds = %d", msg.ExecutionDurationSeconds.Int64())
	}
	if msg.Status != "COMPLETED" || msg.Commission["USDT"] != "0.01" {
		t.Fatalf("unexpected status/commission fields: %#v", msg)
	}
}

func TestWsOrderFillDetailDecodesV2LowerCamelCasePayloadWithStringID(t *testing.T) {
	payload := []byte(`{
		"id":"195840221",
		"orderCreatedTime":"1719999999000",
		"masterOrderId":"DOGEUSDT-20260601-2061347360656314368",
		"exchange":"BINANCE",
		"category":"spot",
		"symbol":"DOGEUSDT",
		"side":"buy",
		"filledNotional":"25",
		"filledQuantity":"100",
		"averagePrice":"0.25",
		"price":"0.25",
		"status":"FILLED",
		"rejectReason":"",
		"baseCurrency":"DOGE",
		"quoteCurrency":"USDT",
		"orderType":"LIMIT",
		"orderId":"exchange-order-1",
		"quantity":"100",
		"createdAt":"2026-06-01T00:00:00Z",
		"updatedAt":"2026-06-01T00:00:01Z"
	}`)

	var msg WsOrderFillDetail
	if err := json.Unmarshal(payload, &msg); err != nil {
		t.Fatalf("decode order fill detail: %v", err)
	}

	if msg.ID != "195840221" {
		t.Fatalf("ID = %q", msg.ID)
	}
	if msg.MasterOrderID != "DOGEUSDT-20260601-2061347360656314368" {
		t.Fatalf("MasterOrderID = %q", msg.MasterOrderID)
	}
	if msg.OrderID != "exchange-order-1" || msg.OrderType != "LIMIT" {
		t.Fatalf("unexpected order identity fields: %#v", msg)
	}
	if msg.FilledQuantity.String() != "100" || msg.FilledNotional.String() != "25" || msg.AveragePrice.String() != "0.25" {
		t.Fatalf("unexpected fill fields: %#v", msg)
	}
}
