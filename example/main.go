package main

import (
	"context"
	"fmt"
	"log"
	"time"

	sivi "github.com/luoiori/sivi-go-sdk"
)

func main() {
	config, err := sivi.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Create client
	client, err := sivi.NewClient(config)
	print(config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown(context.Background())

	// Example 1: HTTP server metrics - 修改值以便识别
	httpServerAttributes := sivi.NewAttributesBuilder().
		Put("callee_server", "im-server").
		Put("callee_method", "/send/message").
		Put("callee_ip", "192.168.1.1").
		Put("callee_port", "8080").
		Put("biz_code", "10000").
		Put("code", "200").
		Put("code_type", "success").
		Put("time", time.Now().Format(time.DateTime)).
		Put("timestamp", fmt.Sprintf("%d", time.Now().Unix())). // 添加时间戳
		Build()

	counter := client.CounterBuilder("rpc_server_handled_total").Build()
	counter.Add(5, httpServerAttributes) // 改为5

	latency := client.HistogramBuilder("rpc_server_handled_latency").Build()
	latency.Record(25000, httpServerAttributes) // 改为25000

	// Example 2: Redis client metrics
	redisClientAttributes := sivi.NewAttributesBuilder().
		Put("caller_method", "/send/message").
		Put("caller_server", "im-server").
		Put("callee_method", "pub").
		Put("callee_server", "Redis").
		Put("callee_ip", "192.168.1.2").
		Put("callee_port", "6379").
		Put("biz_code", "0").
		Put("code", "0").
		Put("code_type", "success").
		Build()

	clientCounter := client.CounterBuilder("rpc_client_handled_total").Build()
	clientCounter.Add(1, redisClientAttributes)

	clientLatency := client.HistogramBuilder("rpc_client_handled_latency").Build()
	clientLatency.Record(5000, redisClientAttributes)

	// Example 3: WebSocket metrics
	wsAttributes := sivi.NewAttributesBuilder().
		Put("caller_method", "sub").
		Put("caller_server", "im-server").
		Put("callee_method", "$route").
		Put("callee_server", "hawa_frontend").
		Put("data_type", "request").
		Put("biz_code", "0").
		Put("code", "0").
		Put("code_type", "success").
		Build()

	wsCounter := client.CounterBuilder("rpc_client_handled_total").Build()
	wsCounter.Add(1, wsAttributes)

	datasize := client.HistogramBuilder("rpc_client_handled_datasize").Build()
	datasize.Record(20000, wsAttributes)

	log.Println("Metrics recorded, manually triggering export...")

	// 手动触发发送
	ctx := context.Background()
	if err := client.ForceFlush(ctx); err != nil {
		log.Printf("Manual export failed: %v", err)
	} else {
		log.Println("Manual export completed successfully")
	}

	log.Println("Waiting 2 seconds for background exports...")
	time.Sleep(10 * time.Second)
}
