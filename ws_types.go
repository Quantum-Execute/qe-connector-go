package qe_connector

// ClientMessageType 客户端消息类型
type ClientMessageType string
type ClientProtocolVersion string

const (
	ClientDataType            ClientMessageType = "data"
	ClientStatusType          ClientMessageType = "status"
	ClientErrorType           ClientMessageType = "error"
	ClientMasterDetailType    ClientMessageType = "master_data"
	ClientOrderFillDetailType ClientMessageType = "order_data"

	ClientProtocolV1 ClientProtocolVersion = "v1"
	ClientProtocolV2 ClientProtocolVersion = "v2"
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

// ---- 服务端 V2 推送的 DTO 结构（对应 V2 API master_data / order_data 契约） ----

// WsMasterOrderDetail 服务端通过 WS 推送的 V2 母单详情（master_data）。
// 字段与 V2 API 契约一致，使用 lowerCamelCase JSON tag。
type WsMasterOrderDetail struct {
	CreatedAt                string            `json:"createdAt"`
	UpdatedAt                string            `json:"updatedAt"`
	MasterOrderID            string            `json:"masterOrderId"`
	ClientOrderID            string            `json:"clientOrderId"`
	ApiKeyID                 string            `json:"apiKeyId"`
	TradingAccount           string            `json:"tradingAccount"`
	Exchange                 string            `json:"exchange"`
	MarketType               string            `json:"marketType"`
	Category                 string            `json:"category"`
	Symbol                   string            `json:"symbol"`
	BaseCurrency             string            `json:"baseCurrency"`
	QuoteCurrency            string            `json:"quoteCurrency"`
	Side                     string            `json:"side"`
	MarginType               string            `json:"marginType"`
	ReduceOnly               bool              `json:"reduceOnly"`
	IsMargin                 bool              `json:"isMargin"`
	Algorithm                string            `json:"algorithm"`
	TotalQuantity            FlexDecimalString `json:"totalQuantity"`
	OrderNotional            FlexDecimalString `json:"orderNotional"`
	StartTimeMs              FlexInt64         `json:"startTimeMs"`
	ExecutionDurationSeconds FlexInt64         `json:"executionDurationSeconds"`
	WorstPrice               FlexDecimalString `json:"worstPrice"`
	MustComplete             bool              `json:"mustComplete"`
	MakerRateLimit           FlexDecimalString `json:"makerRateLimit"`
	POVLimit                 FlexDecimalString `json:"povLimit"`
	POVMinLimit              FlexDecimalString `json:"povMinLimit"`
	UpTolerance              FlexDecimalString `json:"upTolerance"`
	LowTolerance             FlexDecimalString `json:"lowTolerance"`
	StrictUpBound            bool              `json:"strictUpBound"`
	TailOrderProtection      bool              `json:"tailOrderProtection"`
	EnableMake               bool              `json:"enableMake"`
	IsTargetPosition         bool              `json:"isTargetPosition"`
	Notes                    string            `json:"notes"`
	Status                   string            `json:"status"`
	RejectReason             string            `json:"rejectReason"`
	FinishedMs               FlexInt64         `json:"finishedMs"`
	CumFilledQty             FlexDecimalString `json:"cumFilledQty"`
	CumFilledNotional        FlexDecimalString `json:"cumFilledNotional"`
	AvgFilledPrice           FlexDecimalString `json:"avgFilledPrice"`
	MakerRate                FlexDecimalString `json:"makerRate"`
	CompletedQuantity        FlexDecimalString `json:"completedQuantity"`
	Commission               map[string]string `json:"commission,omitempty"`
}

// WsOrderFillDetail 服务端通过 WS 推送的 V2 子单/成交详情（order_data）。
// 字段与 V2 API 契约一致，使用 lowerCamelCase JSON tag。
type WsOrderFillDetail struct {
	ID               string            `json:"id"`
	OrderCreatedTime string            `json:"orderCreatedTime"`
	MasterOrderID    string            `json:"masterOrderId"`
	Exchange         string            `json:"exchange"`
	Category         string            `json:"category"`
	Symbol           string            `json:"symbol"`
	Side             string            `json:"side"`
	FilledNotional   FlexDecimalString `json:"filledNotional"`
	FilledQuantity   FlexDecimalString `json:"filledQuantity"`
	AveragePrice     FlexDecimalString `json:"averagePrice"`
	Price            FlexDecimalString `json:"price"`
	Status           string            `json:"status"`
	RejectReason     string            `json:"rejectReason"`
	BaseCurrency     string            `json:"baseCurrency"`
	QuoteCurrency    string            `json:"quoteCurrency"`
	OrderType        string            `json:"orderType"`
	OrderID          string            `json:"orderId"`
	Quantity         FlexDecimalString `json:"quantity"`
	CreatedAt        string            `json:"createdAt"`
	UpdatedAt        string            `json:"updatedAt"`
}

// WebSocketEventHandlers 事件处理器集合
type WebSocketEventHandlers struct {
	// 算法侧原始消息回调（保留向后兼容）
	OnMasterOrder func(msg *MasterOrderMessage) error
	OnOrder       func(msg *OrderMessage) error
	OnFill        func(msg *FillMessage) error

	// 服务端推送 DTO 回调（对应 WS 实际推送的 V2 API 字段）
	OnMasterOrderDetail func(msg *WsMasterOrderDetail) error
	OnOrderFillDetail   func(msg *WsOrderFillDetail) error

	OnStatus       func(data string) error
	OnError        func(err error)
	OnConnected    func()
	OnDisconnected func()
	OnRawMessage   func(msg *ClientPushMessage) error
}
