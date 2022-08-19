package events

import (
	"context"
	"fmt"

	// "log"
	"time"

	"github.com/my_tdm/user-service-shu/config"
	"github.com/my_tdm/user-service-shu/pkg/logger"
	messagebroker "github.com/my_tdm/user-service-shu/pkg/messagebroker"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	kafkaWrite *kafka.Writer
	log        logger.Logger
}

func NewKafkaProducerBroker(conf config.Config, log logger.Logger, topic string) messagebroker.Publisher {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)

	return &KafkaPublisher{
		kafkaWrite: &kafka.Writer{
			Addr:                   kafka.TCP(connString),
			Topic:                  topic,
			BatchTimeout:           time.Millisecond * 10,
			AllowAutoTopicCreation: true,
		},
	}
}

//Start ...
func (p *KafkaPublisher) Start() error {
	return nil
}

// Stop ...
func (p *KafkaPublisher) Stop() error {
	err := p.kafkaWrite.Close()
	if err != nil {
		return err
	}

	return nil
}

// Publish ...
func (p *KafkaPublisher) Publish(key, body []byte, logBody string) error {
	message := kafka.Message{
		Key:   key,
		Value: body,
	}

	if err := p.kafkaWrite.WriteMessages(context.Background(), message); err != nil {
		return err
	}

	//	p.log.Info("Message published(key/body): " + string(key) + "/" + logBody)
	return nil
}
