package qe_connector

import "testing"

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
