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

## API 参考

### 交易所 API 管理

#### 查询交易所 API 列表

查询当前用户绑定的所有交易所 API 账户。

**请求参数：**
- `page` (int32) - 页码，可选
- `pageSize` (int32) - 每页数量，可选
- `exchange` (string) - 交易所名称筛选，可选

**响应字段：**
- `items` - API 列表，包含以下字段：
  - `id` - API 记录的唯一标识
  - `createdAt` - API 添加时间
  - `accountName` - 账户名称（如：账户1、账户2）
  - `exchange` - 交易所名称（如：Binance、OKX、Bybit）
  - `apiKey` - 交易所 API Key（部分隐藏）
  - `verificationMethod` - API 验证方式（如：OAuth、API）
  - `balance` - 账户余额（美元）
  - `status` - API 状态：正常、异常（不可用）
  - `isValid` - API 是否有效
  - `isTradingEnabled` - 是否开启交易权限
  - `isDefault` - 是否为该交易所的默认账户
  - `isPm` - 是否为 Pm 账户
- `total` - API 总数
- `page` - 当前页码
- `pageSize` - 每页显示数量

**示例代码：**

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

// 打印结果
for _, api := range result.Items {
    log.Printf("账户: %s, 交易所: %s, 状态: %s, 余额: $%.2f",
        api.AccountName,
        api.Exchange,
        api.Status,
        api.Balance,
    )
}
```

#### 添加交易所 API

添加新的交易所 API 账户。

**请求参数：**
- `accountName` (string) - 账户名称，必填
- `exchange` (string) - 交易所名称（如：Binance、OKX、Bybit），必填
- `apiKey` (string) - 交易所 API Key，必填
- `apiSecret` (string) - 交易所 API Secret，必填
- `passphrase` (string) - API 密码短语（部分交易所需要），可选
- `verificationMethod` (string) - API 验证方式（如：OAuth、API），可选
- `enableTrading` (bool) - 是否开启交易权限，可选

**响应字段：**
- `id` - 新创建的 API ID
- `success` - 添加是否成功
- `message` - 操作结果消息

**示例代码：**

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

if result.Success {
    log.Printf("API Key 添加成功，ID: %s", result.Id)
} else {
    log.Printf("API Key 添加失败：%s", result.Message)
}
```

### 交易订单管理

#### 创建主订单

创建新的主订单并提交到算法侧执行。

**请求参数：**
- `algorithm` (string) - 交易算法（如：VWAP、TWAP），必填
- `algorithmType` (string) - 算法分类，暂时固定传 "TIME_WEIGHTED"，必填
- `exchange` (string) - 交易所名称（如：Binance、OKX），必填
- `symbol` (string) - 交易对符号（如：BTCUSDT），必填
- `marketType` (string) - 市场类型（SPOT:现货, FUTURES:合约），必填
- `side` (string) - 买卖方向（BUY:买入, SELL:卖出），必填
- `apiKeyId` (string) - 指定使用的 API 密钥 ID，必填
- `totalQuantity` (float64) - 要交易的总数量（按币值时使用），与 orderNotional 二选一
- `orderNotional` (float64) - 按价值下单时的金额（USDT），与 totalQuantity 二选一
- `strategyType` (string) - 策略类型（如：AGGRESSIVE、PASSIVE），可选
- `startTime` (string) - 开始执行时间（ISO 8601格式），可选
- `endTime` (string) - 结束时间（ISO 8601格式），TWAP-2 时必填
- `executionDuration` (int32) - 执行时长（秒），TWAP-1 时必填
- `worstPrice` (float64) - 最差成交价，必须完成为 true 时可选
- `limitPriceString` (string) - 限价字符串，超出范围停止交易，填 "-1" 不限制
- `upTolerance` (string) - 允许超出 schedule 的容忍度（如 "0.1" 表示 10%）
- `lowTolerance` (string) - 允许落后 schedule 的容忍度
- `strictUpBound` (bool) - 是否严格小于 upTolerance，不建议开启
- `mustComplete` (bool) - 是否必须在 duration 内执行完，默认 true
- `makerRateLimit` (float64) - 要求 maker 占比超过该值（0-1）
- `povLimit` (float64) - 占市场成交量比例限制（0-1）
- `marginType` (string) - 合约交易保证金类型（CROSS:全仓, ISOLATED:逐仓）
- `reduceOnly` (bool) - 合约交易时是否仅减仓，默认 false
- `notes` (string) - 订单备注，可选
- `clientId` (string) - 客户端唯一标识符，可选

**响应字段：**
- `masterOrderId` - 创建成功的主订单 ID
- `success` - 创建是否成功
- `message` - 创建结果消息

**示例代码：**

```go
// TWAP 订单示例 - 在 30 分钟内分批买入价值 $10,000 的 BTC
result, err := client.NewCreateMasterOrderService().
    Algorithm("TWAP").
    AlgorithmType("TIME_WEIGHTED").
    Exchange("binance").
    Symbol("BTCUSDT").
    MarketType("SPOT").
    Side("BUY").
    ApiKeyId("your-api-key-id").
    OrderNotional(10000).              // $10,000 名义价值
    StartTime("2024-01-01T10:00:00Z").
    EndTime("2024-01-01T10:30:00Z").
    ExecutionDuration("1800").         // 30 分钟 = 1800 秒
    MustComplete(true).
    WorstPrice(60000).                 // 最差价格 $60,000
    UpTolerance("0.1").                // 允许超出 10%
    LowTolerance("0.1").               // 允许落后 10%
    ClientId("my-order-001").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

if result.Success {
    log.Printf("主订单创建成功，ID: %s", result.MasterOrderId)
}

// VWAP 订单示例 - 根据市场成交量分布执行
result, err := client.NewCreateMasterOrderService().
    Algorithm("VWAP").
    AlgorithmType("VOLUME_WEIGHTED").
    Exchange("binance").
    Symbol("ETHUSDT").
    MarketType("SPOT").
    Side("SELL").
    ApiKeyId("your-api-key-id").
    TotalQuantity(5.0).               // 卖出 5 ETH
    StrategyType("AGGRESSIVE").       // 激进策略
    StartTime("2024-01-01T09:00:00Z").
    EndTime("2024-01-01T17:00:00Z").
    MakerRateLimit(0.3).              // 最少 30% Maker 订单
    PovLimit(0.1).                    // 最多占市场成交量的 10%
    LimitPriceString("3000").         // 最低卖价 $3000
    Notes("VWAP 卖出 ETH").
    Do(context.Background())

// 期货订单示例
result, err := client.NewCreateMasterOrderService().
    Algorithm("TWAP").
    AlgorithmType("TIME_WEIGHTED").
    Exchange("binance").
    Symbol("BTCUSDT").
    MarketType("FUTURES").
    Side("BUY").
    ApiKeyId("your-api-key-id").
    TotalQuantity(0.1).               // 0.1 BTC
    MarginType("CROSS").              // 全仓模式
    ReduceOnly(false).                // 非只减仓
    WorstPrice(55000).                // 最差价格限制
    UpTolerance("0.001").             // 向上容差 0.1%
    LowTolerance("0.001").            // 向下容差 0.1%
    StrictUpBound(false).             // 不严格限制上界
    ExecutionDuration("3600").        // 1 小时
    Do(context.Background())
```

#### 查询主订单列表

获取用户的主订单列表。

**请求参数：**
- `page` (int32) - 页码，可选
- `pageSize` (int32) - 每页数量，可选
- `status` (string) - 订单状态筛选，可选
- `exchange` (string) - 交易所名称筛选，可选
- `symbol` (string) - 交易对筛选，可选
- `startTime` (string) - 开始时间筛选，可选
- `endTime` (string) - 结束时间筛选，可选

**响应字段：**
- `items` - 主订单列表，每个订单包含：
  - `masterOrderId` - 主订单 ID
  - `algorithm` - 算法
  - `algorithmType` - 算法类型
  - `exchange` - 交易所
  - `symbol` - 交易对
  - `marketType` - 市场类型
  - `side` - 买卖方向
  - `totalQuantity` - 总数量
  - `filledQuantity` - 已成交数量
  - `averagePrice` - 平均成交价
  - `status` - 状态
  - `executionDuration` - 执行时长（秒）
  - `priceLimit` - 价格限制
  - `startTime` - 开始时间
  - `endTime` - 结束时间
  - `createdAt` - 创建时间
  - `updatedAt` - 更新时间
  - `notes` - 备注
  - `marginType` - 保证金类型（U:U本位, C:币本位）
  - `reduceOnly` - 是否仅减仓
  - `strategyType` - 策略类型
  - `orderNotional` - 订单金额（USDT）
  - `mustComplete` - 是否必须完成
  - `makerRateLimit` - 最低 Maker 率
  - `povLimit` - 最大市场成交量占比
  - `clientId` - 客户端 ID
  - `date` - 发单日期（格式：YYYYMMDD）
  - `ticktimeInt` - 发单时间（格式：093000000 表示 9:30:00.000）
  - `limitPriceString` - 限价（字符串）
  - `upTolerance` - 上容忍度
  - `lowTolerance` - 下容忍度
  - `strictUpBound` - 严格上界
  - `ticktimeMs` - 发单时间戳（epoch 毫秒）
  - `category` - 交易品种（spot 或 perp）
  - `filledAmount` - 成交金额
  - `totalValue` - 成交总值
  - `base` - 基础币种
  - `quote` - 计价币种
  - `completionProgress` - 完成进度（0-1）
  - `reason` - 原因（如取消原因）
- `total` - 总数
- `page` - 当前页码
- `pageSize` - 每页数量

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
    Status("ACTIVE").
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
        order.CreatedAt,
        order.Date,
        order.UpTolerance,
        order.LowTolerance,
    )
}
```

#### 查询成交记录

获取用户的成交记录。

**请求参数：**
- `page` (int32) - 页码，可选
- `pageSize` (int32) - 每页数量，可选
- `masterOrderId` (string) - 主订单 ID 筛选，可选
- `subOrderId` (string) - 子订单 ID 筛选，可选
- `symbol` (string) - 交易对筛选，可选
- `startTime` (string) - 开始时间筛选，可选
- `endTime` (string) - 结束时间筛选，可选

**响应字段：**
- `items` - 成交记录列表，每条记录包含：
  - `id` - 记录 ID
  - `orderCreatedTime` - 订单创建时间
  - `masterOrderId` - 主订单 ID
  - `exchange` - 交易所
  - `category` - 市场类型
  - `symbol` - 交易对
  - `side` - 方向
  - `filledValue` - 成交价值
  - `filledQuantity` - 成交数量
  - `avgPrice` - 平均价格
  - `price` - 成交价格
  - `fee` - 手续费
  - `tradingAccount` - 交易账户
  - `status` - 状态
  - `rejectReason` - 拒绝原因
  - `base` - 基础币种
  - `quote` - 计价币种
  - `type` - 订单类型
- `total` - 总数
- `page` - 当前页码
- `pageSize` - 每页数量

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

#### 取消主订单

取消指定的主订单。

**请求参数：**
- `masterOrderId` (string) - 要取消的主订单 ID，必填
- `reason` (string) - 取消原因，可选

**响应字段：**
- `success` - 取消是否成功
- `message` - 取消结果消息

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
        if api.Balance < 100 {
            log.Printf("警告: 账户 %s 余额不足 ($%.2f)", api.AccountName, api.Balance)
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
            
            if order.Status == "COMPLETED" || order.Status == "CANCELLED" {
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

## 常见问题

### 1. 如何获取 API 密钥？

请登录 Quantum Execute 平台，在用户设置中创建 API 密钥。

### 2. 测试环境和生产环境的区别？

- 生产环境：`https://api.quantumexecute.com`
- 测试环境：`https://testapi.quantumexecute.com`

测试环境使用模拟数据，不会产生实际交易。

### 3. 如何处理时间格式？

时间格式使用 ISO 8601 标准，例如：
- UTC 时间：`2024-01-01T10:00:00Z`
- 带时区：`2024-01-01T18:00:00+08:00`

### 4. 订单类型说明

- **TWAP (Time Weighted Average Price)**：时间加权平均价格算法，在指定时间段内平均分配订单
- **VWAP (Volume Weighted Average Price)**：成交量加权平均价格算法，根据市场成交量分布执行订单
- **POV (Percentage of Volume)**：成交量百分比算法，保持占市场成交量的固定比例
- **IMPLEMENTATION_SHORTFALL**：执行缺口算法，最小化执行成本

### 5. 容忍度参数说明

- `upTolerance`：允许超出计划进度的容忍度，如 "0.1" 表示允许超出 10%
- `lowTolerance`：允许落后计划进度的容忍度
- `strictUpBound`：是否严格限制在 upTolerance 以内，开启后可能导致小订单被过度拆分

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