package events

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/my_tdm/post-service/config"
	handler "github.com/my_tdm/post-service/events/handler"
	"github.com/my_tdm/post-service/pkg/logger"
	broker "github.com/my_tdm/post-service/pkg/messagebroker"
	"github.com/my_tdm/post-service/storage"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	kafkaConsumer *kafka.Reader
	eventHandler  *handler.EventHandler
	log           logger.Logger
}

func NewKafkaConsumer(db *sqlx.DB, conf *config.Config, log logger.Logger, topic string) broker.Consumer {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)
	return &KafkaConsumer{
		kafkaConsumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{connString},
			Topic:          topic,
			MinBytes:       10e3, // 10 KB
			MaxBytes:       10e6, // 10 MB
			Partition:      0,
			CommitInterval: 0,
		}),
		eventHandler: handler.NewEventHandlerFunc(*conf, storage.NewStoragePg(db), log),
		log:          log,
	}
}

func (k KafkaConsumer) Start() {
	fmt.Println(">>> Kafka consumer started.")
	for {
		m, err := k.kafkaConsumer.ReadMessage(context.Background())
		if err != nil {
			k.log.Error("Error on consuming a message:", logger.Error(err))
			break
		}

		err = k.eventHandler.Handle(m.Value)
		if err != nil {
			k.log.Error("faild to handle consumed topic:",
				logger.String("on topic", m.Topic), logger.Error(err))
		} else {
			fmt.Println()
			k.log.Info("Successfully consumed message",
				logger.String("on topic", m.Topic),
				logger.String("message", "success"))
			fmt.Println()
		}
	}

	err := k.kafkaConsumer.Close()
	if err != nil {
		k.log.Error("Error on closing kafka reader:", logger.Error(err))
	}
}
