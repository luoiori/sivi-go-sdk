package sivi

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Client struct {
	meter    metric.Meter
	provider *sdkmetric.MeterProvider
	config   *Config
}

func NewClient(config *Config) (*Client, error) {
	ctx := context.Background()
	url := config.Sivi.SDK.MetricURL

	insecure := false
	if strings.HasPrefix(url, "http://") {
		insecure = true
		url = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	}
	var endpoint, path string
	if idx := strings.Index(url, "/"); idx != -1 {
		endpoint = url[:idx]
		path = url[idx:]
	} else {
		endpoint = url
		path = "/v1/metrics"
	}

	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithURLPath(path),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
		otlpmetrichttp.WithTemporalitySelector(func(sdkmetric.InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}),
		otlpmetrichttp.WithAggregationSelector(func(ik sdkmetric.InstrumentKind) sdkmetric.Aggregation {
			if ik == sdkmetric.InstrumentKindHistogram {
				return sdkmetric.AggregationExplicitBucketHistogram{
					Boundaries: []float64{200, 500, 1000, 3000},
				}
			}
			return sdkmetric.DefaultAggregationSelector(ik)
		}),
	}
	if insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	exporter, err := otlpmetrichttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	fmt.Printf("Exporter created with endpoint: %s%s\n", endpoint, path)

	// Create resource with global attributes
	res := resource.NewWithAttributes(
		"metrics",
		attribute.String("app", config.Sivi.SDK.App),
		attribute.String("app_id", strconv.Itoa(config.Sivi.SDK.AppID)),
		attribute.String("server", config.Sivi.SDK.Server),
		attribute.String("profile", config.Sivi.SDK.Profile),
	)

	// Create periodic reader with configurable interval
	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Duration(config.Sivi.SDK.Period)*time.Second))

	// Create meter provider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	otel.SetMeterProvider(provider)

	meter := provider.Meter("sivi-go-sdk")

	return &Client{
		meter:    meter,
		provider: provider,
		config:   config,
	}, nil
}

func (c *Client) CounterBuilder(name string) CounterBuilder {
	return CounterBuilder{
		meter: c.meter,
		name:  name,
	}
}

func (c *Client) HistogramBuilder(name string) HistogramBuilder {
	return HistogramBuilder{
		meter: c.meter,
		name:  name,
	}
}

func (c *Client) Shutdown(ctx context.Context) error {
	return c.provider.Shutdown(ctx)
}

type GaugeBuilder struct {
	meter metric.Meter
	name  string
}

func (gb GaugeBuilder) Build() Gauge {
	gauge, _ := gb.meter.Float64ObservableGauge(gb.name)
	return Gauge{gauge: gauge}
}

type Gauge struct {
	gauge metric.Float64ObservableGauge
}

func (c *Client) ForceFlush(ctx context.Context) error {
	return c.provider.ForceFlush(ctx)
}

type CounterBuilder struct {
	meter metric.Meter
	name  string
}

func (cb CounterBuilder) Build() Counter {
	counter, _ := cb.meter.Int64Counter(cb.name)
	return Counter{counter: counter}
}

type Counter struct {
	counter metric.Int64Counter
}

func (c Counter) Add(value int64, attrs attribute.Set) {
	c.counter.Add(context.Background(), value, metric.WithAttributeSet(attrs))
}

type HistogramBuilder struct {
	meter metric.Meter
	name  string
}

func (hb HistogramBuilder) Build() Histogram {
	histogram, _ := hb.meter.Float64Histogram(hb.name)
	return Histogram{histogram: histogram}
}

type Histogram struct {
	histogram metric.Float64Histogram
}

func (h Histogram) Record(value float64, attrs attribute.Set) {
	h.histogram.Record(context.Background(), value, metric.WithAttributeSet(attrs))
}
