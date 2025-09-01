package qe_connector

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

// MasterOrderMessage 主订单消息
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

// OrderMessage 订单消息
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

// FillMessage 成交消息
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

// WebSocketEventHandlers 事件处理器集合
type WebSocketEventHandlers struct {
	OnMasterOrder  func(msg *MasterOrderMessage) error
	OnOrder        func(msg *OrderMessage) error
	OnFill         func(msg *FillMessage) error
	OnStatus       func(data string) error
	OnError        func(err error)
	OnConnected    func()
	OnDisconnected func()
	OnRawMessage   func(msg *ClientPushMessage) error
}
