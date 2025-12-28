package algorithm_dto

// TCAAnalysisResponse TCA分析响应结构体，字段名与后端返回格式一致（PascalCase，与Excel表头一致）
type TCAAnalysisResponse struct {
	MasterOrderID           string  `json:"MasterOrderID"`
	StartTime               string  `json:"StartTime"`
	EndTime                 string  `json:"EndTime"`
	FinishedTime            string  `json:"FinishedTime"`
	Strategy                string  `json:"Strategy"`
	Symbol                  string  `json:"Symbol"`
	Category                string  `json:"Category"`
	Side                    string  `json:"Side"`
	Date                    string  `json:"Date"`
	MasterOrderQty          float64 `json:"MasterOrderQty"`
	MasterOrderNotional     float64 `json:"MasterOrderNotional"`
	ArrivalPrice            float64 `json:"ArrivalPrice"`
	ExcutedRate             float64 `json:"ExcutedRate"`
	FillQty                 float64 `json:"FillQty"`
	TakeFillNotional        float64 `json:"TakeFillNotional"`
	MakeFillNotional        float64 `json:"MakeFillNotional"`
	FillNotional            float64 `json:"FillNotional"`
	MakerRate               float64 `json:"MakerRate"`
	ChildOrderCnt           int     `json:"ChildOrderCnt"`
	AverageFillPrice        float64 `json:"AverageFillPrice"`
	Slippage                float64 `json:"Slippage"`
	SlippagePct             float64 `json:"Slippage_pct"`
	TwapSlippagePct         float64 `json:"TWAP_Slippage_pct"`
	VwapSlippagePct         float64 `json:"VWAP_Slippage_pct"`
	Spread                  float64 `json:"Spread"`
	TwapSlippagePctFartouch float64 `json:"TWAP_Slippage_pct_Fartouch"`
	VwapSlippagePctFartouch float64 `json:"VWAP_Slippage_pct_Fartouch"`
	IntervalReturn          float64 `json:"IntervalReturn"`
	ParticipationRate       float64 `json:"ParticipationRate"`
}
