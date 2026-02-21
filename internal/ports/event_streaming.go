package ports

import "context"

// EventProducer defines the interface for producing events to a streaming platform (e.g. Redpanda/Kafka)
type EventProducer interface {
	Produce(ctx context.Context, key, value []byte) error
}

// EventConsumer defines the interface for consuming events
type EventConsumer interface {
	ConsumeLoop(ctx context.Context, handler func(ctx context.Context, key, value []byte))
}
