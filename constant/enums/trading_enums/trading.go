package trading_enums

type (
	MasterOrderStatus     string
	Algorithm             string
	StrategyType          string
	MarketType            string
	OrderSide             string
	MarginType            string
	Exchange              string
	Category              string
	TradingPairMarketType string
)

// 母单状态枚举
const (
	MasterOrderStatusNew       MasterOrderStatus = "NEW"       // 执行中
	MasterOrderStatusCompleted MasterOrderStatus = "COMPLETED" // 已完成
)

// 算法枚举
const (
	AlgorithmTWAP      Algorithm = "TWAP"      // TWAP算法
	AlgorithmVWAP      Algorithm = "VWAP"      // VWAP算法
	AlgorithmPOV       Algorithm = "POV"       // POV算法
	AlgorithmBoostVWAP Algorithm = "BoostVWAP" // BoostVWAP算法（高频alpha发单）
	AlgorithmBoostTWAP Algorithm = "BoostTWAP" // BoostTWAP算法（高频alpha发单）
)

// 策略类型枚举
const (
	StrategyTypeTWAP1 StrategyType = "TWAP-1" // TWAP策略版本1
	StrategyTypeTWAP2 StrategyType = "TWAP-2" // TWAP策略版本2
	StrategyTypePOV   StrategyType = "POV"    // TWAP策略版本2
)

// 市场类型枚举
const (
	MarketTypeSpot MarketType = "SPOT" // 现货市场
	MarketTypePerp MarketType = "PERP" // 合约市场
)

// 订单方向枚举
const (
	OrderSideBuy  OrderSide = "buy"  // 买入
	OrderSideSell OrderSide = "sell" // 卖出
)

// 保证金类型枚举
const (
	MarginTypeU MarginType = "U" // U本位
	MarginTypeC MarginType = "C" // 币本位
)

// 交易所枚举
const (
	ExchangeBinance Exchange = "Binance" // 币安
	ExchangeOKX     Exchange = "OKX"     // OKX
	ExchangeLTP     Exchange = "LTP"     // LTP
	ExchangeDeribit Exchange = "Deribit" // Deribit
)

// 币对品种枚举（与市场类型对应）
const (
	CategorySpot Category = "spot" // 现货品种
	CategoryPerp Category = "perp" // 合约品种
)

// 币对品种枚举（与市场类型对应）
const (
	TradingPairFutures TradingPairMarketType = "FUTURES" // 现货品种
	TradingPairSpot    TradingPairMarketType = "SPOT"    // 合约品种
)
