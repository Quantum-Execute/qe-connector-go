package qe_connector

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

func TestClient_NewListExchangeApisService(t *testing.T) {
	if os.Getenv("QE_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set QE_INTEGRATION_TEST=1 to enable")
	}
	ctx := context.Background()
	client := NewClient("your-api-key", "your-secret-key")
	do, err := client.NewListExchangeApisService().
		Exchange(trading_enums.ExchangeBinance).
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
	if os.Getenv("QE_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set QE_INTEGRATION_TEST=1 to enable")
	}
	ctx := context.Background()
	client := NewTestClient("", "")
	do, err := client.NewGetMasterOrdersService().
		Page(1).
		PageSize(10).
		Status(trading_enums.MasterOrderStatusCompleted).
		Exchange("Binance").
		Do(ctx)
	if err != nil {
		t.Errorf("err should be nil, but got %v", err)
		return
	}
	t.Logf("%#v", do)
}

func TestClient_NewGetOrderFillsService(t *testing.T) {
	if os.Getenv("QE_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set QE_INTEGRATION_TEST=1 to enable")
	}
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
	if os.Getenv("QE_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set QE_INTEGRATION_TEST=1 to enable")
	}
	ctx := context.Background()
	client := NewTestClient("", "")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// 根据提供的JSON示例创建订单
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			do, err := client.NewCreateMasterOrderService().
				MarketType(trading_enums.MarketTypeSpot).
				Symbol("BTCUSDT").
				Exchange(trading_enums.ExchangeBinance).
				Side(trading_enums.OrderSideBuy).
				StartTime(time.Now().In(loc).Format("2006-01-02T15:04:05-07:00")).
				Algorithm(trading_enums.AlgorithmTWAP).
				ExecutionDuration(5).
				ApiKeyId("6ead9160857441e3846f39b85daf55a3").
				ReduceOnly(false).
				MustComplete(true).
				OrderNotional(200).
				StrategyType(trading_enums.StrategyTypeTWAP1).
				Do(ctx)
			if err != nil {
				t.Errorf("err should be nil, but got %v", err)
				return
			}
			time.Sleep(1 * time.Second)
			client.NewCancelMasterOrderService().MasterOrderId(do.MasterOrderId).Do(ctx)
			t.Logf("%#v", do)
		}()
	}
	wg.Wait()
}

func TestClient_NewWebSocketService(t *testing.T) {
	if os.Getenv("QE_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set QE_INTEGRATION_TEST=1 to enable")
	}
	client := NewTestClient("", "")
	wsService := client.NewWebSocketService()
	// 设置事件处理器
	handlers := &WebSocketEventHandlers{
		OnConnected: func() {
			t.Logf("WebSocket connected")
		},
		OnDisconnected: func() {
			t.Logf("WebSocket disconnected")
		},
		OnError: func(err error) {
			t.Logf("WebSocket error: %v\n", err)
		},
		OnStatus: func(data string) error {
			t.Logf("Status message: %s\n", data)
			return nil
		},
		OnMasterOrder: func(msg *MasterOrderMessage) error {
			t.Logf("Master Order Update:\n")
			t.Logf("  - Master Order ID: %s\n", msg.MasterOrderID)
			t.Logf("  - Symbol: %s\n", msg.Symbol)
			t.Logf("  - Side: %s\n", msg.Side)
			t.Logf("  - Quantity: %.8f\n", msg.Qty)
			t.Logf("  - Status: %s\n", msg.Status)
			t.Logf("  - Strategy: %s\n", msg.Strategy)
			if msg.Reason != "" {
				t.Logf("  - Reason: %s\n", msg.Reason)
			}
			return nil
		},
		OnOrder: func(msg *OrderMessage) error {
			t.Logf("Order Update:\n")
			t.Logf("  - Order ID: %s\n", msg.OrderID)
			t.Logf("  - Master Order ID: %s\n", msg.MasterOrderID)
			t.Logf("  - Symbol: %s\n", msg.Symbol)
			t.Logf("  - Side: %s\n", msg.Side)
			t.Logf("  - Price: %.8f\n", msg.Price)
			t.Logf("  - Quantity: %.8f\n", msg.Quantity)
			t.Logf("  - Status: %s\n", msg.Status)
			t.Logf("  - Filled Qty: %.8f\n", msg.FillQty)
			t.Logf("  - Cumulative Filled: %.8f\n", msg.CumFilledQty)
			if msg.Reason != "" {
				t.Logf("  - Reason: %s\n", msg.Reason)
			}
			return nil
		},
		OnFill: func(msg *FillMessage) error {
			t.Logf("Fill Update:\n")
			t.Logf("  - Order ID: %s\n", msg.OrderID)
			t.Logf("  - Master Order ID: %s\n", msg.MasterOrderID)
			t.Logf("  - Symbol: %s\n", msg.Symbol)
			t.Logf("  - Side: %s\n", msg.Side)
			t.Logf("  - Fill Price: %.8f\n", msg.FillPrice)
			t.Logf("  - Filled Qty: %.8f\n", msg.FilledQty)
			t.Logf("  - Fill Time: %s\n", time.Unix(msg.FillTime/1000, 0).Format("2006-01-02 15:04:05"))
			return nil
		},
		OnRawMessage: func(msg *ClientPushMessage) error {
			// 可选：处理原始消息
			// t.Logf("Raw message - Type: %s, MessageId: %s\n", msg.Type, msg.MessageId)
			return nil
		},
	}

	wsService.SetHandlers(handlers)

	// 连接 WebSocket
	t.Logf("Connecting to WebSocket...")

	if err := wsService.Connect("ae234d8c24c14a1b8ce0546fefad199a"); err != nil {
		log.Fatalf("Failed to connect WebSocket: %v", err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	t.Logf("WebSocket client is running. Press Ctrl+C to exit.")
	t.Logf("Waiting for order updates...")
	// 等待信号
	<-sigChan

	t.Logf("\nShutting down...")

	// 关闭 WebSocket
	if err := wsService.Close(); err != nil {
		t.Logf("Error closing WebSocket: %v", err)
	}

	t.Logf("WebSocket client stopped.")
}

func TestPubFun(t *testing.T) {
	if os.Getenv("QE_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set QE_INTEGRATION_TEST=1 to enable")
	}
	ctx := context.Background()
	client := NewClient("", "")
	do, err := client.NewTradingPairsService().Do(ctx)
	if err != nil {
		t.Errorf("err should be nil, but got %v", err)
		return
	}
	t.Logf("%#v", do)
	err = client.NewPingServer().Do(ctx)
	if err != nil {
		t.Errorf("err should be nil, but got %v", err)
		return
	}
	timestampMill, err := client.NewTimestampService().Do(ctx)
	if err != nil {
		t.Errorf("err should be nil, but got %v", err)
		return
	}
	t.Logf("timestampMill: %v", timestampMill)
}
