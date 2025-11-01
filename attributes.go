package sivi

import "go.opentelemetry.io/otel/attribute"

type AttributesBuilder struct {
	attrs []attribute.KeyValue
}

func NewAttributesBuilder() *AttributesBuilder {
	return &AttributesBuilder{
		attrs: make([]attribute.KeyValue, 0),
	}
}

func (b *AttributesBuilder) Put(key, value string) *AttributesBuilder {
	b.attrs = append(b.attrs, attribute.String(key, value))
	return b
}

func (b *AttributesBuilder) Build() attribute.Set {
	return attribute.NewSet(b.attrs...)
}
