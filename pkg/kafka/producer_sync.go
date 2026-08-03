package kafka

import "github.com/IBM/sarama"

type SyncProducer struct {
	sarama.SyncProducer
}

func NewSyncProducer(brokers []string) (*SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	cfg.Producer.Retry.Max = 5
	cfg.Net.MaxOpenRequests = 1    // one request at a time
	cfg.Producer.Idempotent = true // no duplicates on retry

	// enable this when using the async client
	// false by default, just added to make it explicit
	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = false

	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{
		SyncProducer: syncProducer,
	}, nil
}

// Publish sends value to a topic.
// Messages with the same key land on the same partition,
// pass one when ordering matters (e.g., the stock symbol).
func (p *SyncProducer) Publish(topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	if key != nil {
		msg.Key = sarama.ByteEncoder(key)
	}

	_, _, err := p.SendMessage(msg)
	return err
}
