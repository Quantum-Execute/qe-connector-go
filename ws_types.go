package qe_connector

import "time"

// ClientMessageType 客户端消息类型
type ClientMessageType string

const (
	ClientDataType            ClientMessageType = "data"
	ClientStatusType          ClientMessageType = "status"
	ClientErrorType           ClientMessageType = "error"
	ClientMasterDetailType    ClientMessageType = "master_data"
	ClientOrderFillDetailType ClientMessageType = "order_data"
)

// ClientPushMessage 客户端推送消息
type ClientPushMessage struct {
	Type      ClientMessageType `json:"type"`
	MessageId string            `json:"messageId"`
	UserId    string            `json:"userId"`
	Data      string            `json:"data"`
}

// ThirdPartyMessageType 第三方消息类型
type ThirdPartyMessageType string

const (
	MasterOrderType ThirdPartyMessageType = "master_order"
	OrderType       ThirdPartyMessageType = "order"
	FillType        ThirdPartyMessageType = "fill"
)

// MasterOrderMessage 主订单消息（算法侧原始格式，保留向后兼容）
type MasterOrderMessage struct {
	Type          ThirdPartyMessageType `json:"type"`
	MasterOrderID string                `json:"master_order_id"`
	ClientID      string                `json:"client_id"`
	Strategy      string                `json:"strategy"`
	Symbol        string                `json:"symbol"`
	Side          string                `json:"side"`
	Qty           float64               `json:"qty"`
	DurationSecs  float64               `json:"duration_secs"`
	Category      string                `json:"category"`
	Action        string                `json:"action"`
	ReduceOnly    bool                  `json:"reduce_only"`
	Status        string                `json:"status"`
	Date          float64               `json:"date"`
	TicktimeInt   int64                 `json:"ticktime_int"`
	TicktimeMs    int64                 `json:"ticktime_ms"`
	Reason        string                `json:"reason"`
	Timestamp     int64                 `json:"timestamp"`
}

// OrderMessage 订单消息（算法侧原始格式，保留向后兼容）
type OrderMessage struct {
	Type              ThirdPartyMessageType `json:"type"`
	MasterOrderID     string                `json:"master_order_id"`
	OrderID           string                `json:"order_id"`
	Symbol            string                `json:"symbol"`
	Category          string                `json:"category"`
	Side              string                `json:"side"`
	Price             float64               `json:"price"`
	Quantity          float64               `json:"quantity"`
	Status            string                `json:"status"`
	CreatedTime       int64                 `json:"created_time"`
	FillQty           float64               `json:"fill_qty"`
	FillPrice         float64               `json:"fill_price"`
	CumFilledQty      float64               `json:"cum_filled_qty"`
	QuantityRemaining float64               `json:"quantity_remaining"`
	AckTime           int64                 `json:"ack_time"`
	LastFillTime      int64                 `json:"last_fill_time"`
	CancelTime        int64                 `json:"cancel_time"`
	PriceType         string                `json:"price_type"`
	Reason            string                `json:"reason"`
	Timestamp         int64                 `json:"timestamp"`
}

// FillMessage 成交消息（算法侧原始格式，保留向后兼容）
type FillMessage struct {
	Type          ThirdPartyMessageType `json:"type"`
	MasterOrderID string                `json:"master_order_id"`
	OrderID       string                `json:"order_id"`
	Symbol        string                `json:"symbol"`
	Category      string                `json:"category"`
	Side          string                `json:"side"`
	FillPrice     float64               `json:"fill_price"`
	FilledQty     float64               `json:"filled_qty"`
	FillTime      int64                 `json:"fill_time"`
	Timestamp     int64                 `json:"timestamp"`
}

// BaseThirdPartyMessage 基础第三方消息接口
type BaseThirdPartyMessage struct {
	Type ThirdPartyMessageType `json:"type"`
}

// WebSocketHandler 处理函数类型定义
type WebSocketHandler func(data []byte) error

// ---- 服务端推送的 DTO 结构（对应 backend MasterOrderDTO / OrderFillDTO） ----

// WsMasterOrderDetail 服务端通过 WS 推送的母单详情（master_data）
// 字段与 backend-server MasterOrderDTO 一一对应，使用 snake_case JSON tag。
type WsMasterOrderDetail struct {
	ID                       uint64     `json:"id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	MasterOrderID            string     `json:"master_order_id"`
	UserUid                  uint64     `json:"user_uid"`
	Algorithm                string     `json:"algorithm"`
	Exchange                 string     `json:"exchange"`
	MarketType               string     `json:"market_type"`
	TradingAccount           string     `json:"trading_account"`
	AccountName              *string    `json:"account_name"`
	BaseCurrency             string     `json:"base_currency"`
	QuoteCurrency            string     `json:"quote_currency"`
	TradingPair              string     `json:"trading_pair"`
	Side                     string     `json:"side"`
	IsClosePosition          bool       `json:"is_close_position"`
	TotalQuantity            float64    `json:"total_quantity"`
	OrderNotional            float64    `json:"order_notional"`
	StartTime                *time.Time `json:"start_time"`
	EndTime                  *time.Time `json:"end_time"`
	ExecutionDuration        int64      `json:"execution_duration"`
	ExecutionDurationSeconds *int64     `json:"execution_duration_seconds,omitempty"`
	AlgorithmType            string     `json:"algorithm_type"`
	StrategyType             string     `json:"strategy_type"`
	Status                   string     `json:"status"`
	SubmitTime               time.Time  `json:"submit_time"`
	CompletedQuantity        float64    `json:"completed_quantity"`
	CompletionProgress       float64    `json:"completion_progress"`
	FilledAmount             float64    `json:"filled_amount"`
	TotalValue               float64    `json:"total_value"`
	AvgPrice                 float64    `json:"avg_price"`
	RejectReason             *string    `json:"reject_reason"`
	MarginType               string     `json:"margin_type"`
	ReduceOnly               bool       `json:"reduce_only"`
	LimitPrice               float64    `json:"limit_price"`
	WorstPrice               float64    `json:"worst_price"`
	MustComplete             bool       `json:"must_complete"`
	MakerRateLimit           float64    `json:"maker_rate_limit"`
	POVLimit                 float64    `json:"pov_limit"`
	POVMinLimit              float64    `json:"pov_min_limit"`
	ClientID                 *string    `json:"client_id"`
	ClientOrderID            *string    `json:"client_order_id"`
	Date                     int64      `json:"date"`
	TicktimeInt              *int64     `json:"ticktime_int"`
	TicktimeMs               *int64     `json:"ticktime_ms"`
	DurationSecs             *float64   `json:"duration_secs"`
	MidPrice                 *float64   `json:"mid_price"`
	Category                 string     `json:"category"`
	LimitPriceString         string     `json:"limit_price_string"`
	UpTolerance              string     `json:"up_tolerance"`
	LowTolerance             string     `json:"low_tolerance"`
	StrictUpBound            bool       `json:"strict_up_bound"`
	TakerMakerRate           float64    `json:"taker_maker_rate"`
	TailOrderProtection      bool       `json:"tail_order_protection"`
	IsMargin                 bool       `json:"is_margin"`
	EnableMake               bool       `json:"enable_make"`
	FinishedMs               int64      `json:"finished_ms"`
	FinishedMsSynced         bool       `json:"finished_ms_synced"`
}

// WsOrderFillDetail 服务端通过 WS 推送的子单/成交详情（order_data）
// 字段与 backend-server OrderFillDTO 一一对应，使用 snake_case JSON tag。
type WsOrderFillDetail struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	MasterOrderID    string  `json:"master_order_id"`
	BaseCurrency     string  `json:"base_currency"`
	QuoteCurrency    string  `json:"quote_currency"`
	OrderID          string  `json:"order_id"`
	Symbol           string  `json:"symbol"`
	Category         string  `json:"category"`
	Side             string  `json:"side"`
	Price            float64 `json:"price"`
	Quantity         float64 `json:"quantity"`
	Status           string  `json:"status"`
	RejectReason     *string `json:"reject_reason"`
	OrderType        string  `json:"order_type"`
	OrderCreatedTime *int64  `json:"order_created_time"`

	FilledQuantity float64 `json:"filled_quantity"`
	FilledValue    float64 `json:"filled_value"`
	AvgPrice       float64 `json:"avg_price"`

	Exchange       string  `json:"exchange"`
	TradingAccount string  `json:"trading_account"`
	AccountName    *string `json:"account_name"`

	FillPrice float64   `json:"fill_price"`
	FilledQty float64   `json:"filled_qty"`
	FillTime  time.Time `json:"fill_time"`
}

// WebSocketEventHandlers 事件处理器集合
type WebSocketEventHandlers struct {
	// 算法侧原始消息回调（保留向后兼容）
	OnMasterOrder func(msg *MasterOrderMessage) error
	OnOrder       func(msg *OrderMessage) error
	OnFill        func(msg *FillMessage) error

	// 服务端推送 DTO 回调（对应 WS 实际推送的 DB 查询结果）
	OnMasterOrderDetail func(msg *WsMasterOrderDetail) error
	OnOrderFillDetail   func(msg *WsOrderFillDetail) error

	OnStatus       func(data string) error
	OnError        func(err error)
	OnConnected    func()
	OnDisconnected func()
	OnRawMessage   func(msg *ClientPushMessage) error
}
