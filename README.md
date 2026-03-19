# Quantum Execute Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/Quantum-Execute/qe-connector-go.svg)](https://pkg.go.dev/github.com/Quantum-Execute/qe-connector-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

这是 Quantum Execute 公共 API 的官方 Go SDK，为开发者提供了一个轻量级、易于使用的接口来访问 Quantum Execute 的交易服务。

## 功能特性

- ✅ 完整的 Quantum Execute API 支持
- ✅ 交易所 API 密钥管理
- ✅ 主订单创建与管理（TWAP、VWAP 等算法）
- ✅ 订单查询和成交明细
- ✅ ListenKey 创建与管理
- ✅ 交易对信息查询
- ✅ 服务器连通性测试
- ✅ 服务器时间同步
- ✅ 安全的 HMAC-SHA256 签名认证
- ✅ 支持生产环境和测试环境
- ✅ 链式调用 API 设计
- ✅ 完整的错误处理

## 安装

```bash
go get github.com/Quantum-Execute/qe-connector-go
```

## 快速开始

### 初始化客户端

```go
package main

import (
    "context"
    "log"
    qe "github.com/Quantum-Execute/qe-connector-go"
)

func main() {
    // 创建生产环境客户端
    client := qe.NewClient("your-api-key", "your-api-secret")
    
    // 或创建测试环境客户端
    // client := qe.NewTestClient("your-api-key", "your-api-secret")
    
    // 启用调试模式
    client.Debug = true
}
```

## API 参考

### 公共接口

#### 服务器连通性测试

##### Ping 服务器

测试与 Quantum Execute 服务器的连通性。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| 无需参数 | - | - | - |

**响应：**

成功时无返回内容，失败时返回错误信息。

**示例代码：**

```go
// 测试服务器连通性
err := client.NewPingServer().Do(context.Background())
if err != nil {
    log.Printf("服务器连接失败: %v", err)
} else {
    log.Println("服务器连接正常")
}
```

#### 获取服务器时间

##### 查询服务器时间戳

获取 Quantum Execute 服务器的当前时间戳（毫秒）。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| 无需参数 | - | - | - |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| serverTimeMilli | int64 | 服务器时间戳（毫秒） |

**示例代码：**

```go
// 获取服务器时间戳
timestamp, err := client.NewTimestampService().Do(context.Background())
if err != nil {
    log.Fatal(err)
}

log.Printf("服务器时间戳: %d", timestamp)
log.Printf("服务器时间: %s", time.Unix(timestamp/1000, 0).Format("2006-01-02 15:04:05"))
```

#### 交易对管理

##### 查询交易对列表

获取支持的交易对信息，包括现货和合约交易对。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| page | int32 | 否 | 页码 |
| pageSize | int32 | 否 | 每页数量 |
| exchange | string | 否 | 交易所名称筛选，可选值：Binance、OKX、LTP、Deribit |
| marketType | string | 否 | 市场类型筛选，可选值：SPOT（现货）、FUTURES（合约） |
| isCoin | bool | 否 | 是否查询币本位合约可用交易对。传 `true` 时返回币本位合约可用交易对，仅 Binance 可用 |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| items | array | 交易对列表 |
| ├─ id | int | 交易对 ID |
| ├─ symbol | string | 交易对符号（如：BTCUSDT） |
| ├─ baseAsset | string | 基础币种（如：BTC） |
| ├─ quoteAsset | string | 计价币种（如：USDT） |
| ├─ exchange | string | 交易所名称 |
| ├─ marketType | string | 市场类型（SPOT/FUTURES） |
| ├─ contractType | string | 合约类型（仅合约交易对） |
| ├─ deliveryDate | string | 交割日期（仅合约交易对） |
| ├─ status | string | 交易对状态 |
| ├─ createdAt | string | 创建时间 |
| ├─ updatedAt | string | 更新时间 |
| page | int | 当前页码 |
| pageSize | int | 每页数量 |
| total | string | 总数 |

**示例代码：**

```go
// 获取所有交易对
pairs, err := client.NewTradingPairsService().Do(context.Background())
if err != nil {
    log.Fatal(err)
}

// 获取币安现货交易对
pairs, err := client.NewTradingPairsService().
    Exchange(trading_enums.ExchangeBinance).
    MarketType(trading_enums.TradingPairSpot).
    Page(1).
    PageSize(50).
    Do(context.Background())

// 获取OKX现货交易对
pairs, err := client.NewTradingPairsService().
    Exchange(trading_enums.ExchangeOKX).
    MarketType(trading_enums.TradingPairSpot).
    Page(1).
    PageSize(50).
    Do(context.Background())

// 获取LTP现货交易对
pairs, err := client.NewTradingPairsService().
    Exchange(trading_enums.ExchangeLTP).
    MarketType(trading_enums.TradingPairSpot).
    Page(1).
    PageSize(50).
    Do(context.Background())

// 获取合约交易对
pairs, err := client.NewTradingPairsService().
    MarketType(trading_enums.TradingPairFutures).
    Page(1).
    PageSize(100).
    Do(context.Background())

// 获取币本位合约可用交易对（仅 Binance）
pairs, err := client.NewTradingPairsService().
    IsCoin(true).
    Page(1).
    PageSize(100).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

// 打印交易对信息
for _, pair := range pairs.Items {
    log.Printf(`
交易对信息：
    符号: %s
    基础币种: %s
    计价币种: %s
    交易所: %s
    市场类型: %s
    状态: %s
    创建时间: %s
`,
        pair.Symbol,
        pair.BaseAsset,
        pair.QuoteAsset,
        pair.Exchange,
        pair.MarketType,
        pair.Status,
        pair.CreatedAt,
    )
    
    // 如果是合约交易对，显示额外信息
    if pair.MarketType == "FUTURES" {
        log.Printf("    合约类型: %s", pair.ContractType)
        if pair.DeliveryDate != "" {
            log.Printf("    交割日期: %s", pair.DeliveryDate)
        }
    }
}
```

### 交易所 API 管理

#### 查询交易所 API 列表

查询当前用户绑定的所有交易所 API 账户。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| page | int32 | 否 | 页码 |
| pageSize | int32 | 否 | 每页数量 |
| exchange | string | 否 | 交易所名称筛选，可选值：Binance、OKX、LTP、Deribit |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| items | array | API 列表 |
| ├─ id | string | API 记录的唯一标识 |
| ├─ createdAt | string | API 添加时间 |
| ├─ accountName | string | 账户名称（如：账户1、账户2） |
| ├─ exchange | string | 交易所名称（如：Binance、OKX、LTP、Deribit） |
| ├─ apiKey | string | 交易所 API Key（部分隐藏） |
| ├─ verificationMethod | string | API 验证方式（如：OAuth、API） |
| ├─ status | string | API 状态：正常、异常（不可用） |
| ├─ isValid | bool | API 是否有效 |
| ├─ isTradingEnabled | bool | 是否开启交易权限 |
| ├─ isDefault | bool | 是否为该交易所的默认账户 |
| ├─ isPm | bool | 是否为 Pm 账户 |
| total | int32 | API 总数 |
| page | int32 | 当前页码 |
| pageSize | int32 | 每页显示数量 |

**示例代码：**

```go
// 获取所有交易所 API 密钥
result, err := client.NewListExchangeApisService().Do(context.Background())
if err != nil {
    log.Fatal(err)
}

// 带分页和过滤 - 查询币安API
result, err := client.NewListExchangeApisService().
    Page(1).
    PageSize(10).
    Exchange(trading_enums.ExchangeBinance).
    Do(context.Background())

// 查询OKX API
result, err := client.NewListExchangeApisService().
    Page(1).
    PageSize(10).
    Exchange(trading_enums.ExchangeOKX).
    Do(context.Background())

// 查询LTP API
result, err := client.NewListExchangeApisService().
    Page(1).
    PageSize(10).
    Exchange(trading_enums.ExchangeLTP).
    Do(context.Background())

// 打印结果
for _, api := range result.Items {
    log.Printf("账户: %s, 交易所: %s, 状态: %s",
        api.AccountName,
        api.Exchange,
        api.Status,
    )
}
```

### 交易所余额/持仓/账户查询

所有接口均需 HMAC-SHA256 签名鉴权，公共必传参数为 `bindingId`（交易所 API Key 绑定 UUID，通过"查询交易所 API 列表"接口获取）。

**限频规则：** 同一 `bindingId` 10 分钟内最多 60 次请求。

---

#### 获取 Binance 现货账户余额

**GET** `/user/exchange-apis/account-balance`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| balances | array | 现货余额列表 |
| ├─ asset | string | 资产符号（如 BTC、USDT） |
| ├─ free | string | 可用余额 |
| ├─ locked | string | 冻结余额 |
| exchange | string | 交易所名称 |
| accountType | string | 账户类型（SPOT） |
| updateTime | string | 余额更新时间 |

**示例代码：**

```go
result, err := client.NewGetAccountBalanceService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, b := range result.Balances {
    log.Printf("资产: %s, 可用: %s, 冻结: %s", b.Asset, b.Free, b.Locked)
}
```

---

#### 获取 Binance 合约账户余额

**GET** `/user/exchange-apis/margin-balance`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| balances | array | 合约余额列表 |
| ├─ asset | string | 资产符号 |
| ├─ walletBalance | string | 钱包余额 |
| ├─ availableBalance | string | 可用余额 |
| ├─ crossWalletBalance | string | 全仓钱包余额 |
| ├─ crossUnPnl | string | 全仓未实现盈亏 |
| ├─ marginBalance | string | 保证金余额 |
| ├─ maxWithdrawAmount | string | 最大可转出余额 |
| exchange | string | 交易所名称 |
| accountType | string | 账户类型（FUTURES） |
| updateTime | string | 余额更新时间 |

**示例代码：**

```go
result, err := client.NewGetMarginBalanceService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, b := range result.Balances {
    log.Printf("资产: %s, 可用: %s", b.Asset, b.AvailableBalance)
}
```

---

#### 获取 Binance PAPI PV1 余额

**GET** `/user/exchange-apis/pv1-balance`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| exchange | string | 交易所名称 |
| accountType | string | 账户类型 |
| balances | array | 余额列表 |
| ├─ asset | string | 资产符号 |
| ├─ totalWalletBalance | string | 总钱包余额 |
| ├─ crossMarginBorrowed | string | 全仓杠杆借入 |
| ├─ crossMarginFree | string | 全仓杠杆可用 |
| ├─ crossMarginInterest | string | 全仓杠杆利息 |
| ├─ crossMarginLocked | string | 全仓杠杆锁定 |
| ├─ umWalletBalance | string | U本位合约钱包余额 |
| ├─ umUnrealizedPnl | string | U本位合约未实现盈亏 |
| ├─ cmWalletBalance | string | 币本位合约钱包余额 |
| ├─ cmUnrealizedPnl | string | 币本位合约未实现盈亏 |
| ├─ updateTime | int64 | 更新时间（毫秒时间戳） |
| ├─ negativeBalance | string | 负余额 |

**示例代码：**

```go
result, err := client.NewGetPv1BalanceService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, b := range result.Balances {
    log.Printf("资产: %s, 总余额: %s", b.Asset, b.TotalWalletBalance)
}
```

---

#### 获取 OKX 账户余额

**GET** `/user/exchange-apis/okx-account-balance`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| data | array | OKX 账户余额信息列表 |
| ├─ totalEq | string | 总权益 |
| ├─ availEq | string | 可用权益 |
| ├─ adjEq | string | 调整后权益 |
| ├─ imr | string | 初始保证金 |
| ├─ mmr | string | 维持保证金 |
| ├─ mgnRatio | string | 保证金率 |
| ├─ notionalUsd | string | 名义美元价值 |
| ├─ ordFroz | string | 订单冻结 |
| ├─ upl | string | 未实现盈亏 |
| ├─ uTime | string | 更新时间 |
| ├─ details | array | 各币种余额详情 |
| │ ├─ ccy | string | 币种 |
| │ ├─ eq | string | 权益 |
| │ ├─ eqUsd | string | 权益（USD） |
| │ ├─ availBal | string | 可用余额 |
| │ ├─ availEq | string | 可用权益 |
| │ ├─ cashBal | string | 现金余额 |
| │ ├─ frozenBal | string | 冻结余额 |
| │ ├─ upl | string | 未实现盈亏 |
| │ ├─ liab | string | 负债 |
| │ ├─ interest | string | 利息 |
| exchange | string | 交易所名称 |

**示例代码：**

```go
result, err := client.NewGetOkxAccountBalanceService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, d := range result.Data {
    log.Printf("总权益: %s, 可用权益: %s", d.TotalEq, d.AvailEq)
}
```

---

#### 获取 Binance FAPI 持仓方向双开状态

**GET** `/user/exchange-apis/fapi-position-side-dial`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| dualSidePosition | bool | 是否双向持仓（true=双向，false=单向） |
| exchange | string | 交易所名称 |
| accountType | string | 账户类型 |

**示例代码：**

```go
result, err := client.NewGetFapiPositionSideDialService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("双向持仓: %v", result.DualSidePosition)
```

---

#### 获取 Binance PAPI UM 持仓方向双开状态

**GET** `/user/exchange-apis/papi-um-position-side-dual`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| dualSidePosition | bool | 是否双向持仓 |
| exchange | string | 交易所名称 |
| accountType | string | 账户类型 |

**示例代码：**

```go
result, err := client.NewGetPapiUmPositionSideDualService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("双向持仓: %v", result.DualSidePosition)
```

---

#### 获取 OKX 持仓信息

**GET** `/user/exchange-apis/okx-account-positions`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| data | array | 持仓信息列表 |
| ├─ instId | string | 产品 ID（如 BTC-USDT-SWAP） |
| ├─ instType | string | 产品类型（SWAP/FUTURES/OPTION） |
| ├─ pos | string | 持仓数量 |
| ├─ posSide | string | 持仓方向（long/short/net） |
| ├─ avgPx | string | 平均开仓价 |
| ├─ markPx | string | 标记价格 |
| ├─ liqPx | string | 预估强平价 |
| ├─ upl | string | 未实现盈亏 |
| ├─ uplRatio | string | 未实现盈亏率 |
| ├─ lever | string | 杠杆倍数 |
| ├─ mgnMode | string | 保证金模式（cross/isolated） |
| ├─ imr | string | 初始保证金 |
| ├─ mmr | string | 维持保证金 |
| ├─ margin | string | 保证金 |
| ├─ notionalUsd | string | 仓位美元价值 |
| ├─ adl | string | 自动减仓指标 |
| ├─ ccy | string | 币种 |
| ├─ pnl | string | 平仓盈亏 |
| ├─ realizedPnl | string | 已实现盈亏 |
| ├─ fee | string | 手续费 |
| ├─ fundingFee | string | 资金费 |
| ├─ cTime | string | 创建时间 |
| ├─ uTime | string | 更新时间 |
| exchange | string | 交易所名称 |

**示例代码：**

```go
result, err := client.NewGetOkxAccountPositionsService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, p := range result.Data {
    log.Printf("产品: %s, 持仓量: %s, 未实现盈亏: %s", p.InstId, p.Pos, p.Upl)
}
```

---

#### 获取 OKX 最大可开仓数量

**GET** `/user/exchange-apis/okx-account-max-size`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |
| instId | string | 是 | 产品 ID（如 BTC-USDT-SWAP） |
| tdMode | string | 是 | 交易模式：cross（全仓）、isolated（逐仓）、cash（现货） |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| data | array | 可开仓数量信息列表 |
| ├─ ccy | string | 币种 |
| ├─ instId | string | 产品 ID |
| ├─ maxBuy | string | 最大可买数量 |
| ├─ maxSell | string | 最大可卖数量 |
| exchange | string | 交易所名称 |

**示例代码：**

```go
result, err := client.NewGetOkxAccountMaxSizeService().
    BindingId("your-binding-id").
    InstId("BTC-USDT-SWAP").
    TdMode("cross").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, d := range result.Data {
    log.Printf("最大买入: %s, 最大卖出: %s", d.MaxBuy, d.MaxSell)
}
```

---

#### 获取 LTP 持仓信息

**GET** `/user/exchange-apis/ltp-position`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |
| sym | string | 否 | 交易对符号，不填则查询所有持仓 |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| data | array | 持仓信息列表 |
| ├─ positionId | string | 持仓 ID |
| ├─ portfolioId | string | 投资组合 ID |
| ├─ sym | string | 交易对 |
| ├─ positionSide | string | 持仓方向 |
| ├─ positionQty | string | 持仓数量 |
| ├─ positionValue | string | 持仓价值 |
| ├─ positionMargin | string | 持仓保证金 |
| ├─ positionMm | string | 持仓维持保证金 |
| ├─ unrealizedPnl | string | 未实现盈亏 |
| ├─ unrealizedPnlRate | string | 未实现盈亏率 |
| ├─ avgPrice | string | 平均开仓价 |
| ├─ markPrice | string | 标记价格 |
| ├─ liqPrice | string | 强平价格 |
| ├─ leverage | string | 杠杆倍数 |
| ├─ maxLeverage | string | 最大杠杆 |
| ├─ riskLevel | string | 风险等级 |
| ├─ fee | string | 手续费 |
| ├─ fundingFee | string | 资金费 |
| ├─ tpslOrder | string | 止盈止损订单（JSON 字符串） |
| ├─ createAt | string | 创建时间 |
| ├─ updateAt | string | 更新时间 |
| exchange | string | 交易所名称 |

**示例代码：**

```go
// 查询所有持仓
result, err := client.NewGetLtpPositionService().
    BindingId("your-binding-id").
    Do(context.Background())

// 查询指定交易对持仓
result, err = client.NewGetLtpPositionService().
    BindingId("your-binding-id").
    Sym("BTCUSDT").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, p := range result.Data {
    log.Printf("交易对: %s, 持仓量: %s", p.Sym, p.PositionQty)
}
```

---

#### 获取 Deribit 持仓信息

**GET** `/user/exchange-apis/deribit-position`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| data | array | 持仓信息列表 |
| ├─ instrumentName | string | 合约名称 |
| ├─ direction | string | 持仓方向（buy/sell） |
| ├─ size | float64 | 持仓数量 |
| ├─ averagePrice | float64 | 平均开仓价 |
| ├─ markPrice | float64 | 标记价格 |
| ├─ indexPrice | float64 | 指数价格 |
| ├─ floatingProfitLoss | float64 | 浮动盈亏 |
| ├─ totalProfitLoss | float64 | 总盈亏 |
| ├─ initialMargin | float64 | 初始保证金 |
| ├─ maintenanceMargin | float64 | 维持保证金 |
| ├─ estimatedLiquidationPrice | float64 | 预估强平价格 |
| ├─ leverage | int32 | 杠杆倍数 |
| ├─ kind | string | 合约类型（future/option） |
| ├─ sizeCurrency | float64 | 币本位持仓数量 |
| ├─ delta | float64 | Delta |
| ├─ realizedFunding | float64 | 已实现资金费 |
| ├─ realizedProfitLoss | float64 | 已实现盈亏 |
| ├─ settlementPrice | float64 | 结算价格 |
| exchange | string | 交易所名称 |

**示例代码：**

```go
result, err := client.NewGetDeribitPositionService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, p := range result.Data {
    log.Printf("合约: %s, 方向: %s, 数量: %f", p.InstrumentName, p.Direction, p.Size)
}
```

---

#### 获取 Binance PAPI UM 账户

**GET** `/user/exchange-apis/um-account`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| tradeGroupId | int32 | 交易组 ID |
| assets | array | 资产列表 |
| ├─ asset | string | 资产符号 |
| ├─ crossWalletBalance | string | 全仓钱包余额 |
| ├─ crossUnPnl | string | 全仓未实现盈亏 |
| ├─ maintMargin | string | 维持保证金 |
| ├─ initialMargin | string | 初始保证金 |
| ├─ positionInitialMargin | string | 持仓初始保证金 |
| ├─ openOrderInitialMargin | string | 挂单初始保证金 |
| ├─ updateTime | int64 | 更新时间（毫秒时间戳） |
| positions | array | 持仓列表 |
| ├─ symbol | string | 交易对 |
| ├─ positionAmt | string | 持仓数量 |
| ├─ positionSide | string | 持仓方向 |
| ├─ entryPrice | string | 开仓价 |
| ├─ breakEvenPrice | string | 盈亏平衡价 |
| ├─ unrealizedProfit | string | 未实现盈亏 |
| ├─ leverage | string | 杠杆 |
| ├─ initialMargin | string | 初始保证金 |
| ├─ maintMargin | string | 维持保证金 |
| exchange | string | 交易所名称 |
| accountType | string | 账户类型 |
| updateTime | string | 更新时间 |

**示例代码：**

```go
result, err := client.NewGetUmAccountService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("交易组ID: %d, 持仓数: %d", result.TradeGroupId, len(result.Positions))
```

---

#### 获取 Binance PAPI CM 账户

**GET** `/user/exchange-apis/cm-account`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

与 PAPI UM 账户一致（不含 `tradeGroupId` 字段）。

**示例代码：**

```go
result, err := client.NewGetCmAccountService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("资产数: %d, 持仓数: %d", len(result.Assets), len(result.Positions))
```

---

#### 获取 Binance PAPI PV1 账户

**GET** `/user/exchange-apis/pv1-account`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| exchange | string | 交易所名称 |
| accountType | string | 账户类型 |
| uniMmr | string | 统一维持保证金率 |
| accountEquity | string | 账户权益 |
| actualEquity | string | 实际权益 |
| accountInitialMargin | string | 账户初始保证金 |
| accountMaintMargin | string | 账户维持保证金 |
| accountStatus | string | 账户状态 |
| virtualMaxWithdrawAmount | string | 虚拟最大可提取金额 |
| totalAvailableBalance | string | 总可用余额 |
| totalMarginOpenLoss | string | 总保证金未平仓亏损 |
| updateTime | string | 更新时间 |

**示例代码：**

```go
result, err := client.NewGetPv1AccountService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("账户权益: %s, 统一维持保证金率: %s", result.AccountEquity, result.UniMmr)
```

---

#### 获取 Binance DAPI 账户

**GET** `/user/exchange-apis/dapi-account`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| assets | array | 资产列表 |
| ├─ asset | string | 资产符号 |
| ├─ walletBalance | string | 钱包余额 |
| ├─ unrealizedProfit | string | 未实现盈亏 |
| ├─ marginBalance | string | 保证金余额 |
| ├─ availableBalance | string | 可用余额 |
| ├─ maxWithdrawAmount | string | 最大可提取金额 |
| ├─ crossWalletBalance | string | 全仓钱包余额 |
| ├─ crossUnPnl | string | 全仓未实现盈亏 |
| positions | array | 持仓列表 |
| ├─ symbol | string | 交易对 |
| ├─ positionAmt | string | 持仓数量 |
| ├─ positionSide | string | 持仓方向 |
| ├─ entryPrice | string | 开仓价 |
| ├─ breakEvenPrice | string | 盈亏平衡价 |
| ├─ unrealizedProfit | string | 未实现盈亏 |
| ├─ leverage | string | 杠杆 |
| ├─ isolated | bool | 是否逐仓 |
| canDeposit | bool | 是否可充值 |
| canTrade | bool | 是否可交易 |
| canWithdraw | bool | 是否可提现 |
| feeTier | int32 | 手续费等级 |
| exchange | string | 交易所名称 |
| accountType | string | 账户类型 |
| updateTime | string | 更新时间 |

**示例代码：**

```go
result, err := client.NewGetDapiAccountService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("资产数: %d, 持仓数: %d", len(result.Assets), len(result.Positions))
```

---

#### 获取 Binance FAPI 账户

**GET** `/user/exchange-apis/fapi-account`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| totalWalletBalance | string | 总钱包余额 |
| totalUnrealizedProfit | string | 总未实现盈亏 |
| totalMarginBalance | string | 总保证金余额 |
| availableBalance | string | 可用余额 |
| maxWithdrawAmount | string | 最大可提取金额 |
| totalInitialMargin | string | 总初始保证金 |
| totalMaintMargin | string | 总维持保证金 |
| totalCrossWalletBalance | string | 总全仓钱包余额 |
| totalCrossUnPnl | string | 总全仓未实现盈亏 |
| assets | array | 资产列表 |
| ├─ asset | string | 资产符号 |
| ├─ walletBalance | string | 钱包余额 |
| ├─ unrealizedProfit | string | 未实现盈亏 |
| ├─ marginBalance | string | 保证金余额 |
| ├─ availableBalance | string | 可用余额 |
| ├─ maxWithdrawAmount | string | 最大可提取金额 |
| ├─ marginAvailable | bool | 是否可作保证金 |
| positions | array | 持仓列表 |
| ├─ symbol | string | 交易对 |
| ├─ positionAmt | string | 持仓数量 |
| ├─ positionSide | string | 持仓方向 |
| ├─ entryPrice | string | 开仓价 |
| ├─ breakEvenPrice | string | 盈亏平衡价 |
| ├─ unrealizedProfit | string | 未实现盈亏 |
| ├─ leverage | string | 杠杆 |
| ├─ isolated | bool | 是否逐仓 |
| ├─ notional | string | 名义价值 |
| ├─ isolatedWallet | string | 逐仓钱包余额 |
| exchange | string | 交易所名称 |
| accountType | string | 账户类型 |

**示例代码：**

```go
result, err := client.NewGetFapiAccountService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("总钱包余额: %s, 可用余额: %s", result.TotalWalletBalance, result.AvailableBalance)
```

---

#### 获取 Binance 全仓杠杆账户详情

**GET** `/user/exchange-apis/cross-margin-account-detail`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| exchange | string | 交易所名称 |
| accountType | string | 账户类型 |
| borrowEnabled | bool | 是否允许借入 |
| tradeEnabled | bool | 是否允许交易 |
| transferEnabled | bool | 是否允许转账 |
| marginLevel | string | 保证金水平 |
| totalAssetOfBtc | string | 总资产（BTC 计） |
| totalLiabilityOfBtc | string | 总负债（BTC 计） |
| totalNetAssetOfBtc | string | 总净资产（BTC 计） |
| userAssets | array | 各币种资产详情 |
| ├─ asset | string | 资产符号 |
| ├─ free | string | 可用余额 |
| ├─ locked | string | 锁定余额 |
| ├─ borrowed | string | 借入金额 |
| ├─ interest | string | 利息 |
| ├─ netAsset | string | 净资产 |

**示例代码：**

```go
result, err := client.NewGetCrossMarginAccountDetailService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("保证金水平: %s, 总净资产(BTC): %s", result.MarginLevel, result.TotalNetAssetOfBtc)
```

---

#### 获取 LTP 账户信息

**GET** `/user/exchange-apis/ltp-account`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| data | array | 账户信息列表 |
| ├─ portfolioId | string | 投资组合 ID |
| ├─ exchangeType | string | 交易所类型 |
| ├─ equity | string | 权益 |
| ├─ maintainMargin | string | 维持保证金 |
| ├─ positionValue | string | 持仓价值 |
| ├─ uniMmr | string | 统一维持保证金率 |
| ├─ riskRatio | string | 风险率 |
| ├─ accountStatus | string | 账户状态 |
| ├─ availableMargin | string | 可用保证金 |
| ├─ validMargin | string | 有效保证金 |
| ├─ frozenMargin | string | 冻结保证金 |
| ├─ upnl | string | 未实现盈亏 |
| ├─ positionMode | string | 持仓模式 |
| exchange | string | 交易所名称 |

**示例代码：**

```go
result, err := client.NewGetLtpAccountService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, a := range result.Data {
    log.Printf("组合ID: %s, 权益: %s, 风险率: %s", a.PortfolioId, a.Equity, a.RiskRatio)
}
```

---

#### 获取 LTP 投资组合资产

**GET** `/user/exchange-apis/ltp-portfolio-asset`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| data | array | 资产信息列表 |
| ├─ portfolioId | string | 投资组合 ID |
| ├─ coin | string | 币种 |
| ├─ exchangeType | string | 交易所类型 |
| ├─ available | string | 可用余额 |
| ├─ frozen | string | 冻结余额 |
| ├─ equity | string | 权益 |
| ├─ balance | string | 余额 |
| ├─ borrow | string | 借币 |
| ├─ debt | string | 负债 |
| ├─ marginValue | string | 保证金价值 |
| ├─ indexPrice | string | 指数价格 |
| ├─ maxTransferable | string | 最大可转出金额 |
| ├─ upnl | string | 未实现盈亏 |
| ├─ equityValue | string | 权益价值 |
| ├─ createAt | string | 创建时间 |
| ├─ updateAt | string | 更新时间 |
| exchange | string | 交易所名称 |

**示例代码：**

```go
result, err := client.NewGetLtpPortfolioAssetService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, a := range result.Data {
    log.Printf("币种: %s, 可用: %s, 权益: %s", a.Coin, a.Available, a.Equity)
}
```

---

#### 获取 Deribit 账户信息

**GET** `/user/exchange-apis/deribit-account`

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| bindingId | string | 是 | 交易所 API Key 绑定 UUID |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| data | array | 账户信息列表（按币种分组） |
| ├─ currency | string | 币种（如 BTC、ETH） |
| ├─ equity | float64 | 权益 |
| ├─ balance | float64 | 余额 |
| ├─ availableFunds | float64 | 可用资金 |
| ├─ availableWithdrawalFunds | float64 | 可用提现资金 |
| ├─ marginBalance | float64 | 保证金余额 |
| ├─ initialMargin | float64 | 初始保证金 |
| ├─ maintenanceMargin | float64 | 维持保证金 |
| ├─ lockedBalance | float64 | 锁定余额 |
| ├─ totalPl | float64 | 总盈亏 |
| ├─ sessionUpl | float64 | 会话未实现盈亏 |
| ├─ sessionRpl | float64 | 会话已实现盈亏 |
| ├─ futuresPl | float64 | 期货盈亏 |
| ├─ optionsValue | float64 | 期权价值 |
| ├─ optionsDelta | float64 | 期权 Delta |
| ├─ optionsGamma | float64 | 期权 Gamma |
| ├─ optionsVega | float64 | 期权 Vega |
| ├─ optionsTheta | float64 | 期权 Theta |
| ├─ deltaTotal | float64 | 总 Delta |
| ├─ marginModel | string | 保证金模型 |
| ├─ portfolioMarginingEnabled | bool | 是否启用投资组合保证金 |
| ├─ crossCollateralEnabled | bool | 是否启用交叉抵押 |
| exchange | string | 交易所名称 |

**示例代码：**

```go
result, err := client.NewGetDeribitAccountService().
    BindingId("your-binding-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, a := range result.Data {
    log.Printf("币种: %s, 权益: %f, 可用资金: %f", a.Currency, a.Equity, a.AvailableFunds)
}
```

---

### 交易订单管理

#### 创建主订单

创建新的主订单并提交到算法侧执行。

**请求参数：**

| 参数名 | 类型      | 是否必传 | 描述                                                                                                                                                                                         |
|--------|---------|--------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **基础参数** |
| strategyType | string  | 是 | 策略类型，可选值：TWAP-1、POV                                                                                                                                                                        |
| algorithm | string  | 是 | 交易算法。strategyType=TWAP-1时，可选值：TWAP、VWAP、BoostVWAP、BoostTWAP；strategyType=POV时，可选值：POV                                                                                                      |
| exchange | string  | 是 | 交易所名称，可选值：Binance、OKX、LTP、Deribit                                                                                                                                                          |
| symbol | string  | 是 | 交易对符号（如：BTCUSDT）（可用交易对查询）                                                                                                                                                                  |
| marketType | string  | 是 | 可选值：SPOT（现货）、PERP（永续合约）                                                                                                                                                                    |
| side | string  | 是 | 1.如果isTargetPosition=False：side代表交易方向，可选值：buy（买入）、sell（卖出）；合约交易时可与reduceOnly组合，reduceOnly=True时：buy代表买入平空，sell代表卖出平多。2.如果isTargetPosition=True：side代表仓位方向，可选值：buy（多头）、sell（空头）。【仅合约交易时需传入】 |
| apiKeyId | string  | 是 | 指定使用的 API Key ID，这将决定您本次下单使用哪个交易所账户执行                                                                                                                                                      |
| **数量参数（二选一）** |
| totalQuantity | float64 | 否* | 要交易的总数量，与 orderNotional 二选一，输入范围：>0。Deribit 下单 BTCUSD/ETHUSD 时该字段单位为 USD；Binance 下单 `perp_cm` 时该字段单位为张，且必须为整数                                                                                                       |
| orderNotional | float64 | 否* | 按价值下单时的金额，以计价币种为单位（如ETHUSDT为USDT数量），与 totalQuantity 二选一，输入范围：>0。Deribit 下单 BTCUSD/ETHUSD 与 Binance 下单 `perp_cm` 时当前字段不可用                                                                                           |
| **下单模式参数** |
| isTargetPosition | bool    | 否 | 是否为目标仓位下单，默认为 false                                                                                                                                                                        |
| **时间参数** |
| startTime | string  | 否 | 交易执行的启动时间，传入格式：ISO 8601(2025-09-03T01:30:00+08:00)，若不传入，则立即执行                                                                                                                              |
| executionDuration | int32     | 否 | 订单最大执行时长，分钟，范围>=1                                                                                                                                                                          |
| executionDurationSeconds | int32     | 否 | 执行时长（秒），仅 TWAP-1 使用。当提供此字段且>0时，优先使用此字段。必须大于10秒                                                                                                                                                     |
| **TWAP/VWAP 算法参数** |
| mustComplete | bool    | 否 | 是否一定要在executionDuration之内执行完毕，选false则不会追赶进度，默认：true                                                                                                                                        |
| makerRateLimit | float64  | 否 | 要求maker占比超过该值，范围：0-1（包含0和1。输入0.1代表10%），默认：-1(算法智能计算推荐值执行)                                                                                                                                  |
| povLimit | string  | 否 | 占市场成交量比例上限，优先级低于mustComplete，范围：0-1（包含0和1。输入0.1代表10%），默认：0.8                                                                                                                               |
| limitPrice | float64       | 否 | 最高/低允许交易的价格，买入时该字段象征最高买入价，卖出时该字段象征最低卖出价，若市价超出范围则停止交易，范围：>0，默认：-1，代表无限制                                                                                                                     |
| upTolerance | string  | 否 | 允许超出目标进度的最大容忍度，范围：0-1（不包含0和1，最小输入0.0001，最大输入0.9999。输入0.1代表可以超前目标进度10%），默认：-1（即无容忍）                                                                                                         |
| lowTolerance | string  | 否 | 允许落后目标进度的最大容忍度，范围：0-1（不包含0和1，最小输入0.0001，最大输入0.9999。输入0.1代表可落后目标进度10%），默认：-1（即无容忍）                                                                                                          |
| strictUpBound | bool    | 否 | 是否严格小于uptolerance，开启后会更加严格贴近交易进度执行，同时可能会把母单拆很细；如需严格控制交易进度则建议开启，其他场景建议不开启，默认：false                                                                                                          |
| tailOrderProtection | bool    | 否 | 订单余量小于交易所最小发单量时，是否必须taker扫完，如果false，则订单余量小于交易所最小发单量时，订单结束执行；如果true，则订单余量随最近一笔下单全额执行（可能会提高Taker率），默认：true                                                                                   |
| **POV 算法参数** |
| makerRateLimit | float64  | 否 | 要求maker占比超过该值（包含0和1，输入0.1代表10%），输入范围：0-1（输入0.1代表10%），默认：-1(算法智能计算推荐值执行)                                                                                                                    |
| povLimit | string  | 否 | 占市场成交量比例上限（包含0和0.5，一般建议小于0.15），输入范围：0-0.5（povMinLimit < max(povLimit-0.01,0)），默认：0                                                                                                         |
| povMinLimit | float64  | 否 | 占市场成交量比例下限，范围：小于max(POVLimit-0.01,0)，默认：0（即无下限）                                                                                                                                            |
| limitPrice | float64       | 否 | 最高/低允许交易的价格，买入时该字段象征最高买入价，卖出时该字段象征最低卖出价，若市价超出范围则停止交易，范围：>0，默认：-1，代表无限制                                                                                                                     |
| strictUpBound | bool    | 否 | 是否追求严格小于povLimit，开启后可能会把很小的母单也拆的很细，比如50u拆成10个5u，不建议开启，算法的每个order会权衡盘口流动性，默认：false                                                                                                          |
| tailOrderProtection | bool    | 否 | 订单余量小于交易所最小发单量时，是否必须taker扫完，如果false，则订单余量小于交易所最小发单量时，订单结束执行；如果true，则订单余量随最近一笔下单全额执行（可能会提高Taker率），默认：true                                                                                   |
| **其他参数** |
| reduceOnly | bool    | 否 | 合约交易时是否仅减仓，默认值：false                                                                                                                                                                       |
| marginType | string  | 否* | **永续合约必传参数** - 合约交易保证金类型，可选值：U（U本位）、C（币本位）。当 marketType 为 PERP（永续合约）时必传；其中 `C` 对应 Binance 币本位合约                                                                                                      |
| isMargin | bool    | 否 | 是否使用现货杠杆。- 默认为false - 仅现货可使用该字段                                                                                                                                                            |
| notes | string  | 否 | 订单备注                                                                                                                                                                                       |
| enableMake | bool  | 否 | 是否允许挂单，如果关闭则全部吃单 - 默认：true                                                                                                                                                                          |
| clientOrderId | string  | 否 | 用户自定义的订单ID，用于后续通过client_order_id查询母单详情（可选）                                                                                                                                                                          |

*注：totalQuantity 和 orderNotional 必须传其中一个，但当 isTargetPosition 为 true 时，totalQuantity 必填代表目标仓位数量且 orderNotional 不可填  
*注：当使用 Deribit 账户下单 BTCUSD 或 ETHUSD 合约时，只能使用 totalQuantity 作为数量输入字段，且数量单位为 USD；orderNotional 当前不可用。  
*注：当使用 Binance 账户下单 `perp_cm`（币本位合约）时，只能使用 totalQuantity 作为数量输入字段，数量单位为张，且输入值必须为整数；orderNotional 当前不可用。  
*注：使用BoostVWAP、BoostTWAP时，代表使用高频alpha发单。仅Binance交易所可用，不适用于其他交易所。现货支持的交易对：BTCUSDT,ETHUSDT,SOLUSDT,BNBUSDT,LTCUSDT,AVAXUSDT,XLMUSDT,XRPUSDT,DOGEUSDT,CRVUSDT。永续合约支持的交易对：BTCUSDT,ETHUSDT,SOLUSDT,BNBUSDT,LTCUSDT,AVAXUSDT,XLMUSDT,XRPUSDT,DOGEUSDT,ADAUSDT,BCHUSDT,FILUSDT,1000SATSUSDT,CRVUSDT。

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| masterOrderId | string | 创建成功的主订单 ID |
| success | bool | 创建是否成功 |
| message | string | 创建结果消息 |

**示例代码：**

```go
import (
    "context"
    "log"
    qe "github.com/Quantum-Execute/qe-connector-go"
    "github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

// TWAP 订单示例 - 在 30 分钟内分批买入价值 $10,000 的 BTC
result, err := client.NewCreateMasterOrderService().
    Algorithm(trading_enums.AlgorithmTWAP).
    Exchange(trading_enums.ExchangeBinance).
    Symbol("BTCUSDT").
    MarketType(trading_enums.MarketTypeSpot).
    Side(trading_enums.OrderSideBuy).
    ApiKeyId("your-api-key-id").
    OrderNotional(10000).              // $10,000 名义价值
    StartTime("2024-01-01T10:00:00Z").
    ExecutionDuration(30).             // 30 分钟
    MustComplete(true).
    LimitPrice(60000).                 // 最高价格 $60,000
    UpTolerance("0.1").                // 允许超出 10%
    LowTolerance("0.1").               // 允许落后 10%
    StrictUpBound(false).              // 不追求严格小于uptolerance
    TailOrderProtection(true).
    StrategyType(trading_enums.StrategyTypeTWAP1).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

if result.Success {
    log.Printf("主订单创建成功，ID: %s", result.MasterOrderId)
}
```

**目标仓位下单示例：**

```go
// 目标仓位下单示例 - 买入 1.5 BTC 到目标仓位
result, err := client.NewCreateMasterOrderService().
    Algorithm(trading_enums.AlgorithmTWAP).
    Exchange(trading_enums.ExchangeBinance).
    Symbol("BTCUSDT").
    MarketType(trading_enums.MarketTypeSpot).
    Side(trading_enums.OrderSideBuy).
    ApiKeyId("your-api-key-id").
    TotalQuantity(1.5).                   // 目标数量 1.5 BTC
    IsTargetPosition(true).               // 启用目标仓位模式
    StartTime("2024-01-01T10:00:00Z").
    ExecutionDuration(60).                // 60 分钟
    MustComplete(true).
    LimitPrice(65000).                    // 最高价格 $65,000
    UpTolerance("0.1").
    LowTolerance("0.1").
    StrictUpBound(false).                 // 不追求严格小于uptolerance
    TailOrderProtection(true).
    StrategyType(trading_enums.StrategyTypeTWAP1).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

if result.Success {
    log.Printf("目标仓位订单创建成功，ID: %s", result.MasterOrderId)
}
```

**POV 算法示例：**

```go
// POV 订单示例 - 按市场成交量比例买入 BTC
result, err := client.NewCreateMasterOrderService().
    Algorithm(trading_enums.AlgorithmPOV).
    Exchange(trading_enums.ExchangeBinance).
    Symbol("BTCUSDT").
    MarketType(trading_enums.MarketTypeSpot).
    Side(trading_enums.OrderSideBuy).
    ApiKeyId("your-api-key-id").
    TotalQuantity(1.5).                // 买入 1.5 BTC
    ExecutionDuration(60).             // 60 分钟
    PovLimit(0.1).                     // 占市场成交量 10%
    PovMinLimit(0.05).                 // 最低占市场成交量 5%
    StrictUpBound(false).              // 不追求严格小于povLimit
    LimitPrice(65000).                 // 最高价格 $65,000
    TailOrderProtection(true).
    StrategyType(trading_enums.StrategyTypeTWAP1).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

if result.Success {
    log.Printf("POV 订单创建成功，ID: %s", result.MasterOrderId)
}
```

#### 查询主订单列表

获取用户的主订单列表。

**请求参数：**

| 参数名       | 类型     | 是否必传 | 描述                                  |
|-----------|--------|------|-------------------------------------|
| page      | int32  | 否    | 页码                                  |
| pageSize  | int32  | 否    | 每页数量                                |
| status    | string | 否    | 订单状态筛选，可选值：NEW（执行中）、COMPLETED（已完成）  |
| exchange  | string | 否    | 交易所名称筛选，可选值：Binance、OKX、LTP、Deribit |
| symbol    | string | 否    | 交易对筛选                               |
| startTime | string | 否    | 开始时间筛选                              |
| endTime   | string | 否    | 结束时间筛选                              |

**响应字段：**

| 字段名                    | 类型      | 描述                                                                                                                                                     |
|------------------------|---------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| items                  | array   | 主订单列表                                                                                                                                                  |
| ├─ masterOrderId       | string  | 主订单 ID                                                                                                                                                 |
| ├─ algorithm           | string  | 算法                                                                                                                                                     |
| ├─ algorithmType       | string  | 算法类型                                                                                                                                                   |
| ├─ exchange            | string  | 交易所                                                                                                                                                    |
| ├─ symbol              | string  | 交易对                                                                                                                                                    |
| ├─ marketType          | string  | 市场类型                                                                                                                                                   |
| ├─ side                | string  | 买卖方向                                                                                                                                                   |
| ├─ totalQuantity       | float64 | 按币数下单的总数量，按金额下单时，该值为0，下单数量应查看orderNotional字段                                                                                                           |         
| ├─ filledQuantity      | float64 | 1.按币数下单时，该字段代表已成交币数。2.按金额下单时，该字段值代表已成交金额                                                                                                               |      
| ├─ averagePrice        | float64 | 平均成交价                                                                                                                                                  |
| ├─ status              | string  | 状态：NEW（创建，未执行）、WAITING（等待中）、PROCESSING（执行中，且未完成）、PAUSED（已暂停）、CANCEL（取消中）、CANCELLED（已取消）、COMPLETED（已完成）、REJECTED（已拒绝）、EXPIRED（已过期）、CANCEL_REJECT（取消被拒绝） |
| ├─ executionDuration   | int32   | 执行时长（分钟）                                                                                                                                               |
| ├─ executionDurationSeconds | int32   | 执行时长（秒，仅 TWAP-1 使用；当提供且>0时优先使用；必须>10秒）                                                                                                                         |
| ├─ priceLimit          | float64 | 价格限制                                                                                                                                                   |
| ├─ startTime           | string  | 开始时间                                                                                                                                                   |
| ├─ endTime             | string  | 结束时间                                                                                                                                                   |
| ├─ createdAt           | string  | 创建时间                                                                                                                                                   |
| ├─ updatedAt           | string  | 更新时间                                                                                                                                                   |
| ├─ notes               | string  | 备注                                                                                                                                                     |
| ├─ marginType          | string  | 保证金类型（U:U本位，C:币本位）                                                                                                                                  |
| ├─ reduceOnly          | bool    | 是否仅减仓                                                                                                                                                  |
| ├─ strategyType        | string  | 策略类型                                                                                                                                                   |
| ├─ orderNotional       | string  | 订单金额（按成交额提交的下单数量）                                                                                                                                      |
| ├─ mustComplete        | bool    | 是否必须完成                                                                                                                                                 |
| ├─ makerRateLimit      | float64 | 最低 Maker 率                                                                                                                                             |
| ├─ povLimit            | float64 | 最大市场成交量占比                                                                                                                                              |
| ├─ clientId            | string  | 客户端 ID                                                                                                                                                 |
| ├─ date                | string  | 发单日期（格式：YYYYMMDD）                                                                                                                                      |
| ├─ ticktimeInt         | string  | 发单时间（格式：093000000 表示 9:30:00.000）                                                                                                                      |
| ├─ limitPriceString    | string  | 限价（字符串）                                                                                                                                                |
| ├─ upTolerance         | string  | 上容忍度                                                                                                                                                   |
| ├─ lowTolerance        | string  | 下容忍度                                                                                                                                                   |
| ├─ strictUpBound       | bool    | 严格上界                                                                                                                                                   |
| ├─ ticktimeMs          | string  | 发单时间戳（epoch 毫秒）                                                                                                                                        |   
| ├─ category            | string  | 交易品种（spot、perp 或 perp_cm，其中 perp_cm 表示 Binance 币本位合约）                                                                                     |   
| ├─ filledAmount        | float64 | 成交币数                                                                                                                                                   |
| ├─ totalValue          | float64 | 成交总值                                                                                                                                                   |
| ├─ base                | string  | 基础币种                                                                                                                                                   |
| ├─ quote               | string  | 计价币种                                                                                                                                                   |
| ├─ completionProgress  | float64 | 完成进度（0-100）返回50代表50%                                                                                                                                   |
| ├─ reason              | string  | 原因（如取消原因）                                                                                                                                              |
| ├─ tailOrderProtection | bool    | 尾单保护开关                                                                                                                                                 |
| ├─ enableMake          | bool    | 是否允许挂单                                                                                                                                                 |
| ├─ makerRate           | float64 | 被动成交率                                                                                                                                                  |
| ├─ clientOrderId       | string  | 用户自定义的订单ID                                                                                                                                              |
| total                  | int32   | 总数                                                                                                                                                     |
| page                   | int32   | 当前页码                                                                                                                                                   |
| pageSize               | int32   | 每页数量                                                                                                                                                   |

**示例代码：**

```go
// 查询所有主订单
orders, err := client.NewGetMasterOrdersService().
    Page(1).
    PageSize(20).
    Do(context.Background())

// 查询特定状态的订单
orders, err := client.NewGetMasterOrdersService().
    Page(1).
    PageSize(20).
    Status(trading_enums.MasterOrderStatusNew).  // 查询执行中的订单
    Symbol("BTCUSDT").
    StartTime("2024-01-01T00:00:00Z").
    EndTime("2024-01-31T23:59:59Z").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

// 打印订单信息
for _, order := range orders.Items {
    log.Printf(`
订单信息：
    ID: %s
    算法: %s (%s)
    交易对: %s %s
    方向: %s
    状态: %s
    完成度: %.2f%%
    平均价格: %.2f
    已成交: %.4f / %.4f
    成交金额: $%.2f
    Maker率: %.2f%%
    创建时间: %s
    发单日期: %s
    上容忍度: %s
    下容忍度: %s
`,
        order.MasterOrderId,
        order.Algorithm,
        order.StrategyType,
        order.Symbol,
        order.MarketType,
        order.Side,
        order.Status,
        order.CompletionProgress*100,
        order.AveragePrice,
        order.FilledQuantity,
        order.TotalQuantity,
        order.FilledAmount,
        order.TakerMakerRate*100,
        order.CreatedAt,
        order.Date,
        order.UpTolerance,
        order.LowTolerance,
    )
}
```

#### 获取母单详情

获取指定母单的详细信息。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| masterOrderId | string | 是 | 母单 ID |

**响应字段：**

成功时返回 `masterOrder` 字段（结构与 `MasterOrderInfo` 一致）。

**示例代码：**

```go
detail, err := client.NewGetMasterOrderDetailService().
    MasterOrderId("your-master-order-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("母单详情: %+v", detail.MasterOrder)
```

#### 通过client_order_id获取母单详情

通过用户指定的client_order_id获取母单的详细信息。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| clientOrderId | string | 是 | 用户自定义的订单ID |

**响应字段：**

成功时返回 `masterOrder` 字段（结构与 `MasterOrderInfo` 一致）。

**示例代码：**

```go
detail, err := client.NewGetMasterOrderDetailByClientOrderIdService().
    ClientOrderId("your-client-order-id").
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}
log.Printf("母单详情: %+v", detail.MasterOrder)
```

#### 查询成交记录

获取用户的成交记录。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| page | int32 | 否 | 页码 |
| pageSize | int32 | 否 | 每页数量 |
| masterOrderId | string | 否 | 主订单 ID 筛选 |
| subOrderId | string | 否 | 子订单 ID 筛选 |
| symbol | string | 否 | 交易对筛选 |
| status | string | 否 | 订单状态筛选，多个状态用逗号分隔，如：PLACED,FILLED。支持的状态：PLACED（已下单）、REJECTED（已拒单）、CANCELLED（算法已撤单）、FILLED（完全成交）、Cancelack（交易已撤单）、CANCEL_REJECTED（拒绝撤单） |
| startTime | string | 否 | 开始时间筛选 |
| endTime | string | 否 | 结束时间筛选 |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| items | array | 成交记录列表 |
| ├─ id | string | 记录 ID |
| ├─ orderCreatedTime | string | 订单创建时间 |
| ├─ masterOrderId | string | 主订单 ID |
| ├─ exchange | string | 交易所 |
| ├─ category | string | 市场类型 |
| ├─ symbol | string | 交易对 |
| ├─ side | string | 方向 |
| ├─ filledValue | float64 | 成交价值 |
| ├─ filledQuantity | string | 成交数量 |
| ├─ avgPrice | float64 | 平均价格 |
| ├─ price | float64 | 成交价格 |
| ├─ fee | float64 | 手续费 |
| ├─ tradingAccount | string | 交易账户 |
| ├─ status | string | 状态 |
| ├─ rejectReason | string | 拒绝原因 |
| ├─ base | string | 基础币种 |
| ├─ quote | string | 计价币种 |
| ├─ type | string | 订单类型 |
| total | int32 | 总数 |
| page | int32 | 当前页码 |
| pageSize | int32 | 每页数量 |

**示例代码：**

```go
// 查询特定主订单的成交明细
fills, err := client.NewGetOrderFillsService().
    MasterOrderId("your-master-order-id").
    Page(1).
    PageSize(50).
    Do(context.Background())

// 查询指定时间段的所有成交
fills, err := client.NewGetOrderFillsService().
    Symbol("BTCUSDT").
    StartTime("2024-01-01T00:00:00Z").
    EndTime("2024-01-01T23:59:59Z").
    Page(1).
    PageSize(100).
    Do(context.Background())

// 查询特定状态的成交记录
fills, err := client.NewGetOrderFillsService().
    Status("FILLED,CANCELLED").  // 查询已成交和已撤单的记录
    Symbol("BTCUSDT").
    Page(1).
    PageSize(50).
    Do(context.Background())

// 查询被拒绝的订单
fills, err := client.NewGetOrderFillsService().
    Status("REJECTED").
    StartTime("2024-01-01T00:00:00Z").
    EndTime("2024-01-31T23:59:59Z").
    Page(1).
    PageSize(100).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

// 统计成交信息
var totalValue, totalFee float64
for _, fill := range fills.Items {
    log.Printf(`
成交详情：
    时间: %s
    交易对: %s
    方向: %s
    成交价格: $%.2f
    成交数量: %.6f
    成交金额: $%.2f
    手续费: $%.4f
    账户: %s
    类型: %s
`,
        fill.OrderCreatedTime,
        fill.Symbol,
        fill.Side,
        fill.Price,
        fill.FilledQuantity,
        fill.FilledValue,
        fill.Fee,
        fill.TradingAccount,
        fill.Type,
    )
    totalValue += fill.FilledValue
    totalFee += fill.Fee
}

log.Printf("总成交额: $%.2f, 总手续费: $%.2f", totalValue, totalFee)
```

#### 查询TCA分析数据

获取TCA分析数据列表（strategy-api：APIKEY签名鉴权）。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| symbol | string | 否 | 交易对筛选 |
| category | string | 否 | 交易品种（spot、perp 或 perp_cm） |
| apikey | string | 否 | ApiKey id 列表，逗号分隔 |
| startTime | int64 | 否 | 开始时间戳（毫秒） |
| endTime | int64 | 否 | 结束时间戳（毫秒） |

**响应字段：**

成功时返回 `[]*algorithm_dto.TCAAnalysisResponse`。

| 字段名 | 类型 | 描述 |
|--------|------|------|
| MasterOrderID | string | 母单 |
| StartTime | string | 母单创建时间 |
| EndTime | string | 母单结束时间 |
| FinishedTime | string | 实际结束时间 |
| Strategy | string | 算法类型 |
| Symbol | string | 交易对 |
| Category | string | 交易类型（spot、perp 或 perp_cm） |
| Side | string | 交易方向 |
| Date | string | 母单创建日期 |
| MasterOrderQty | float64 | 母单下单币数量（如0.001 BTC） |
| MasterOrderNotional | float64 | 母单下单名义金额（如：10 USDT） |
| ArrivalPrice | float64 | 到达价格 |
| ExcutedRate | float64 | 执行率 |
| FillQty | float64 | 成交数量 |
| TakeFillNotional | float64 | Taker订单成交金额 |
| MakeFillNotional | float64 | Maker订单成交金额 |
| FillNotional | float64 | 成交金额 |
| MakerRate | float64 | 挂单率 |
| ChildOrderCnt | int | 子订单数量 |
| AverageFillPrice | float64 | 成交均价 |
| Slippage | float64 | 到达价滑点（绝对值） |
| Slippage_pct | float64 | 到达价滑点 |
| TWAP_Slippage_pct | float64 | TWAP滑点 |
| VWAP_Slippage_pct | float64 | VWAP滑点 |
| Spread | float64 | 相对买卖价差 |
| Slippage_pct_Fartouch | float64 | 到达价滑点（相比对手价） |
| TWAP_Slippage_pct_Fartouch | float64 | TWAP滑点（相比对手价） |
| VWAP_Slippage_pct_Fartouch | float64 | VWAP滑点（相比对手价） |
| IntervalReturn | float64 | 区间理论收益率 |
| ParticipationRate | float64 | 市场参与率 |

**示例代码：**

```go
import (
    "context"
    "log"
    qe "github.com/Quantum-Execute/qe-connector-go"
    "github.com/Quantum-Execute/qe-connector-go/dto/algorithm_dto"
)

// 查询TCA分析数据
items, err := client.NewGetTcaAnalysisService().
    Symbol("BTCUSDT").
    Category("spot").
    Apikey("your-apikey-id").
    StartTime(1735689600000). // 2025-01-01 00:00:00 UTC
    EndTime(1735776000000).   // 2025-01-02 00:00:00 UTC
    Do(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, it := range items {
    // it is *algorithm_dto.TCAAnalysisResponse
    log.Printf(`
TCA分析数据：
    主订单ID: %s
    交易对: %s
    方向: %s
    策略: %s
    开始时间: %s
    结束时间: %s
    完成时间: %s
    成交数量: %.4f
    平均成交价: %.2f
    Maker率: %.2f%%
    TWAP滑点: %.4f%%
    VWAP滑点: %.4f%%
    参与率: %.4f%%
`,
        it.MasterOrderID,
        it.Symbol,
        it.Side,
        it.Strategy,
        it.StartTime,
        it.EndTime,
        it.FinishedTime,
        it.FillQty,
        it.AverageFillPrice,
        it.MakerRate*100,
        it.TwapSlippagePct*100,
        it.VwapSlippagePct*100,
        it.ParticipationRate*100,
    )
}
```

#### 取消主订单

取消指定的主订单。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| masterOrderId | string | 是 | 要取消的主订单 ID |
| reason | string | 否 | 取消原因 |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| success | bool | 取消是否成功 |
| message | string | 取消结果消息 |

**示例代码：**

```go
result, err := client.NewCancelMasterOrderService().
    MasterOrderId("your-master-order-id").
    Reason("用户手动取消").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

if result.Success {
    log.Println("订单取消成功")
} else {
    log.Printf("订单取消失败：%s", result.Message)
}
```

#### 创建 ListenKey

创建一个随机的UUID作为ListenKey，绑定当前用户信息，有效期24小时。ListenKey用于WebSocket连接，可以实时接收用户相关的交易数据推送。

**请求参数：**

| 参数名 | 类型 | 是否必传 | 描述 |
|--------|------|----------|------|
| 无需参数 | - | - | - |

**响应字段：**

| 字段名 | 类型 | 描述 |
|--------|------|------|
| listenKey | string | 生成的ListenKey |
| expireAt | string | ListenKey过期时间戳（秒） |
| success | bool | 创建是否成功 |
| message | string | 创建结果消息 |

**示例代码：**

```go
// 创建 ListenKey
result, err := client.NewCreateListenKeyService().
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

if result.Success {
    log.Printf("ListenKey创建成功:")
    log.Printf("ListenKey: %s", result.ListenKey)
    log.Printf("过期时间: %s", result.ExpireAt)
    
    // 使用 ListenKey 建立 WebSocket 连接
    // wsURL := fmt.Sprintf("wss://api.quantumexecute.com/ws/%s", result.ListenKey)
} else {
    log.Printf("ListenKey创建失败：%s", result.Message)
}
```

**注意事项：**
- ListenKey 有效期为 24 小时，过期后需要重新创建
- 每个用户同时只能有一个有效的 ListenKey
- ListenKey 用于 WebSocket 连接，可以实时接收交易数据推送
- 建议在应用启动时创建 ListenKey，并在接近过期时重新创建

## 错误处理

SDK 提供了详细的错误信息，包括 API 错误和网络错误：

```go
import "github.com/Quantum-Execute/qe-connector-go/handlers"

result, err := client.NewCreateMasterOrderService().
    // ... 设置参数
    Do(context.Background())

if err != nil {
    // 检查是否为 API 错误
    if apiErr, ok := err.(*handlers.APIError); ok {
        log.Printf("API 错误 - 代码: %d, 原因: %s, 消息: %s, TraceID: %s",
            apiErr.Code,
            apiErr.Reason,
            apiErr.Message,
            apiErr.TraceId,
        )
        
        // 根据错误代码处理
        switch apiErr.Code {
        case 400:
            log.Println("请求参数错误")
        case 401:
            log.Println("认证失败")
        case 403:
            log.Println("权限不足")
        case 429:
            log.Println("请求过于频繁")
        default:
            log.Printf("其他错误: %v", apiErr)
        }
    } else {
        log.Printf("网络或其他错误: %v", err)
    }
}
```

## 高级配置

### 自定义 HTTP 客户端

```go
import (
    "net/http"
    "time"
)

// 创建自定义 HTTP 客户端
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

client := qe.NewClient("your-api-key", "your-api-secret")
client.HTTPClient = httpClient
```

### 使用代理

```go
import (
    "net/http"
    "net/url"
)

proxyURL, _ := url.Parse("http://proxy.example.com:8080")
httpClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}

client := qe.NewClient("your-api-key", "your-api-secret")
client.HTTPClient = httpClient
```

### 时间偏移调整

如果遇到时间戳错误，可以调整客户端的时间偏移：

```go
// 设置时间偏移（毫秒）
client.TimeOffset = 1000 // 客户端时间比服务器快 1 秒
```

### 请求重试

```go
// 实现简单的重试逻辑
func retryRequest(fn func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }
        
        // 检查是否应该重试
        if apiErr, ok := err.(*handlers.APIError); ok {
            // 不重试客户端错误
            if apiErr.Code >= 400 && apiErr.Code < 500 {
                return err
            }
        }
        
        // 指数退避
        time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
    }
    return err
}

// 使用重试
err := retryRequest(func() error {
    _, err := client.NewCreateMasterOrderService().
        // ... 设置参数
        Do(context.Background())
    return err
}, 3)
```

## 最佳实践

### 1. API 密钥管理

```go
// 定期检查 API 密钥状态
func checkApiKeyStatus(client *qe.Client) {
    apis, err := client.NewListExchangeApisService().Do(context.Background())
    if err != nil {
        log.Printf("获取 API 列表失败: %v", err)
        return
    }
    
    for _, api := range apis.Items {
        if !api.IsValid {
            log.Printf("警告: API %s (%s) 状态异常", api.Id, api.AccountName)
        }
    }
}
```

### 2. 订单监控

```go
// 监控订单执行状态
func monitorOrder(client *qe.Client, masterOrderId string) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            orders, err := client.NewGetMasterOrdersService().
                Page(1).
                PageSize(1).
                Do(context.Background())
            
            if err != nil {
                log.Printf("查询订单失败: %v", err)
                continue
            }
            
            if len(orders.Items) == 0 {
                log.Println("订单不存在")
                return
            }
            
            order := orders.Items[0]
            log.Printf("订单进度: %.2f%%, 状态: %s",
                order.CompletionProgress*100,
                order.Status,
            )
            
            if order.Status == "COMPLETED" {
                log.Printf("订单已结束，最终状态: %s", order.Status)
                return
            }
        }
    }
}
```

### 3. 批量处理

```go
// 批量获取所有订单
func getAllOrders(client *qe.Client) ([]qe.MasterOrderInfo, error) {
    var allOrders []qe.MasterOrderInfo
    page := int32(1)
    pageSize := int32(100)
    
    for {
        result, err := client.NewGetMasterOrdersService().
            Page(page).
            PageSize(pageSize).
            Do(context.Background())
        
        if err != nil {
            return nil, err
        }
        
        allOrders = append(allOrders, result.Items...)
        
        // 检查是否还有更多数据
        if len(result.Items) < int(pageSize) {
            break
        }
        page++
    }
    
    return allOrders, nil
}
```

### 4. ListenKey 管理

```go
package main

import (
	"context"
	"fmt"
	qe "github.com/Quantum-Execute/qe-connector-go"
	"log"
	"strconv"
	"time"
)

type ListenKeyManager struct {
	client    *qe.Client
	handlers  *qe.WebSocketEventHandlers
	listenKey string
	expireAt  int64
	wsConn    *qe.WebSocketService
}

// 创建 ListenKey 管理器
func NewListenKeyManager(client *qe.Client) *ListenKeyManager {
	handlers := &qe.WebSocketEventHandlers{
		OnConnected: func() {
			log.Printf("WebSocket connected")
		},
		OnDisconnected: func() {
			log.Printf("WebSocket disconnected")
		},
		OnError: func(err error) {
			log.Printf("WebSocket error: %v\n", err)
		},
		OnStatus: func(data string) error {
			log.Printf("Status message: %s\n", data)
			return nil
		},
		OnMasterOrder: func(msg *qe.MasterOrderMessage) error {
			log.Printf("Master Order Update:\n")
			log.Printf("  - Master Order ID: %s\n", msg.MasterOrderID)
			log.Printf("  - Symbol: %s\n", msg.Symbol)
			log.Printf("  - Side: %s\n", msg.Side)
			log.Printf("  - Quantity: %.8f\n", msg.Qty)
			log.Printf("  - Status: %s\n", msg.Status)
			log.Printf("  - Strategy: %s\n", msg.Strategy)
			if msg.Reason != "" {
				log.Printf("  - Reason: %s\n", msg.Reason)
			}
			return nil
		},
		OnOrder: func(msg *qe.OrderMessage) error {
			log.Printf("Order Update:\n")
			log.Printf("  - Order ID: %s\n", msg.OrderID)
			log.Printf("  - Master Order ID: %s\n", msg.MasterOrderID)
			log.Printf("  - Symbol: %s\n", msg.Symbol)
			log.Printf("  - Side: %s\n", msg.Side)
			log.Printf("  - Price: %.8f\n", msg.Price)
			log.Printf("  - Quantity: %.8f\n", msg.Quantity)
			log.Printf("  - Status: %s\n", msg.Status)
			log.Printf("  - Filled Qty: %.8f\n", msg.FillQty)
			log.Printf("  - Cumulative Filled: %.8f\n", msg.CumFilledQty)
			if msg.Reason != "" {
				log.Printf("  - Reason: %s\n", msg.Reason)
			}
			return nil
		},
		OnFill: func(msg *qe.FillMessage) error {
			log.Printf("Fill Update:\n")
			log.Printf("  - Order ID: %s\n", msg.OrderID)
			log.Printf("  - Master Order ID: %s\n", msg.MasterOrderID)
			log.Printf("  - Symbol: %s\n", msg.Symbol)
			log.Printf("  - Side: %s\n", msg.Side)
			log.Printf("  - Fill Price: %.8f\n", msg.FillPrice)
			log.Printf("  - Filled Qty: %.8f\n", msg.FilledQty)
			log.Printf("  - Fill Time: %s\n", time.Unix(msg.FillTime/1000, 0).Format("2006-01-02 15:04:05"))
			return nil
		},
		OnRawMessage: func(msg *qe.ClientPushMessage) error {
			// 可选：处理原始消息
			// log.Printf("Raw message - Type: %s, MessageId: %s\n", msg.Type, msg.MessageId)
			return nil
		},
	}
	m := &ListenKeyManager{
		client:   client,
		handlers: handlers,
	}
	// 启动自动刷新协程
	go m.autoRefresh()
	return m
}

// 创建或刷新 ListenKey
func (m *ListenKeyManager) createListenKey() error {
	result, err := m.client.NewCreateListenKeyService().
		Do(context.Background())

	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("创建 ListenKey 失败: %s", result.Message)
	}

	m.listenKey = result.ListenKey
	m.expireAt, _ = strconv.ParseInt(result.ExpireAt, 10, 64)

	log.Printf("ListenKey 创建成功: %s, 过期时间: %d", m.listenKey, m.expireAt)
	return nil
}

// 启动 WebSocket 连接
func (m *ListenKeyManager) StartWebSocket() error {
	if m.wsConn != nil {
		_ = m.wsConn.Close()
	}
	if err := m.createListenKey(); err != nil {
		return err
	}

	wsService := m.client.NewWebSocketService()
	wsService.SetHandlers(m.handlers)
	if err := wsService.Connect(m.listenKey); err != nil {
		log.Fatalf("Failed to connect WebSocket: %v", err)
	}
	m.wsConn = wsService

	return nil
}

// 自动刷新 ListenKey
func (m *ListenKeyManager) autoRefresh() {
	ticker := time.NewTicker(30 * time.Minute) // 每30分钟检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 检查是否接近过期（提前1小时刷新）
			if m.expireAt == 0 {
				continue
			}
			if time.Now().Unix() > m.expireAt-3600 {
				log.Println("ListenKey 即将过期，开始刷新...")
				if err := m.StartWebSocket(); err != nil {
					log.Printf("刷新 ListenKey 失败: %v", err)
				}
			}
		}
	}
}

// 使用示例
func main() {
	client := qe.NewClient("your-api-key", "your-secret-key")

	manager := NewListenKeyManager(client)

	if err := manager.StartWebSocket(); err != nil {
		log.Fatalf("启动 WebSocket 失败: %v", err)
	}

	// 保持程序运行
	select {}
}

```

### 5. WebSocket 实时数据推送

SDK 提供了完整的 WebSocket 服务，可以实时接收交易数据推送，包括主订单状态更新、子订单变化、成交明细等。

#### 创建 WebSocket 服务

```go
import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    qe "github.com/Quantum-Execute/qe-connector-go"
)

func main() {
    // 创建客户端
    client := qe.NewClient("your-api-key", "your-api-secret")
    
    // 创建 WebSocket 服务（使用默认host）
    wsService := client.NewWebSocketService()
    
    // 或者创建 WebSocket 服务并指定自定义host
    // wsService := client.NewWebSocketService("wss://custom.quantumexecute.com")
    
    // 或者使用 SetHost 方法设置自定义host
    // wsService := client.NewWebSocketService()
    // wsService.SetHost("wss://another-custom.quantumexecute.com")
    
    // 设置事件处理器
    handlers := &qe.WebSocketEventHandlers{
        OnConnected: func() {
            log.Println("WebSocket 连接成功")
        },
        OnDisconnected: func() {
            log.Println("WebSocket 连接断开")
        },
        OnError: func(err error) {
            log.Printf("WebSocket 错误: %v", err)
        },
        OnStatus: func(data string) error {
            log.Printf("状态消息: %s", data)
            return nil
        },
        OnMasterOrder: func(msg *qe.MasterOrderMessage) error {
            log.Printf("主订单更新:")
            log.Printf("  - 主订单 ID: %s", msg.MasterOrderID)
            log.Printf("  - 交易对: %s", msg.Symbol)
            log.Printf("  - 方向: %s", msg.Side)
            log.Printf("  - 数量: %.8f", msg.Qty)
            log.Printf("  - 状态: %s", msg.Status)
            log.Printf("  - 策略: %s", msg.Strategy)
            if msg.Reason != "" {
                log.Printf("  - 原因: %s", msg.Reason)
            }
            return nil
        },
        OnOrder: func(msg *qe.OrderMessage) error {
            log.Printf("子订单更新:")
            log.Printf("  - 订单 ID: %s", msg.OrderID)
            log.Printf("  - 主订单 ID: %s", msg.MasterOrderID)
            log.Printf("  - 交易对: %s", msg.Symbol)
            log.Printf("  - 方向: %s", msg.Side)
            log.Printf("  - 价格: %.8f", msg.Price)
            log.Printf("  - 数量: %.8f", msg.Quantity)
            log.Printf("  - 状态: %s", msg.Status)
            log.Printf("  - 已成交: %.8f", msg.FillQty)
            log.Printf("  - 累计成交: %.8f", msg.CumFilledQty)
            if msg.Reason != "" {
                log.Printf("  - 原因: %s", msg.Reason)
            }
            return nil
        },
        OnFill: func(msg *qe.FillMessage) error {
            log.Printf("成交更新:")
            log.Printf("  - 订单 ID: %s", msg.OrderID)
            log.Printf("  - 主订单 ID: %s", msg.MasterOrderID)
            log.Printf("  - 交易对: %s", msg.Symbol)
            log.Printf("  - 方向: %s", msg.Side)
            log.Printf("  - 成交价格: %.8f", msg.FillPrice)
            log.Printf("  - 成交数量: %.8f", msg.FilledQty)
            log.Printf("  - 成交时间: %s", time.Unix(msg.FillTime/1000, 0).Format("2006-01-02 15:04:05"))
            return nil
        },
        OnRawMessage: func(msg *qe.ClientPushMessage) error {
            // 可选：处理原始消息
            log.Printf("原始消息 - 类型: %s, 消息ID: %s", msg.Type, msg.MessageId)
            return nil
        },
    }
    
    wsService.SetHandlers(handlers)
    
    // 连接 WebSocket
    log.Println("正在连接 WebSocket...")
    
    if err := wsService.Connect(""); err != nil {
        log.Fatalf("WebSocket 连接失败: %v", err)
    }
    
    // 设置信号处理，优雅关闭
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    log.Println("WebSocket 客户端正在运行，按 Ctrl+C 退出")
    log.Println("等待订单更新...")
    
    // 等待信号
    <-sigChan
    
    log.Println("\n正在关闭...")
    
    // 关闭 WebSocket
    if err := wsService.Close(); err != nil {
        log.Printf("关闭 WebSocket 时出错: %v", err)
    }
    
    log.Println("WebSocket 客户端已停止")
}
```

#### 消息类型说明

**客户端推送消息类型：**

| 消息类型 | 描述 |
|----------|------|
| data | 数据消息 |
| status | 状态消息 |
| error | 错误消息 |
| master_data | 主订单数据 |
| order_data | 订单数据 |

**第三方消息类型：**

| 消息类型 | 描述 |
|----------|------|
| master_order | 主订单消息 |
| order | 子订单消息 |
| fill | 成交消息 |

#### 配置选项

```go
// 设置自定义host
wsService.SetHost("wss://custom.quantumexecute.com")

// 设置重连延迟
wsService.SetReconnectDelay(10 * time.Second)

// 设置心跳间隔
wsService.SetPingInterval(2 * time.Second)

// 设置 Pong 超时时间
wsService.SetPongTimeout(15 * time.Second)

// 设置日志记录器
logger := log.New(os.Stdout, "[WS] ", log.LstdFlags)
wsService.SetLogger(logger)
```

#### 自定义 WebSocket Host

SDK 支持自定义 WebSocket 连接地址，适用于以下场景：

- **测试环境**：连接到测试服务器
- **私有部署**：连接到私有部署的服务器
- **负载均衡**：连接到特定的服务器实例

**使用方式：**

```go
// 方式1：在创建时指定host
wsService := client.NewWebSocketService("wss://custom.quantumexecute.com")

// 方式2：使用 SetHost 方法
wsService := client.NewWebSocketService()
wsService.SetHost("wss://another-custom.quantumexecute.com")

// 方式3：链式调用
wsService := client.NewWebSocketService().
    SetHost("wss://custom.quantumexecute.com").
    SetReconnectDelay(5 * time.Second).
    SetPingInterval(1 * time.Second)
```

**注意事项：**
- Host 地址必须包含协议（`wss://` 或 `ws://`）
- 确保自定义host支持相同的API路径格式：`/api/ws?listen_key={listenKey}`
- 如果未设置自定义host，将使用默认地址：`wss://test.quantumexecute.com`

#### 连接状态管理

```go
// 检查连接状态
if wsService.IsConnected() {
    log.Println("WebSocket 已连接")
} else {
    log.Println("WebSocket 未连接")
}

// 手动重连
if err := wsService.Connect(listenKey); err != nil {
    log.Printf("重连失败: %v", err)
}
```

#### 错误处理

```go
handlers := &qe.WebSocketEventHandlers{
    OnError: func(err error) {
        log.Printf("WebSocket 错误: %v", err)
        
        // 根据错误类型进行处理
        if strings.Contains(err.Error(), "connection refused") {
            log.Println("连接被拒绝，可能是服务器不可用")
        } else if strings.Contains(err.Error(), "authentication failed") {
            log.Println("认证失败，请检查 ListenKey 是否有效")
        }
    },
}
```

#### 生产环境使用建议

1. **自动重连机制**：SDK 已内置自动重连功能，无需手动实现
2. **ListenKey 管理**：定期检查 ListenKey 有效性，接近过期时主动刷新
3. **错误监控**：实现完善的错误日志记录和监控
4. **负载均衡**：考虑使用多个 WebSocket 连接分散负载
5. **消息去重**：根据 `messageId` 实现消息去重处理

## 常见问题

### 1. 如何获取 API 密钥？

请登录 Quantum Execute 平台，在用户设置中创建 API 密钥。

### 2. 如何处理时间格式？

时间格式使用 ISO 8601 标准，例如：
- UTC 时间：`2024-01-01T10:00:00Z`
- 带时区：`2024-01-01T18:00:00+08:00`

### 3. 订单类型说明

- **TWAP (Time Weighted Average Price)**：时间加权平均价格算法，在指定时间段内平均分配订单
- **VWAP (Volume Weighted Average Price)**：成交量加权平均价格算法，根据市场成交量分布执行订单
- **POV (Percentage of Volume)**：成交量百分比算法，保持占市场成交量的固定比例

### 4. 枚举值说明

**算法类型 (Algorithm)：**

| 枚举常量 | 枚举值 | 描述 |
|----------|--------|------|
| `trading_enums.AlgorithmTWAP` | TWAP | TWAP算法 |
| `trading_enums.AlgorithmVWAP` | VWAP | VWAP算法 |
| `trading_enums.AlgorithmPOV` | POV | POV算法 |
| `trading_enums.AlgorithmBoostVWAP` | BoostVWAP | BoostVWAP算法（高频alpha发单） |
| `trading_enums.AlgorithmBoostTWAP` | BoostTWAP | BoostTWAP算法（高频alpha发单） |

**市场类型 (MarketType)：**

| 枚举常量 | 枚举值 | 描述 |
|----------|--------|------|
| `trading_enums.MarketTypeSpot` | SPOT | 现货市场 |
| `trading_enums.MarketTypePerp` | PERP | 合约市场 |

**订单方向 (OrderSide)：**

| 枚举常量 | 枚举值 | 描述 |
|----------|--------|------|
| `trading_enums.OrderSideBuy` | buy | 买入 |
| `trading_enums.OrderSideSell` | sell | 卖出 |

**交易所 (Exchange)：**

| 枚举常量 | 枚举值 | 描述 |
|----------|--------|------|
| `trading_enums.ExchangeBinance` | Binance | 币安 |
| `trading_enums.ExchangeOKX` | OKX | OKX |
| `trading_enums.ExchangeLTP` | LTP | LTP |
| `trading_enums.ExchangeDeribit` | Deribit | Deribit |

**保证金类型 (MarginType)：**

| 枚举常量 | 枚举值 | 描述 |
|----------|--------|------|
| `trading_enums.MarginTypeU` | U | U本位 |
| `trading_enums.MarginTypeC` | C | 币本位 |

**母单状态 (MasterOrderStatus)：**

| 枚举常量 | 枚举值 | 描述 |
|----------|--------|------|
| `trading_enums.MasterOrderStatusNew` | NEW | 执行中 |
| `trading_enums.MasterOrderStatusCompleted` | COMPLETED | 已完成 |

**策略类型 (StrategyType)：**

| 枚举常量 | 枚举值 | 描述 |
|----------|--------|------|
| `trading_enums.StrategyTypeTWAP1` | TWAP-1 | TWAP策略版本1 |
| `trading_enums.StrategyTypeTWAP2` | TWAP-2 | TWAP策略版本2 |
| `trading_enums.StrategyTypePOV` | POV | POV策略 |

**交易对市场类型 (TradingPairMarketType)：**

| 枚举常量 | 枚举值 | 描述 |
|----------|--------|------|
| `trading_enums.TradingPairSpot` | SPOT | 现货交易对 |
| `trading_enums.TradingPairFutures` | FUTURES | 合约交易对 |

### 5. 容忍度参数说明

- `upTolerance`：允许超出计划进度的容忍度，如 "0.1" 表示允许超出 10%
- `lowTolerance`：允许落后计划进度的容忍度
- `strictUpBound`：是否严格限制在 upTolerance 以内，开启后可能导致小订单被过度拆分

### 6. 新增字段说明

- `povMinLimit`：POV算法专用，占市场成交量比例下限，范围：小于max(POVLimit-0.01,0)
- `tailOrderProtection`：尾单保护，如果为true则尾单必须taker扫完，如果false则允许省一点，小于交易所最小发单量

### 7. ListenKey 相关说明

**什么是 ListenKey？**
- ListenKey 是一个用于 WebSocket 连接的身份验证令牌
- 每个 ListenKey 绑定到特定的用户账户
- 用于实时接收用户相关的交易数据推送

**ListenKey 的生命周期：**
- 有效期：24 小时
- 每个用户同时只能有一个有效的 ListenKey
- 过期后需要重新创建

**使用建议：**
- 在应用启动时创建 ListenKey
- 定期检查过期时间，提前刷新
- 实现自动重连机制
- 妥善处理 WebSocket 连接异常

### 8. WebSocket 相关说明

**WebSocket 连接地址：**
- `wss://test.quantumexecute.com/api/ws?listen_key={listenKey}`

**支持的消息类型：**
- 主订单状态更新（`master_order`）
- 子订单变化（`order`）
- 成交明细（`fill`）
- 系统状态消息（`status`）
- 错误消息（`error`）

**连接管理：**
- SDK 自动处理心跳检测和重连
- 支持自定义重连延迟、心跳间隔等参数
- 提供连接状态查询接口

**消息处理：**
- 支持结构化消息解析
- 提供原始消息访问接口
- 支持自定义错误处理逻辑

**性能优化建议：**
- 避免在消息处理器中执行耗时操作
- 使用 goroutine 处理消息，避免阻塞主连接
- 合理设置心跳参数，平衡实时性和资源消耗

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 联系我们

- 官网：[https://test.quantumexecute.com](https://test.quantumexecute.com)
- 邮箱：support@quantumexecute.com
- GitHub：[https://github.com/Quantum-Execute/qe-connector-go](https://github.com/Quantum-Execute/qe-connector-go)