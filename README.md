# Quantum Execute Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/Quantum-Execute/qe-connector-go.svg)](https://pkg.go.dev/github.com/Quantum-Execute/qe-connector-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

这是 Quantum Execute 公共 API 的官方 Go SDK，为开发者提供了一个轻量级、易于使用的接口来访问 Quantum Execute 的交易服务。

## 功能特性

- ✅ 完整的 Quantum Execute API 支持
- ✅ 交易所 API 密钥管理
- ✅ 主订单创建与管理（TWAP、VWAP 等算法）
- ✅ 订单查询和成交明细
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

### 管理交易所 API

#### 列出交易所 API 密钥

```go
// 获取所有交易所 API 密钥
result, err := client.NewListExchangeApisService().Do(context.Background())
if err != nil {
    log.Fatal(err)
}

// 带分页和过滤
result, err := client.NewListExchangeApisService().
    Page(1).
    PageSize(10).
    Exchange("binance").
    Do(context.Background())
```

#### 添加交易所 API 密钥

```go
result, err := client.NewAddExchangeApiService().
    AccountName("我的币安账户").
    Exchange("binance").
    ApiKey("your-exchange-api-key").
    ApiSecret("your-exchange-api-secret").
    EnableTrading(true).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}
log.Printf("API Key 添加成功，ID: %s", result.Id)
```

### 创建主订单

#### TWAP 订单示例

```go
// 创建一个 TWAP 订单，在 30 分钟内分批买入价值 $10,000 的 BTC
result, err := client.NewCreateMasterOrderService().
    Algorithm("TWAP").
    AlgorithmType("TIME_WEIGHTED").
    Exchange("binance").
    Symbol("BTCUSDT").
    MarketType("SPOT").
    Side("BUY").
    ApiKeyId("your-api-key-id").  // 从 ListExchangeApis 获取
    OrderNotional(10000).          // $10,000 名义价值
    StartTime("2024-01-01T10:00:00Z").
    EndTime("2024-01-01T10:30:00Z").
    ExecutionDuration("60").       // 每 60 秒执行一次
    MustComplete(true).            // 必须完成全部订单
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}
log.Printf("主订单创建成功，ID: %s", result.MasterOrderId)
```

#### VWAP 订单示例

```go
// 创建 VWAP 订单，根据市场成交量分布执行
result, err := client.NewCreateMasterOrderService().
    Algorithm("VWAP").
    AlgorithmType("VOLUME_WEIGHTED").
    Exchange("binance").
    Symbol("ETHUSDT").
    MarketType("SPOT").
    Side("SELL").
    ApiKeyId("your-api-key-id").
    TotalQuantity(5.0).           // 卖出 5 ETH
    StrategyType("AGGRESSIVE").    // 激进策略
    StartTime("2024-01-01T09:00:00Z").
    EndTime("2024-01-01T17:00:00Z").
    MakerRateLimit(0.3).          // 最多 30% Maker 订单
    PovLimit(0.1).                // 最多占市场成交量的 10%
    Do(context.Background())
```

#### 期货订单示例

```go
// 创建期货 TWAP 订单
result, err := client.NewCreateMasterOrderService().
    Algorithm("TWAP").
    AlgorithmType("TIME_WEIGHTED").
    Exchange("binance").
    Symbol("BTCUSDT").
    MarketType("FUTURES").
    Side("BUY").
    ApiKeyId("your-api-key-id").
    TotalQuantity(0.1).           // 0.1 BTC
    MarginType("CROSS").          // 全仓模式
    WorstPrice(55000).            // 最差价格限制
    UpTolerance("0.001").         // 向上容差 0.1%
    LowTolerance("0.001").        // 向下容差 0.1%
    Do(context.Background())
```

### 查询订单

#### 查询主订单列表

```go
orders, err := client.NewGetMasterOrdersService().
    Page(1).
    PageSize(20).
    Status("ACTIVE").
    Symbol("BTCUSDT").
    StartTime("2024-01-01T00:00:00Z").
    EndTime("2024-01-31T23:59:59Z").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, order := range orders.Items {
    log.Printf("订单 %s: %s %s %s, 状态: %s, 完成度: %.2f%%",
        order.MasterOrderId,
        order.Side,
        order.Symbol,
        order.Algorithm,
        order.Status,
        order.CompletionProgress*100,
    )
}
```

#### 查询成交明细

```go
fills, err := client.NewGetOrderFillsService().
    MasterOrderId("your-master-order-id").
    Page(1).
    PageSize(50).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, fill := range fills.Items {
    log.Printf("成交: %s %s @ %f, 数量: %f, 手续费: %f",
        fill.Side,
        fill.Symbol,
        fill.Price,
        fill.FilledQuantity,
        fill.Fee,
    )
}
```

### 取消订单

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
}
```

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

## API 文档

完整的 API 文档请参考 [Quantum Execute API 文档](https://docs.quantumexecute.com)

## 支持的算法类型

- **TWAP (Time Weighted Average Price)**: 时间加权平均价格算法
- **VWAP (Volume Weighted Average Price)**: 成交量加权平均价格算法
- **POV (Percentage of Volume)**: 成交量百分比算法
- **IMPLEMENTATION_SHORTFALL**: 执行缺口算法
- 更多算法类型请参考 API 文档

## 常见问题

### 1. 如何获取 API 密钥？

请登录 Quantum Execute 平台，在用户设置中创建 API 密钥。

### 2. 测试环境和生产环境的区别？

- 生产环境：`https://api.quantumexecute.com`
- 测试环境：`https://testapi.quantumexecute.com`

测试环境使用模拟数据，不会产生实际交易。

### 3. 如何处理大量数据的分页？

```go
var allOrders []qe.MasterOrderInfo
page := int32(1)

for {
    result, err := client.NewGetMasterOrdersService().
        Page(page).
        PageSize(100).
        Do(context.Background())
    
    if err != nil {
        log.Fatal(err)
    }
    
    allOrders = append(allOrders, result.Items...)
    
    if len(result.Items) < int(result.PageSize) {
        break
    }
    page++
}
```

## 示例代码

更多示例代码请参考 [examples](./examples) 目录。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 联系我们

- 官网：[https://quantumexecute.com](https://quantumexecute.com)
- 邮箱：support@quantumexecute.com
- GitHub：[https://github.com/Quantum-Execute/qe-connector-go](https://github.com/Quantum-Execute/qe-connector-go)

