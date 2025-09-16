package qe_connector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketService WebSocket 服务
type WebSocketService struct {
	c              *Client
	listenKey      string
	conn           *websocket.Conn
	handlers       *WebSocketEventHandlers
	isConnected    bool
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
	reconnectDelay time.Duration
	pingInterval   time.Duration
	pongTimeout    time.Duration
	wg             sync.WaitGroup
}

// NewWebSocketService 创建 WebSocket 服务
func NewWebSocketService(c *Client) *WebSocketService {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketService{
		c:              c,
		handlers:       &WebSocketEventHandlers{},
		reconnectDelay: 5 * time.Second,
		pingInterval:   1 * time.Second,
		pongTimeout:    10 * time.Second,
		ctx:            ctx,
		cancel:         cancel,
	}
}

// SetHandlers 设置事件处理器
func (ws *WebSocketService) SetHandlers(handlers *WebSocketEventHandlers) *WebSocketService {
	ws.handlers = handlers
	return ws
}

// Connect 连接 WebSocket
func (ws *WebSocketService) Connect(listenKey string) error {
	ws.mu.Lock()
	ws.listenKey = listenKey
	ws.mu.Unlock()

	return ws.connect()
}

// connect 内部连接方法
func (ws *WebSocketService) connect() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.isConnected {
		return nil
	}

	// 构建 WebSocket URL
	wsURL := ws.getWebSocketURL()

	ws.c.debug("Connecting to WebSocket: %s", wsURL)

	// 创建 WebSocket 连接
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		ws.c.debug("Failed to connect WebSocket: %v", err)
		return fmt.Errorf("failed to connect websocket: %w", err)
	}

	ws.conn = conn
	ws.isConnected = true

	// 设置 Pong 处理器
	ws.conn.SetPongHandler(func(string) error {
		ws.conn.SetReadDeadline(time.Now().Add(ws.pongTimeout))
		return nil
	})

	// 调用连接成功回调
	if ws.handlers.OnConnected != nil {
		ws.handlers.OnConnected()
	}

	// 启动读取和心跳协程
	ws.wg.Add(1)
	go ws.readMessages()

	return nil
}

// getWebSocketURL 获取 WebSocket URL
func (ws *WebSocketService) getWebSocketURL() string {
	baseURL := "wss://test.quantumexecute.com"

	return fmt.Sprintf("%s/api/ws?listen_key=%s", baseURL, ws.listenKey)
}

// readMessages 读取消息
func (ws *WebSocketService) readMessages() {
	defer ws.wg.Done()

	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			ws.mu.Lock()
			conn := ws.conn
			ws.mu.Unlock()
			if conn == nil {
				time.Sleep(1 * time.Second)
				continue
			}
			_, message, err := conn.ReadMessage()
			if err != nil {
				ws.c.debug("WebSocket read error: %v", err)
				ws.reconnect()
				continue
			}

			ws.c.debug("Received message: %s", string(message))
			if string(message) == "pong" {
				return
			}

			// 处理消息
			go ws.handleMessage(message)
		}
	}
}

// handleMessage 处理消息
func (ws *WebSocketService) handleMessage(data []byte) {

	ws.c.debug("Received message: %s", string(data))
	// 首先解析客户端推送消息
	var clientMsg ClientPushMessage
	if err := json.Unmarshal(data, &clientMsg); err != nil {
		ws.c.debug("Failed to unmarshal client message: %v", err)
		if ws.handlers.OnError != nil {
			ws.handlers.OnError(err)
		}
		return
	}

	// 调用原始消息处理器
	if ws.handlers.OnRawMessage != nil {
		if err := ws.handlers.OnRawMessage(&clientMsg); err != nil {
			ws.c.debug("Raw message handler error: %v", err)
		}
	}

	// 根据消息类型处理
	switch clientMsg.Type {
	case ClientStatusType:
		if ws.handlers.OnStatus != nil {
			if err := ws.handlers.OnStatus(clientMsg.Data); err != nil {
				ws.c.debug("Status handler error: %v", err)
			}
		}

	case ClientErrorType:
		if ws.handlers.OnError != nil {
			ws.handlers.OnError(fmt.Errorf("server error: %s", clientMsg.Data))
		}

	case ClientMasterDetailType, ClientOrderFillDetailType:
		// 解析第三方消息类型
		var baseMsg BaseThirdPartyMessage
		if err := json.Unmarshal([]byte(clientMsg.Data), &baseMsg); err != nil {
			ws.c.debug("Failed to unmarshal base message: %v", err)
			if ws.handlers.OnError != nil {
				ws.handlers.OnError(err)
			}
			return
		}

		// 根据第三方消息类型分发
		switch baseMsg.Type {
		case MasterOrderType:
			if ws.handlers.OnMasterOrder != nil {
				var msg MasterOrderMessage
				if err := json.Unmarshal([]byte(clientMsg.Data), &msg); err != nil {
					ws.c.debug("Failed to unmarshal master order message: %v", err)
					if ws.handlers.OnError != nil {
						ws.handlers.OnError(err)
					}
					return
				}
				if err := ws.handlers.OnMasterOrder(&msg); err != nil {
					ws.c.debug("Master order handler error: %v", err)
				}
			}

		case OrderType:
			if ws.handlers.OnOrder != nil {
				var msg OrderMessage
				if err := json.Unmarshal([]byte(clientMsg.Data), &msg); err != nil {
					ws.c.debug("Failed to unmarshal order message: %v", err)
					if ws.handlers.OnError != nil {
						ws.handlers.OnError(err)
					}
					return
				}
				if err := ws.handlers.OnOrder(&msg); err != nil {
					ws.c.debug("Order handler error: %v", err)
				}
			}

		case FillType:
			if ws.handlers.OnFill != nil {
				var msg FillMessage
				if err := json.Unmarshal([]byte(clientMsg.Data), &msg); err != nil {
					ws.c.debug("Failed to unmarshal fill message: %v", err)
					if ws.handlers.OnError != nil {
						ws.handlers.OnError(err)
					}
					return
				}
				if err := ws.handlers.OnFill(&msg); err != nil {
					ws.c.debug("Fill handler error: %v", err)
				}
			}
		}
	}
}

// handleDisconnect 处理断开连接
func (ws *WebSocketService) handleDisconnect() {
	ws.mu.Lock()
	if !ws.isConnected {
		ws.mu.Unlock()
		return
	}

	ws.isConnected = false
	if ws.conn != nil {
		ws.conn.Close()
		ws.conn = nil
	}
	ws.mu.Unlock()

	// 调用断开连接回调
	if ws.handlers.OnDisconnected != nil {
		ws.handlers.OnDisconnected()
	}

	// 尝试重连
	if ws.ctx.Err() == nil {
		go ws.reconnect()
	}
}

// reconnect 重连
func (ws *WebSocketService) reconnect() {
	for {
		select {
		case <-ws.ctx.Done():
			return
		case <-time.After(ws.reconnectDelay):
			ws.c.debug("Attempting to reconnect...")
			if err := ws.connect(); err != nil {
				ws.c.debug("Reconnect failed: %v", err)
				continue
			}
			ws.c.debug("Reconnected successfully")
			return
		}
	}
}

// Close 关闭连接
func (ws *WebSocketService) Close() error {
	ws.cancel()

	ws.mu.Lock()
	if ws.conn != nil {
		err := ws.conn.Close()
		ws.conn = nil
		ws.isConnected = false
		ws.mu.Unlock()

		// 等待所有协程退出
		ws.wg.Wait()

		return err
	}
	ws.mu.Unlock()

	ws.wg.Wait()
	return nil
}

// IsConnected 是否已连接
func (ws *WebSocketService) IsConnected() bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.isConnected
}

// SetReconnectDelay 设置重连延迟
func (ws *WebSocketService) SetReconnectDelay(delay time.Duration) *WebSocketService {
	ws.reconnectDelay = delay
	return ws
}

// SetPingInterval 设置心跳间隔
func (ws *WebSocketService) SetPingInterval(interval time.Duration) *WebSocketService {
	ws.pingInterval = interval
	return ws
}

// SetPongTimeout 设置 Pong 超时时间
func (ws *WebSocketService) SetPongTimeout(timeout time.Duration) *WebSocketService {
	ws.pongTimeout = timeout
	return ws
}

// SetLogger 设置日志记录器
func (ws *WebSocketService) SetLogger(logger *log.Logger) *WebSocketService {
	ws.c.Logger = logger
	return ws
}
