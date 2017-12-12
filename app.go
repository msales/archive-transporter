package transporter

import (
	"context"
)

// Consumer represents a queued Kafka consumer.
type Consumer interface {
	// Close closes the Consumer.
	Close()
	// GetNextMessage gets the next message from the queue.
	GetNextMessage(ctx context.Context, topic string) ([]byte, error)
	// GetNextBatch gets the next count messages from the queue.
	GetNextBatch(ctx context.Context, topic string, count int) ([][]byte, error)
}

// Application represents the transporter application.
type Application struct {
	Consumer Consumer
}

// NewApplication creates an instance of Application.
func NewApplication() *Application {
	return &Application{}
}

// Close closes all application connections.
func (a *Application) Close() {
	a.Consumer.Close()
}

// GetNextMessage gets the next message from the queue.
func (a *Application) GetNextMessage(ctx context.Context, topic string) ([]byte, error) {
	return a.Consumer.GetNextMessage(ctx, topic)
}

// GetNextBatch gets the next count messages from the queue.
func (a *Application) GetNextBatch(ctx context.Context, topic string, count int) ([][]byte, error) {
	return a.Consumer.GetNextBatch(ctx, topic, count)
}

// IsHealthy checks the health of the Application.
func (a *Application) IsHealthy() error {
	return nil
}
