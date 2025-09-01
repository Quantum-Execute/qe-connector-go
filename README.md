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

### 交易订单管理

#### 创建主订单

创建新的主订单并提交到算法侧执行。

**请求参数：**

**基础参数（必填）：**
- `algorithm` (string) - 交易算法，可选值：`TWAP`、`VWAP`、`POV`，必填
- `exchange` (string) - 交易所名称，可选值：`Binance`，必填
- `symbol` (string) - 交易对符号（如：BTCUSDT），必填
- `marketType` (string) - 市场类型，可选值：`SPOT`（现货）、`PERP`（合约），必填
- `side` (string) - 买卖方向，可选值：`buy`（买入）、`sell`（卖出），必填
- `apiKeyId` (string) - 指定使用的 API 密钥 ID，必填

**数量参数（二选一）：**
- `totalQuantity` (string) - 要交易的总数量，支持字符串表示以避免精度问题，与 orderNotional 二选一，范围：>0
- `orderNotional` (string) - 按价值下单时的金额，以计价币种为单位（如ETHUSDT为USDT数量），与 totalQuantity 二选一，范围：>0

**时间参数：**
- `startTime` (string) - 开始执行时间（ISO 8601格式），可选
- `endTime` (string) - 结束时间（ISO 8601格式），TWAP-2 时必填
- `executionDuration` (int) - 订单的有效时间（分钟），TWAP-1 时必填，范围：>10

**算法特定参数：**

**TWAP/VWAP 算法参数：**
- `mustComplete` (bool) - 是否一定要在duration之内执行完，选false则不会追进度，默认：true
- `makerRateLimit` (string) - 要求maker占比超过该值（优先级低于mustcomplete），范围：0-1，默认："0"
- `povLimit` (string) - 占市场成交量比例限制，优先级低于mustcomplete，范围：0-1，默认："0.8"
- `limitPrice` (string) - 最高/低允许交易的价格，买的话就是最高价，卖就是最低价，超出范围停止交易，填"-1"不限制，范围：>0，默认："-1"
- `upTolerance` (string) - 允许超出schedule的容忍度，比如0.1就是执行过程中允许比目标进度超出母单数量的10%，范围：>0且<1，默认：-1
- `lowTolerance` (string) - 允许落后schedule的容忍度，范围：>0且<1，默认：-1
- `strictUpBound` (bool) - 是否追求严格小于uptolerance，开启后可能会把很小的母单也拆的很细，不建议开启，默认：false
- `takeMakeFeeDiff` (string) - 客户账户的taker maker commission fee的差，范围：-0.1~0.1，默认："0.000175"
- `tailOrderProtection` (bool) - 尾单必须taker扫完，如果false则允许省一点，小于交易所最小发单量，默认：true

**POV 算法参数：**
- `makerRateLimit` (string) - 要求maker占比超过该值，范围：0-1，默认："0"
- `povLimit` (string) - 占市场成交量比例限制，范围：0-0.5，默认："0.05"
- `povMinLimit` (string) - 占市场成交量比例下限，范围：小于max(POVLimit-0.01,0)，默认："0"
- `limitPrice` (string) - 最高/低允许交易的价格，买的话就是最高价，卖就是最低价，超出范围暂停交易，填"-1"不限制，范围：>0，默认："-1"
- `strictUpBound` (bool) - 是否追求严格小于povlimit，开启后可能会把很小的母单也拆的很细，不建议开启，默认：false
- `takeMakeFeeDiff` (string) - 客户账户的taker maker commission fee的差，范围：-0.1~0.1，默认："0.000175"
- `tailOrderProtection` (bool) - 尾单必须taker扫完，如果false则允许省一点，小于交易所最小发单量，默认：true

**其他参数：**

- `reduceOnly` (bool) - 合约交易时是否仅减仓，默认：false
- `marginType` (string) - 合约交易保证金类型，可选值：`U`（U本位）、`C`（币本位）
- `notes` (string) - 订单备注，可选

**响应字段：**
- `masterOrderId` - 创建成功的主订单 ID
- `success` - 创建是否成功
- `message` - 创建结果消息

**示例代码：**

```go
// TWAP 订单示例 - 在 30 分钟内分批买入价值 $10,000 的 BTC
result, err := client.NewCreateMasterOrderService().
    Algorithm("TWAP").
    Exchange("Binance").
    Symbol("BTCUSDT").
    MarketType("SPOT").
    Side("buy").
    ApiKeyId("your-api-key-id").
    OrderNotional("10000").            // $10,000 名义价值
    StartTime("2024-01-01T10:00:00Z").
    EndTime("2024-01-01T10:30:00Z").
    ExecutionDuration(30).             // 30 分钟
    MustComplete(true).
    LimitPrice("60000").               // 最高价格 $60,000
    UpTolerance("0.1").                // 允许超出 10%
    LowTolerance("0.1").               // 允许落后 10%
    TailOrderProtection(true).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

if result.Success {
    log.Printf("主订单创建成功，ID: %s", result.MasterOrderId)
}
```

**POV 算法示例：**

```go
// POV 订单示例 - 按市场成交量比例买入 BTC
result, err := client.NewCreateMasterOrderService().
    Algorithm("POV").
    Exchange("Binance").
    Symbol("BTCUSDT").
    MarketType("SPOT").
    Side("buy").
    ApiKeyId("your-api-key-id").
    TotalQuantity("1.5").              // 买入 1.5 BTC
    ExecutionDuration(60).             // 60 分钟
    PovLimit("0.1").                   // 占市场成交量 10%
    PovMinLimit("0.05").               // 最低占市场成交量 5%
    LimitPrice("65000").               // 最高价格 $65,000
    TailOrderProtection(true).
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
- `page` (int32) - 页码，可选
- `pageSize` (int32) - 每页数量，可选
- `status` (string) - 订单状态筛选，可选值：`NEW`（执行中）、`COMPLETED`（已完成），可选
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
  - `status` - 状态，可选值：`NEW`（执行中）、`COMPLETED`（已完成）
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
    Status("NEW").                    // 查询执行中的订单
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

#### 创建 ListenKey

创建一个随机的UUID作为ListenKey，绑定当前用户信息，有效期24小时。ListenKey用于WebSocket连接，可以实时接收用户相关的交易数据推送。

**请求参数：**
- 无需参数

**响应字段：**
- `listenKey` (string) - 生成的ListenKey
- `expireAt` (string) - ListenKey过期时间戳（秒）
- `success` (bool) - 创建是否成功
- `message` (string) - 创建结果消息

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
// ListenKey 管理器
type ListenKeyManager struct {
    client     *qe.Client
    listenKey  string
    expireAt   int64
    wsConn     *websocket.Conn
    reconnect  chan bool
}

// 创建 ListenKey 管理器
func NewListenKeyManager(client *qe.Client) *ListenKeyManager {
    return &ListenKeyManager{
        client:    client,
        reconnect: make(chan bool, 1),
    }
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
    if err := m.createListenKey(); err != nil {
        return err
    }
    
    // 建立 WebSocket 连接
    wsURL := fmt.Sprintf("wss://api.quantumexecute.com/ws/%s", m.listenKey)
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        return err
    }
    
    m.wsConn = conn
    
    // 启动消息处理协程
    go m.handleMessages()
    
    // 启动自动刷新协程
    go m.autoRefresh()
    
    return nil
}

// 处理 WebSocket 消息
func (m *ListenKeyManager) handleMessages() {
    defer m.wsConn.Close()
    
    for {
        _, message, err := m.wsConn.ReadMessage()
        if err != nil {
            log.Printf("WebSocket 读取错误: %v", err)
            m.reconnect <- true
            return
        }
        
        // 处理接收到的消息
        log.Printf("收到消息: %s", string(message))
        
        // 这里可以解析消息并处理业务逻辑
        // 例如：订单状态更新、成交通知等
    }
}

// 自动刷新 ListenKey
func (m *ListenKeyManager) autoRefresh() {
    ticker := time.NewTicker(30 * time.Minute) // 每30分钟检查一次
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // 检查是否接近过期（提前1小时刷新）
            if time.Now().Unix() > m.expireAt-3600 {
                log.Println("ListenKey 即将过期，开始刷新...")
                if err := m.createListenKey(); err != nil {
                    log.Printf("刷新 ListenKey 失败: %v", err)
                }
            }
        case <-m.reconnect:
            log.Println("开始重连 WebSocket...")
            m.wsConn.Close()
            time.Sleep(5 * time.Second)
            if err := m.StartWebSocket(); err != nil {
                log.Printf("重连失败: %v", err)
            }
            return
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
    
    // 创建 WebSocket 服务
    wsService := client.NewWebSocketService()
    
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
- `data` - 数据消息
- `status` - 状态消息
- `error` - 错误消息
- `master_data` - 主订单数据
- `order_data` - 订单数据

**第三方消息类型：**
- `master_order` - 主订单消息
- `order` - 子订单消息
- `fill` - 成交消息

#### 配置选项

```go
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
- `TWAP` - TWAP算法
- `VWAP` - VWAP算法  
- `POV` - POV算法

**市场类型 (MarketType)：**
- `SPOT` - 现货市场
- `PERP` - 合约市场

**订单方向 (OrderSide)：**
- `buy` - 买入
- `sell` - 卖出

**交易所 (Exchange)：**
- `Binance` - 币安

**保证金类型 (MarginType)：**
- `U` - U本位
- `C` - 币本位

**母单状态 (MasterOrderStatus)：**
- `NEW` - 执行中
- `COMPLETED` - 已完成

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

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 联系我们

- 官网：[https://test.quantumexecute.com](https://test.quantumexecute.com)
- 邮箱：support@quantumexecute.com
- GitHub：[https://github.com/Quantum-Execute/qe-connector-go](https://github.com/Quantum-Execute/qe-connector-go)