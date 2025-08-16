package qe_connector

import (
	"context"
	"testing"
)

func TestClient_NewListExchangeApisService(t *testing.T) {
	ctx := context.Background()
	client := NewTestClient("", "")
	do, err := client.NewListExchangeApisService().
		Exchange("Binance").
		Page(1).
		PageSize(10).
		Do(ctx)
	if err != nil {
		t.Errorf("err should be nil, but got %v", err)
		return
	}
	t.Logf("%#v", do)
}

func TestClient_NewGetMasterOrdersService(t *testing.T) {
	ctx := context.Background()
	client := NewTestClient("", "")
	do, err := client.NewGetMasterOrdersService().
		Page(1).
		PageSize(10).
		Status("COMPLETED").
		Exchange("Binance").
		Do(ctx)
	if err != nil {
		t.Errorf("err should be nil, but got %v", err)
		return
	}
	t.Logf("%#v", do)
}

func TestClient_NewGetOrderFillsService(t *testing.T) {
	ctx := context.Background()
	client := NewTestClient("", "")
	do, err := client.NewGetOrderFillsService().
		Page(1).
		PageSize(10).
		Do(ctx)
	if err != nil {
		t.Errorf("err should be nil, but got %v", err)
		return
	}
	t.Logf("%#v", do)
}

func TestClient_NewCreateMasterOrderService(t *testing.T) {
	ctx := context.Background()
	client := NewTestClient("", "")

	// 根据提供的JSON示例创建订单
	do, err := client.NewCreateMasterOrderService().
		MarketType("SPOT").
		Symbol("BTCUSDT").
		Exchange("Binance").
		Side("buy").
		StartTime("2025-08-17T01:11:34+08:00").
		EndTime("2025-08-17T01:44:35+08:00").
		Algorithm("TWAP").
		AlgorithmType("TWAP").
		ExecutionDuration("5").
		ApiKeyId("").
		ReduceOnly(false).
		MustComplete(true).
		OrderNotional(200).
		StrategyType("TWAP-1").
		WorstPrice(-1).
		UpTolerance("-1").
		LowTolerance("-1").
		Do(ctx)
	if err != nil {
		t.Errorf("err should be nil, but got %v", err)
		return
	}
	t.Logf("%#v", do)
}
