package transporter

import (
	"context"

	"github.com/msales/transporter/kafka"
)

// Application represents the transporter application.
type Application struct {
	Consumer *kafka.Consumer
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

// IsHealthy checks the health of the Application.
func (a *Application) IsHealthy() error {
	return nil
}
