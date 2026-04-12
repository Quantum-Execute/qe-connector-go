# 更新日志（Changelog）

本文件记录 `qe-connector-go` 的用户可见变更。

## 1.1.25 - 2026-04-12

### 新增
- **Hyperliquid 支持**：`Exchange` 枚举新增 `ExchangeHyperliquid = "Hyperliquid"`
- **暂停母单接口**：新增 `PauseMasterOrderService`，支持 `PUT /user/trading/master-orders/{masterOrderId}/pause`
  - 参数：`masterOrderId`（必填）、`reason`（可选）
  - 工厂方法：`client.NewPauseMasterOrderService()`
- **恢复母单接口**：新增 `ResumeMasterOrderService`，支持 `PUT /user/trading/master-orders/{masterOrderId}/resume`
  - 参数：`masterOrderId`（必填）
  - 工厂方法：`client.NewResumeMasterOrderService()`
- **修改母单参数接口**：新增 `UpdateMasterOrderParamsService`，支持 `PUT /user/trading/master-orders/{masterOrderId}/update`
  - 必填参数：`masterOrderId`
  - 可选参数：`orderNotional`、`totalQuantity`、`upTolerance`、`lowTolerance`、`enableMake`、`makerRateLimit`、`strictUpBound`、`povLimit`、`povMinLimit`、`limitPrice`、`tailOrderProtection`、`mustComplete`、`executionDurationSeconds`
  - 工厂方法：`client.NewUpdateMasterOrderParamsService()`
- **创建母单入参**：`CreateMasterOrderService` 新增 `WorstPrice()` 方法，作为 `LimitPrice()` 的推荐替代

### 变更
- **成交记录响应**：`OrderFillInfo` 新增字段 `OrderId`（子订单ID）、`Quantity`（下单数量）、`CreatedAt`（数据创建时间）、`UpdatedAt`（最后修改时间）
- **`LimitPrice()` 标记为 Deprecated**：建议使用 `WorstPrice()` 替代，`LimitPrice()` 保留以兼容旧版本

## 1.1.24 - 2026-03-08

### 新增
- **币对品种枚举**：`Category` 新增 `CategoryPerpCm = "perp_cm"`，用于表示 Binance 币本位合约

### 变更
- **创建母单接口**：新增 Binance 币本位合约限制校验
  - 当 `exchange=Binance`、`marketType=PERP` 且 `marginType=C` 时，仅允许使用 `totalQuantity`
  - `orderNotional` 在上述场景下不可用
  - `totalQuantity` 单位为张，且输入值必须为整数

### 文档
- README：补充 `marginType` 对 `C`（币本位）的支持说明
- README：补充 `isCoin=true` 时返回币本位合约可用交易对，仅 Binance 可用
- README：更新母单响应与筛选参数中的 `category` 说明，支持 `perp_cm`
- README：补充 Binance `perp_cm` 场景下 `totalQuantity` / `orderNotional` 的使用限制

## 1.1.23 - 2026-01-29

### 新增
- **通过client_order_id获取母单详情接口**：新增 `GetMasterOrderDetailByClientOrderIdService`，支持 `GET /user/trading/master-orders/by-client-order-id/{clientOrderId}`
- **创建母单入参**：`CreateMasterOrderService` 新增 `ClientOrderId()` 方法，支持用户自定义订单ID

### 变更
- **母单数据响应**：母单列表/母单详情新增 `clientOrderId` 字段返回

### 文档
- README：补充"通过client_order_id获取母单详情"用法，并在"创建主订单"参数中增加 `clientOrderId` 说明

## 1.1.22 - 2026-01-17

### 新增
- **获取母单详情接口**：新增 `GetMasterOrderDetailService`，支持 `GET /user/trading/master-orders/{masterOrderId}`

### 变更
- **创建母单入参**：`CreateMasterOrderService` 新增 `executionDurationSeconds`（秒级执行时长；当提供且>0时优先使用，且必须大于10秒）
- **母单数据响应**：母单列表/母单详情新增 `executionDurationSeconds` 字段返回

### 文档
- README：补充“获取母单详情”用法，并在“创建主订单”参数中增加 `executionDurationSeconds` 说明

## 1.1.21 - 2026-01-07

### 文档
- **Boost 算法支持交易对更新**：更新 `CreateMasterOrderService` 中 BoostVWAP、BoostTWAP 算法支持的交易对说明（仅Binance交易所可用。）
  - 现货支持的交易对：BTCUSDT,ETHUSDT,SOLUSDT,BNBUSDT,LTCUSDT,AVAXUSDT,XLMUSDT,XRPUSDT,DOGEUSDT,CRVUSDT
  - 永续合约支持的交易对：BTCUSDT,ETHUSDT,SOLUSDT,BNBUSDT,LTCUSDT,AVAXUSDT,XLMUSDT,XRPUSDT,DOGEUSDT,ADAUSDT,BCHUSDT,FILUSDT,1000SATSUSDT,CRVUSDT

## 1.1.20 - 2026-01-06

### 变更
- **创建母单入参**：新增 `CreateMasterOrderService` 入参：`enableMake`，是否允许挂单
- **TCA 分析接口响应格式**：更新 `GetTcaAnalysisService` 返回类型
  - 从 `[]*algorithm_dto.AlgorithmTCAAnalysisAllDataDTO` 改为 `[]*algorithm_dto.TCAAnalysisResponse`
  - 响应字段名从 snake_case 改为 PascalCase，与后端接口和 Excel 表头一致
  - 新增响应 DTO：`dto/algorithm_dto/tca_analysis_response.go`
- **母单数据响应格式**：更新 `GetMasterOrdersService` 返回类型
  - 新增字段`makerRate`，被动成交率
  - 新增字段`enableMake`，是否允许挂单
  - 新增字段`tailOrderProtection`，尾单保护开关

### 文档
- **创建母单入参文档更新**：
    - 新增入参字段描述
- **TCA 分析接口文档更新**：
  - 更新响应字段描述表格，使用 Excel 表头字段名（PascalCase）
  - 更新示例代码以使用新的字段名
  - 字段描述添加中文描述
- **母单数据响应接口文档更新**：
  - 更新响应字段描述

## 1.1.20-rc.1 - 2025-12-27

### 变更
- **TCA 分析接口响应格式**：更新 `GetTcaAnalysisService` 返回类型
  - 从 `[]*algorithm_dto.AlgorithmTCAAnalysisAllDataDTO` 改为 `[]*algorithm_dto.TCAAnalysisResponse`
  - 响应字段名从 snake_case 改为 PascalCase，与后端接口和 Excel 表头一致
  - 新增响应 DTO：`dto/algorithm_dto/tca_analysis_response.go`

### 文档
- **TCA 分析接口文档更新**：
  - 更新响应字段描述表格，使用 Excel 表头字段名（PascalCase）
  - 更新示例代码以使用新的字段名
  - 字段描述直接使用 Excel 表头名称

## 1.1.19 - 2025-12-27

### 新增
- **创建母单接口**：兼容已有 `EndTime()` 方法（该字段已废弃，不再使用）

## 1.1.18 - 2025-12-26

### 新增
- **TCA 分析接口**：新增 `GetTcaAnalysisService`，支持查询 TCA（Transaction Cost Analysis）分析数据
  - 接口路径：`GET /user/trading/tca-analysis`（签名鉴权）
  - 支持参数：`symbol`、`category`、`apikey`、`startTime`、`endTime`
  - 返回类型：`[]*algorithm_dto.AlgorithmTCAAnalysisAllDataDTO`
  - 新增本地 DTO：`dto/algorithm_dto/tca_analysis_all_data.go`

### 变更
- **创建母单接口**：移除 `EndTime()` 方法和 `endTime` 字段（该字段已废弃，不再使用）

### 文档
- **参数描述更新**：
  - `makerRateLimit`：补充范围说明（包含0和1，输入0.1代表10%）
  - `povLimit`：补充范围说明（包含0和1，输入0.1代表10%）
  - `upTolerance`：更新范围说明（不包含0和1，最小输入0.0001，最大输入0.9999）
  - `lowTolerance`：更新范围说明（不包含0和1，最小输入0.0001，最大输入0.9999）

## 1.1.17 - 2025-12-15

### 新增
- **Deribit 支持**：母单创建 `exchange` 新增可选值 `Deribit`（`trading_enums.ExchangeDeribit`）。

### 变更
- **Deribit(BTCUSD/ETHUSD) 下单限制**：当 `exchange=Deribit` 且 `symbol` 为 `BTCUSD` 或 `ETHUSD` 时：
  - 仅允许使用 `totalQuantity`（单位：USD）
  - 禁止使用 `orderNotional`

### 文档
- README：`exchange` 可选值补充 `Deribit`，并补充 Deribit BTCUSD/ETHUSD 的数量字段说明。


