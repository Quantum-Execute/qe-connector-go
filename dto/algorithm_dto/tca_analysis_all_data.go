package algorithm_dto

import "time"

// AlgorithmTCAAnalysisAllDataDTO mirrors backend DTO fields for TCA analysis full data.
// Note: This DTO is defined locally in SDK (no dependency on backend-server module).
type AlgorithmTCAAnalysisAllDataDTO struct {
	ID                  int64     `json:"id"`
	MasterOrderID       string    `json:"master_order_id"`
	ClientID            string    `json:"client_id"`
	Strategy            *string   `json:"strategy"`
	Symbol              *string   `json:"symbol"`
	Category            *string   `json:"category"`
	MidPrice            *float64  `json:"mid_price"` // ArrivalPrice
	ArrivalBid          *float64  `json:"arrival_bid"`
	ArrivalAsk          *float64  `json:"arrival_ask"`
	MoCreatedAt         time.Time `json:"mo_created_at"` // StartTime
	MoUpdatedAt         time.Time `json:"mo_updated_at"` // FinishedTime
	OrderQty            *float64  `json:"order_qty"`
	FilledQty           *float64  `json:"filled_qty"`
	Side                string    `json:"side"`
	MakeOrderNotional   *float64  `json:"make_order_notional"`
	MakeFillRate        *float64  `json:"make_fill_rate"`
	TakeOrderNotional   *float64  `json:"take_order_notional"`
	TakeFillRate        *float64  `json:"take_fill_rate"`
	Notional            *float64  `json:"notional"`
	TargetNotional      *float64  `json:"target_notional"`
	HorizonMinutes      *float64  `json:"horizon_minutes"`
	TotalExecutedQty    *float64  `json:"total_executed_qty"`
	ExecutionRate       *float64  `json:"execution_rate"`
	WeightedAvgPrice    *float64  `json:"weighted_avg_price"`
	OrdersCount         int       `json:"orders_count"`
	MatchRecordsMinTime *int64    `json:"match_records_min_time"`
	MatchRecordsMaxTime *int64    `json:"match_records_max_time"`
	ExecutionDuration   int64     `json:"execution_duration"`
	MatchedOrderIds     *string   `json:"matched_order_ids"`
	TwapBenchmark       *float64  `json:"twap_benchmark"`
	VwapBenchmark       *float64  `json:"vwap_benchmark"`
	MarketTotalNotional *float64  `json:"market_total_notional"`
	TwapSlippageBps     *float64  `json:"twap_slippage_bps"`
	VwapSlippageBps     *float64  `json:"vwap_slippage_bps"`
	ReferencePrice      *float64  `json:"reference_price"`
	PnlTheory           *float64  `json:"pnl_theory"`
	PnlRealized         *float64  `json:"pnl_realized"`
	Impact              *float64  `json:"impact"`
	UnrealizedPnl       *float64  `json:"unrealized_pnl"`
	PnlTheoryNorm       *float64  `json:"pnl_theory_norm"`
	PnlRealizedNorm     *float64  `json:"pnl_realized_norm"`
	ImpactNorm          *float64  `json:"impact_norm"`
	UnrealizedPnlNorm   *float64  `json:"unrealized_pnl_norm"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	MasterOrderTime     time.Time `json:"master_order_time"`
	Apikey              *string   `json:"apikey"`
}
