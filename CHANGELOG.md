# 更新日志（Changelog）

本文件记录 `qe-connector-go` 的用户可见变更。

## 1.1.20 - 2025-12-27

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


