# qe-connector-go 交接文档

> 本文面向接手 `qe-connector-go` 的 SDK 维护、后端联调和发布同学。当前仓库是 Quantum Execute Strategy API 的 Go 官方 SDK，重点维护 V1 兼容层、Strategy API V2 母单交易接口、交易所账户/余额代理接口、公共交易对和 WebSocket 推送。

## 1. 项目定位

`qe-connector-go` 是外部用户在 Go 项目中调用 QE 平台 API 的客户端封装，不直接访问交易所，也不保存用户交易所 API Key。调用方传入的是 QE 平台签发的 Strategy API Key / Secret，SDK 负责：

- 生成 `timestamp`、`signature`，设置 `X-MBX-APIKEY`。
- 封装 `/user/...`、`/pub/...`、`/ping`、`/timestamp` 等 HTTP API。
- 兼容历史 V1 service，并新增 V2 service builder。
- 对 V2 JSON body 写接口按后端签名规则重新组装签名串。
- 连接平台 WebSocket，消费母单和子单/成交推送。

和其它项目的关系：

- 后端 V2 契约参考 [backend-server/docs/frontend-v2-api-upgrade.md](../../backend-server/docs/frontend-v2-api-upgrade.md)。
- WebSocket 服务端行为参考 [backend-server/docs/websocket-api.md](../../backend-server/docs/websocket-api.md)。
- 外部文档站交接参考 [QE-API-Key/docs/qe-api-key-handoff.md](../../QE-API-Key/docs/qe-api-key-handoff.md)。
- Python SDK 交接参考 [qe-connector-python/docs/qe-connector-python-handoff.md](../../qe-connector-python/docs/qe-connector-python-handoff.md)。

需要特别注意：

- SDK endpoint 常量写的是 `/user/...`、`/pub/...`，没有硬编码 `/strategy-api`。线上是否带 `/strategy-api` 前缀由 baseURL 或前置网关决定。
- `apiKeyId` 是用户绑定的交易所账户 ID；`X-MBX-APIKEY` 是 Strategy API Key，二者不是同一个概念。
- 算法侧如果只有 `ApiKeyToken`，应先通过后端 `/task/get-api-key` 换取/解密 Strategy API Key 与 Secret，再初始化 SDK；SDK 本身不调用该接口。

## 2. 快速启动

当前模块信息：

| 项 | 当前值 |
| --- | --- |
| module | `github.com/Quantum-Execute/qe-connector-go` |
| Go 版本 | `go 1.24` |
| SDK 版本 | `consts.go` 中 `Version = "1.3.1"` |
| 主要依赖 | `github.com/gorilla/websocket`、`github.com/bitly/go-simplejson` |

本地验证：

```bash
go test ./...
```

聚焦 V2、交易所枚举和 WebSocket：

```bash
go test ./... -run 'TestV2|TestWebSocket|TestExchange|TestBitget'
```

基础使用：

```go
client := qe_connector.NewClient(apiKey, secretKey)

res, err := client.NewCreateMasterOrderV2Service().
    ApiKeyId("exchange-api-binding-id").
    Exchange(trading_enums.ExchangeBinance).
    MarketType(trading_enums.MarketTypePerp).
    Symbol("BTCUSDT").
    Side(trading_enums.OrderSideBuy).
    Algorithm(trading_enums.AlgorithmTWAP).
    ExecutionDurationSeconds(3600).
    TotalQuantity("0.1").
    WorstPrice("90000").
    Do(context.Background())
```

如果接入测试环境：

```go
client := qe_connector.NewTestClient(apiKey, secretKey)
```

如需显式指定网关前缀：

```go
client := qe_connector.NewClient(apiKey, secretKey, "https://api.quantumexecute.com/strategy-api")
```

是否需要带 `/strategy-api` 取决于部署层转发规则；不要在 endpoint 常量和 baseURL 两边重复添加。

## 3. 目录结构

| 路径 | 说明 |
| --- | --- |
| `client.go` | `Client`、HTTP client、V1/V2 service 构造器、响应拆包 |
| `request.go` | V1 query/form request、`WithRecvWindow` |
| `user.go` | V1 用户交易接口：母单、子单、TCA、listenKey 等 |
| `user_v2.go` | V2 用户交易接口、V2 类型、JSON body 签名 |
| `pub.go` | V1 公共交易对 |
| `pub_v2.go` | V2 公共交易对 `GET /pub/v2/trading-pairs` |
| `exchange_balance.go` | 交易所余额、账户、持仓代理接口 |
| `ws_client.go` | WebSocket 连接、重连、心跳、消息分发 |
| `ws_types.go` | WebSocket 推送 DTO 和 handler 定义 |
| `flex_types.go` | 兼容 string/number 的弹性 JSON 类型 |
| `constant/enums/trading_enums/trading.go` | 交易、市场、交易所枚举 |
| `handlers/errors.go`、`handlers/success.go` | Kratos 响应 wrapper |
| `*_test.go` | V1/V2 签名、枚举、Bitget、WS 测试 |

## 4. Client、签名与响应解析

`NewClient(apiKey, secretKey, baseURL...)` 默认使用 `https://api.quantumexecute.com`，`NewTestClient(...)` 默认使用 `https://testapi.quantumexecute.com` 并开启 `Debug`。

HTTP client 做了连接复用优化：

- `Timeout: 30s`
- `MaxIdleConns: 200`
- `MaxIdleConnsPerHost: 64`
- `IdleConnTimeout: 90s`

V1/GET 签名流程在 `parseRequest`：

1. 应用 `RequestOption`，例如 `WithRecvWindow(5000)`。
2. SIGNED 请求自动加入 `timestamp=currentTimestamp()-TimeOffset`。
3. query 与 form 按 `url.Values.Encode()` 编码。
4. 对 `queryString + bodyString` 做 HMAC-SHA256。
5. 设置 `X-MBX-APIKEY` 和 `User-Agent: qe-connector-go/{Version}`。

响应解析在 `callAPI`：

- HTTP status `>= 400` 时按 `handlers.APIError` 解析并返回错误。
- HTTP 2xx 时按 Kratos wrapper 解析。
- `code != 200` 也会转换为 `handlers.APIError`。
- `code == 200` 时只返回 `message` 字段给上层 service，再由 service 反序列化到具体 reply struct。

这个拆包行为必须和 Python SDK、文档站示例保持一致。

## 5. Strategy API V2

V2 endpoint 常量集中在 `user_v2.go`：

| Service | Method / Path |
| --- | --- |
| `ListExchangeApisV2Service` | `GET /user/exchange/v2/exchange-apis` |
| `CreateMasterOrderV2Service` | `POST /user/trading/v2/master-orders` |
| `GetMasterOrdersV2Service` | `GET /user/trading/v2/master-orders` |
| `GetMasterOrderDetailV2Service` | `GET /user/trading/v2/master-orders/{masterOrderId}` |
| `GetMasterOrderDetailByClientOrderIdV2Service` | `GET /user/trading/v2/master-orders/by-client-order-id/{clientOrderId}` |
| `GetOrderFillsV2Service` | `GET /user/trading/v2/order-fills` |
| `GetTCAAnalysisV2Service` | `GET /user/trading/v2/tca-analysis` |
| `CreateListenKeyV2Service` | `POST /user/trading/v2/listen-key` |
| `CancelMasterOrderV2Service` | `PUT /user/trading/v2/master-orders/{masterOrderId}/cancel` |
| `PauseMasterOrderV2Service` | `PUT /user/trading/v2/master-orders/{masterOrderId}/pause` |
| `ResumeMasterOrderV2Service` | `PUT /user/trading/v2/master-orders/{masterOrderId}/resume` |
| `UpdateMasterOrderParamsV2Service` | `PUT /user/trading/v2/master-orders/{masterOrderId}/update` |
| `BatchCancelMasterOrdersV2Service` | `PUT /user/trading/v2/master-orders/batch-cancel` |

公共数据 V2：

| Service | Method / Path | 说明 |
| --- | --- | --- |
| `TradingPairsV2Service` | `GET /pub/v2/trading-pairs` | 支持 `exchange`、`marketType=SPOT/PERP`、`isCoin`，响应无 V1 分页字段 |

V2 创建母单的核心校验：

- 必填：`apiKeyId`、`exchange`、`marketType`、`symbol`、`side`、`algorithm`、`executionDurationSeconds`。
- `executionDurationSeconds` 必须大于 10。
- `totalQuantity` 与 `orderNotional` 必须二选一。
- `isTargetPosition=true` 时必须传 `totalQuantity`，且禁止 `orderNotional`。
- Deribit `BTCUSD` / `ETHUSD` notional 场景受限，沿用后端规则。
- Binance `PERP` 且 `marginType=C` 时必须使用 `totalQuantity`，禁止 `orderNotional`。
- V2 推荐使用 `worstPrice`，不要再新增依赖 `limitPrice` / `limitPriceString`。
- `povLimit` 默认值由当前代码决定：`POV -> "0.05"`，其它算法 -> `"1"`。

V2 列表查询注意：

- `pageSize` 上限为 100，SDK 侧直接拒绝更大的值。
- 母单状态全集包含 `NEW`、`WAITING`、`PROCESSING`、`PAUSED`、`CANCELLED`、`COMPLETED`、`COMPLETED_WITHTAIL`、`REJECTED`、`EXPIRED`。
- 列表过滤只推荐传聚合值：`NEW` 表示运行中，`COMPLETED` 表示非运行中。
- 详情和 WebSocket 推送会返回细分状态。

兼容字段：

- `ExchangeApiV2Info` 中 `apiKeyId` 是新字段。
- `apiKeyUuid`、`id` 是为旧 SDK 调用方保留的 deprecated alias，会从 `apiKeyId` 回填。
- `FlexInt64`、`FlexDecimalString` 用于兼容后端偶发返回 number/string 两种 JSON 形态。

## 6. V2 JSON Body 签名

带业务 JSON body 的 V2 写接口使用 `callAPIV2WithJSONBody`，这是接手时最容易踩坑的地方。没有业务 body 的接口，例如 listenKey 创建，可继续走普通 signed request。

请求形态：

- `timestamp`、`recvWindow`、`signature` 放 URL query。
- 业务字段放 JSON body。
- Header 设置 `Content-Type: application/json` 和 `X-MBX-APIKEY`。

签名串规则：

1. JSON body 先 `json.Marshal`，再用 `Decoder.UseNumber()` 反解成 `map[string]interface{}`。
2. 合并 query 字段和 JSON body 顶层字段。
3. 跳过 body 中的 `timestamp`、`signature`。
4. 标量按后端 `scalarToString` 对齐：string 原样、number 原样、bool 为 `true/false`。
5. 数组/对象用紧凑 JSON 字符串参与签名，例如 `masterOrderIds=["mo_1","mo_2"]`。
6. 最终对 `url.Values.Encode()` 的结果做 HMAC-SHA256。

新增 V2 POST/PUT endpoint 时，优先复用 `callAPIV2WithJSONBody`。如果误走 V1 `callAPI`，业务字段会进入 query/form，容易和 gin 签名中间件的 body+query 规则不一致。

## 7. 已接入交易所与 SDK 覆盖

交易所枚举在 `constant/enums/trading_enums/trading.go`：

| 交易所 | 枚举 | 当前 SDK 覆盖 |
| --- | --- | --- |
| Binance | `ExchangeBinance` | V1/V2 母单、公共交易对、spot/futures/PAPI/PV1/UM/CM/DAPI/cross margin 余额与账户、position side |
| OKX | `ExchangeOKX` | V1/V2 母单、公共交易对、账户余额、持仓、最大可下单量 |
| LTP | `ExchangeLTP` | V1/V2 母单、公共交易对、账户、组合资产、持仓 |
| Deribit | `ExchangeDeribit` | V1/V2 母单、公共交易对、账户、持仓 |
| Hyperliquid | `ExchangeHyperliquid` | V1/V2 母单、公共交易对、spot balance、perp positions |
| Bybit | `ExchangeBybit` | V2 交易所校验、母单、公共交易对；目标仓位模式遵循 V2 数量规则 |
| Bitget | `ExchangeBitget` | V1/V2 母单、公共交易对；币本位 symbol 形如 `BTCUSD_CM` |

余额/账户/持仓接口统一在 `exchange_balance.go`，目前是平台后端代理到交易所，SDK 只做 QE API 签名。常用 service：

| 分类 | Service |
| --- | --- |
| Binance 余额 | `NewGetAccountBalanceService`、`NewGetMarginBalanceService`、`NewGetPv1BalanceService` |
| Binance 账户 | `NewGetUmAccountService`、`NewGetCmAccountService`、`NewGetPv1AccountService`、`NewGetDapiAccountService`、`NewGetFapiAccountService`、`NewGetCrossMarginAccountDetailService` |
| Binance 持仓设置 | `NewGetFapiPositionSideDialService`、`NewGetPapiUmPositionSideDualService` |
| OKX | `NewGetOkxAccountBalanceService`、`NewGetOkxAccountPositionsService`、`NewGetOkxAccountMaxSizeService` |
| LTP | `NewGetLtpAccountService`、`NewGetLtpPortfolioAssetService`、`NewGetLtpPositionService` |
| Deribit | `NewGetDeribitAccountService`、`NewGetDeribitPositionService` |
| Hyperliquid | `NewGetHyperliquidSpotBalanceService`、`NewGetHyperliquidPositionsService` |

如果后端新增交易所：

1. 在 `trading_enums.Exchange` 增加枚举。
2. 同步 V2 交易所白名单和创建母单校验。
3. 如有公共交易对，确认 `TradingPairsV2Service` 不需要特殊序列化。
4. 如有余额/持仓代理接口，在 `exchange_balance.go` 增加 service 和 response struct。
5. 补充 `exchange_enum_test.go`、交易所专项测试和 README 示例。
6. 同步 Python SDK 与 QE-API-Key 文档站。

## 8. WebSocket

Go SDK 的 WebSocket 入口：

```go
ws := client.NewWebSocketService("wss://www.quantumexecute.com")
ws.SetHandlers(&qe_connector.WebSocketEventHandlers{
    OnConnected: func() {},
    OnMasterOrderDetail: func(order *qe_connector.WsMasterOrderDetail) error { return nil },
    OnOrderFillDetail: func(fill *qe_connector.WsOrderFillDetail) error { return nil },
})
err := ws.Connect(listenKey)
```

默认行为：

- 默认协议版本是 V2。
- 默认 host 是 `wss://test.quantumexecute.com`，可通过 `NewWebSocketService(host...)` 或 `SetHost` 覆盖。
- V2 路径是 `/api/ws/v2?listen_key=...`。
- 调 `UseV1()` 后路径变为 `/api/ws?listen_key=...`。
- 断线后会重连，默认 `reconnectDelay=5s`。
- ping 间隔 `1s`，pong timeout `10s`。

推送分发：

- `status` 走 `OnStatus`。
- `error` 走 `OnError`。
- `master_data` 优先解析成 `WsMasterOrderDetail` 并走 `OnMasterOrderDetail`。
- `order_data` 优先解析成 `WsOrderFillDetail` 并走 `OnOrderFillDetail`。
- `OnRawMessage` 可用于调试完整 envelope。

转发注意：

- 前端或平台 API 调用通常走 `/api/ws/v2`，由 Nginx/gin 转发到后端 WebSocket。
- SDK 直接连接时 host 必须是 `wss://...`，不要传 HTTP baseURL。
- listenKey 来自 V2 `CreateListenKeyV2Service` 或 V1 `CreateListenKeyService`，必须和后端用户/Strategy API Key 权限匹配。

## 9. 测试与质量门槛

文档或小范围 endpoint 改动至少跑：

```bash
go test ./...
git diff --check
```

针对性测试：

| 场景 | 建议命令 |
| --- | --- |
| V2 JSON 签名、分页、字段 | `go test ./... -run TestV2` |
| WebSocket V2 | `go test ./... -run TestWebSocket` |
| 交易所枚举 | `go test ./... -run TestExchange` |
| Bitget | `go test ./... -run TestBitget` |

新增测试时优先使用 fake `Client.do` 拦截请求，断言 URL、query、header、body 和签名，不依赖真实网络。

## 10. 新增接口维护 Runbook

新增 GET endpoint：

1. 在对应文件增加 endpoint 常量。
2. 增加 service struct、链式 setter、`Do(ctx, opts...)`。
3. `Do` 内构建 `request{method, endpoint, secType}`，SIGNED 接口使用 `secTypeSigned`。
4. 响应 struct 使用后端 JSON 字段名。
5. 在 `client.go` 增加 `New...Service()` 构造器。
6. 补测试：query 参数、签名头、响应拆包、错误响应。

新增带 JSON body 的 V2 POST/PUT endpoint：

1. 使用 JSON body 的接口走 `callAPIV2WithJSONBody`。
2. body 字段保持 lowerCamelCase，与后端 wire format 一致。
3. Decimal 字段用 string，对外 setter 尽量接收 string，避免浮点精度。
4. 数组/对象参与签名时确认测试覆盖紧凑 JSON。
5. 同步 Python SDK、QE-API-Key 示例和后端字段词典。

新增字段：

- 回包字段优先按后端原名建 struct tag。
- 若后端可能返回 number/string 两种形态，复用或扩展 `FlexDecimalString`、`FlexInt64`。
- 如果字段替代旧字段，旧字段标记 deprecated，但不要在小版本中删除。

## 11. 发版流程

Go SDK 没有构建产物需要提交，核心是源码、README/CHANGELOG 和 git tag。

推荐流程：

1. 更新 `consts.go` 的 `Version`。
2. 更新 `CHANGELOG.md`，把 `Unreleased` 内容归档到新版本。
3. 如外部调用方式变化，同步 README 示例和 QE-API-Key 文档站示例。
4. 跑 `go test ./...`。
5. 确认 `git diff --check` 无空白错误。
6. 合并后打 tag，例如 `v1.3.2`，让 Go module 生态可拉取。

发布前必须确认：

- 没有把真实 API Key、Secret、listenKey、用户 bindingId 写进测试或文档。
- 没有误提交 IDE、本地缓存或临时文件。
- Go 与 Python SDK 对相同 endpoint 的字段名、默认值、校验规则一致。

## 12. 常见风险点

- baseURL 与 `/strategy-api` 前缀重复或缺失，会导致 404 或签名验证路径混乱。
- `apiKeyId` 和 Strategy API Key 容易混淆；下单 body 里传的是交易所绑定 ID。
- V2 写接口如果误把业务字段放 query，可能无法通过后端签名验证。
- `pageSize > 100` 在 V2 会被拒绝，不是静默裁剪。
- 列表 status 过滤只用 `NEW` / `COMPLETED` 聚合值；细分状态主要用于详情和推送展示。
- Python SDK 的方法名是 snake_case，Go SDK 是 service builder；文档站示例需要分别贴真实写法。
