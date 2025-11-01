# Sivi Go SDK

使用OpenTelemetry标准收集和上报指标的Go SDK。

## 特性

- 使用OpenTelemetry标准进行指标收集
- 基于YAML的配置
- 用于创建指标和属性的构建器模式
- 支持计数器和直方图
- 可配置的指标导出间隔(5秒)
- 自定义直方图桶边界(200, 500, 1000, 3000)
- 支持HTTP和HTTPS上报地址

## 安装

```bash
go get github.com/sivi/go-sdk
```

## 快速开始

### 1. 配置

复制配置模板并修改：

```bash
cp config.yaml.example config.yaml
```

编辑 `config.yaml` 文件：

```yaml
sivi:
  sdk:
    app: hawa-backend
    app-id: 100
    server: hex-hawa-im
    profile: test
    metric-url: https://hawatalk.com/v1/metrics
```

### 2. 使用示例

```go
package main

import (
    "context"
    "log"
    
    sivi "github.com/sivi/go-sdk"
)

func main() {
    // 加载配置
    config, err := sivi.LoadConfig("config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建客户端
    client, err := sivi.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Shutdown(context.Background())
    
    // 创建属性
    attributes := sivi.NewAttributesBuilder().
        Put("callee_server", "im-server").
        Put("callee_method", "/send/message").
        Put("code", "200").
        Put("code_type", "success").
        Build()
    
    // 计数器
    counter := client.CounterBuilder("rpc_server_handled_total").Build()
    counter.Add(1, attributes)
    
    // 直方图
    latency := client.HistogramBuilder("rpc_server_handled_latency").Build()
    latency.Record(20000, attributes)
}
```

## API 参考

### 客户端

- `NewClient(config *Config)` - 创建新的SDK客户端
- `CounterBuilder(name string)` - 创建计数器构建器
- `HistogramBuilder(name string)` - 创建直方图构建器
- `Shutdown(ctx context.Context)` - 关闭客户端

### 属性构建器

- `NewAttributesBuilder()` - 创建新的属性构建器
- `Put(key, value string)` - 添加属性
- `Build()` - 构建属性集

### 指标类型

- `Counter.Add(value int64, attrs attribute.Set)` - 增加计数
- `Histogram.Record(value float64, attrs attribute.Set)` - 记录直方图值
