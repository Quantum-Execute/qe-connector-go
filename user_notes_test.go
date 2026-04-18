package qe_connector

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

func TestCreateMasterOrderServiceNotesParam(t *testing.T) {
	client := NewClient("test-api-key", "test-secret", "https://example.com")
	client.do = func(req *http.Request) (*http.Response, error) {
		if got := req.URL.Query().Get("notes"); got != "desk-a" {
			t.Fatalf("expected notes query param %q, got %q", "desk-a", got)
		}

		body := io.NopCloser(strings.NewReader(`{"code":200,"message":{"masterOrderId":"test-master-order","success":true,"message":"ok"}}`))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
			Header:     http.Header{},
		}, nil
	}

	_, err := client.NewCreateMasterOrderService().
		Algorithm(trading_enums.AlgorithmTWAP).
		Exchange(trading_enums.ExchangeBinance).
		Symbol("BTCUSDT").
		MarketType(trading_enums.MarketTypeSpot).
		Side(trading_enums.OrderSideBuy).
		ApiKeyId("test-account-id").
		OrderNotional(1000).
		StartTime("2026-04-18T09:00:00+08:00").
		ExecutionDuration(5).
		StrategyType(trading_enums.StrategyTypeTWAP1).
		Notes("desk-a").
		Do(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
