package handlers

// APISuccess Success resp
type APISuccess struct {
	Code       int         `json:"code"`
	Reason     string      `json:"reason"`
	Message    interface{} `json:"message"`
	TraceId    string      `json:"traceId"`
	ServerTime int64       `json:"serverTime"`
}
