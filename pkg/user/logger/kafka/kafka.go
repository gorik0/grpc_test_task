package kafka

import (
	"context"
	"github.com/IBM/sarama"
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
	async sarama.AsyncProducer
}

func NewAsyncProducer(kafkaBrokers []string) sarama.AsyncProducer {

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewCustomPartitioner()

	producer, _ := sarama.NewAsyncProducer(kafkaBrokers, cfg)

	go func() {
		for {

			err := <-producer.Errors()
			if err != nil {
				log.Println(err)
			}
		}
	}()
	go func() {
		for {

			msg := <-producer.Successes()
			log.Println("SUCCESS ::: ", msg.Value)
		}
	}()

	return producer
}

func (k Kafka) LoNewUser(ctx context.Context) error {
	m := make(map[string]interface{})
	m["timestamp"] = time.Now().Unix()
	m["message_type"] = "INFO"
	m["message"] = "New user add"
	//bytes, err := json.Marshal(m)
	//if err != nil {
	//	return err
	//}
	println("sendinhg")
	//return k.topic[LOG].Send(ctx, &pubsub.Message{
	//	Body: []byte("s"),
	//})
	bytes := []byte("helo")
	k.async.Input() <- &sarama.ProducerMessage{
		Topic: string(LOG),
		Value: sarama.ByteEncoder(bytes),
	}
	return nil
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

	async := NewAsyncProducer(kafkaBrokers)
	kafka := &Kafka{topic: make(map[topicName]*pubsub.Topic), async: async}
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
