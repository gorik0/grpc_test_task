package kafka

import (
	"context"
	"encoding/json"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/kafkapubsub"
	"grpc/pkg/user/logger"
	"log"
	"time"
)

type topicName string

const LOG topicName = "log"

type Kafka struct {
	topic map[topicName]*pubsub.Topic
}

func (k Kafka) LoNewUser(ctx context.Context) error {
	m := make(map[string]interface{})
	m["timestamp"] = time.Now().Unix()
	m["message_type"] = "INFO"
	m["message"] = "New user add"
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return k.topic[LOG].Send(ctx, &pubsub.Message{
		Body: bytes,
	})

}

func (k Kafka) ShutDown(ctx context.Context) error {

	for _, topic := range k.topic {
		if err := topic.Shutdown(ctx); err != nil {
			return err
		}

	}
	return nil
}

var _ logger.Logger = &Kafka{}

func NewKafka(kafkaBrokers []string) (*Kafka, error) {

	kafka := &Kafka{topic: make(map[topicName]*pubsub.Topic)}
	for _, name := range []topicName{
		LOG,
	} {

		<-time.After(time.Second * 5)
		log.Default().Printf("Create kafka %s on broker %s", name, kafkaBrokers)
		topic, err := kafkapubsub.OpenTopic(kafkaBrokers, kafkapubsub.MinimalConfig(), string(name), nil)
		if err != nil {
			return nil, err
		}
		kafka.topic[name] = topic
	}
	return kafka, nil

}
