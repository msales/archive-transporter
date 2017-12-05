package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/msales/pkg/log"
	"github.com/msales/pkg/stats"
)

// ConsumerFunc represents a function that configures the Consumer.
type ConsumerFunc func(*Consumer)

// WithBrokers sets the brokers on the Consumer.
func WithBrokers(brokers []string) ConsumerFunc {
	return func(c *Consumer) {
		c.brokers = brokers
	}
}

// WithGroupID sets the group id on the Consumer.
func WithGroupID(groupID string) ConsumerFunc {
	return func(c *Consumer) {
		c.groupID = groupID
	}
}

// WithTopics sets the topics on the Consumer.
func WithTopics(topics []string) ConsumerFunc {
	return func(c *Consumer) {
		c.topics = topics
	}
}

// WithBufferSize sets the buffer size on the Consumer.
func WithBufferSize(size int) ConsumerFunc {
	return func(c *Consumer) {
		c.bufferSize = size
	}
}

// Consumer represents a buffered Kafka consumer.
type Consumer struct {
	brokers    []string
	groupID    string
	topics     []string
	bufferSize int

	kafka *cluster.Consumer

	buf       map[string]chan *sarama.ConsumerMessage
	bufTicker *time.Ticker
}

// New creates a new Consumer instance.
func New(ctx context.Context, opts ...ConsumerFunc) (*Consumer, error) {
	c := &Consumer{
		bufferSize: 1000,
	}

	for _, o := range opts {
		o(c)
	}

	if c.brokers == nil || len(c.brokers) == 0 {
		return nil, errors.New("consumer: at least one broker is required")
	}

	if len(c.groupID) == 0 {
		return nil, errors.New("consumer: a groupID is required")
	}

	if c.topics == nil || len(c.topics) == 0 {
		return nil, errors.New("consumer: at least one topic is required")
	}

	c.buf = make(map[string]chan *sarama.ConsumerMessage, len(c.topics))
	for _, topic := range c.topics {
		c.buf[topic] = make(chan *sarama.ConsumerMessage, c.bufferSize)
	}
	go c.monitorBuffers(ctx)

	config := cluster.NewConfig();
	config.Consumer.Return.Errors = true
	config.Group.Mode = cluster.ConsumerModePartitions

	consumer, err := cluster.NewConsumer(c.brokers, c.groupID, c.topics, config)
	if err != nil {
		return nil, err
	}
	c.kafka = consumer

	// Read and log errors in the Errors channel
	go c.readErrors(ctx)

	go func() {
		for {
			select {
			case pc, ok := <-c.kafka.Partitions():
				if !ok {
					return
				}

				go c.readPartition(pc)
			}
		}
	}()

	return c, nil
}

// Close closes the Consumer.
func (c *Consumer) Close() {
	if c.bufTicker != nil {
		c.bufTicker.Stop()
	}

	if c.kafka != nil {
		c.kafka.Close()
	}
}

// GetNextMessage gets the next message from the queue.
func (c *Consumer) GetNextMessage(ctx context.Context, topic string) ([]byte, error) {
	ch, ok := c.buf[topic]
	if !ok {
		return nil, nil
	}

	select {
	case msg, ok := <-ch:
		if !ok {
			return []byte{}, nil
		}

		c.kafka.MarkOffset(msg, "")

		return msg.Value, nil

	case <-ctx.Done():
		return []byte{}, nil
	}
}

func (c *Consumer) monitorBuffers(ctx context.Context) {
	c.bufTicker = time.NewTicker(1 * time.Second)
	for range c.bufTicker.C {
		for topic, ch := range c.buf {
			stats.Gauge(ctx, "buffer.used", float64(len(ch)), 1.0, map[string]string{"topic": topic})
			stats.Gauge(ctx, "buffer.used", float64(cap(ch)), 1.0, map[string]string{"topic": topic})
		}
	}
}

func (c *Consumer) readErrors(ctx context.Context) {
	for err := range c.kafka.Errors() {
		log.Error(ctx, "consumer: "+err.Error())
	}
}

func (c *Consumer) readPartition(pc cluster.PartitionConsumer) {
	for msg := range pc.Messages() {
		c.buf[msg.Topic] <- msg
	}
}
