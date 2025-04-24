package publishers

import (
	"github.com/IBM/sarama"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/adapters/config"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/core/ports"
)

type KafkaPublisher struct {
	producer sarama.AsyncProducer
}

func NewKafkaPublisher(config config.Config) (ports.IPublisher, error) {
	sconfig := sarama.NewConfig()
	sconfig.Producer.RequiredAcks = sarama.NoResponse
	sconfig.Producer.Return.Successes = false
	sconfig.Producer.Return.Errors = false

	brokers := []string{config.KafkaBroker}
	producer, err := sarama.NewAsyncProducer(brokers, sconfig)
	if err != nil {
		return nil, err
	}

	return &KafkaPublisher{
		producer: producer,
	}, nil
}

func (p *KafkaPublisher) Shutdown() {
	p.producer.Close()
}

func (p *KafkaPublisher) Publish(key string, data []byte, topic string) {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
		Key:   sarama.ByteEncoder(key),
	}

	p.producer.Input() <- message
}
