package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

// Consumer is a consumer group member.
// Partitions are shared across every process running with the same group id.
type Consumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler func(*sarama.ConsumerMessage) error
}

func NewConsumer(brokers []string, group string, topics []string) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Return.Errors = true

	g, err := sarama.NewConsumerGroup(brokers, group, cfg)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		group:  g,
		topics: topics,
	}, nil
}

// Consume blocks until ctx is cancelled, calling handler for every message.
// A handler error stops the whole consumer, so return nil for anything the
// caller can skip. Offsets are committed automatically after each message.
func (c *Consumer) Consume(ctx context.Context, handler func(*sarama.ConsumerMessage) error) error {
	c.handler = handler

	for {
		// Consume returns on every rebalance; loop until ctx dies.
		if err := c.group.Consume(ctx, c.topics, c); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}

func (c *Consumer) Close() error {
	return c.group.Close()
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := c.handler(msg); err != nil {
			return err
		}
		sess.MarkMessage(msg, "")
	}

	return nil
}
