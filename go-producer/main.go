package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mcolomerc/kafkasr/producer/data"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
)

type Config struct {
	BootstrapServers string
	Topic            string
	SchemaRegistry   string
	NumMessages      int
}

//create kafka producer with default config
func main() {
	defer TimeTrack(time.Now(), "producer")
	config := readEnvironment()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.BootstrapServers})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	send := make(chan int, config.NumMessages)

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					/* fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset) */
					send <- 1
				}
			case kafka.Error:
				log.Printf("Error: %v\n", ev)
			default:
				log.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	srClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(config.SchemaRegistry))
	if err != nil {
		panic(err)
	}

	ser, err := protobuf.NewSerializer(srClient, serde.ValueSerde, protobuf.NewSerializerConfig())

	if err != nil {
		fmt.Printf("Failed to create serializer: %s\n", err)
		os.Exit(1)
	}

	for i := 0; i < config.NumMessages; i++ {
		person := data.GetPerson()

		payload, err := ser.Serialize(config.Topic, &person)
		if err != nil {
			fmt.Printf("Failed to serialize payload: %s\n", err)
			os.Exit(1)
		}

		//Produce Message
		errD := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &config.Topic, Partition: kafka.PartitionAny},
			Key:            []byte(data.GetPersonKey()),
			Value:          payload,
			Headers:        nil,
		}, nil)
		if errD != nil {
			log.Printf("Delivery failed: %v\n", err)
		}
	}
	for i := 0; i < config.NumMessages; i++ {
		<-send
	}

}

func readEnvironment() Config {
	config := Config{}
	//Checking that an environment variable is present or not.
	bootstrapServers, ok := os.LookupEnv("BOOTSTRAP_SERVERS")
	if !ok {
		log.Println("BOOTSTRAP_SERVERS is not present")
		config.BootstrapServers = "localhost:9092"
	} else {
		config.BootstrapServers = bootstrapServers
	}

	topic, ok := os.LookupEnv("TOPIC")
	if !ok {
		fmt.Println("TOPIC is not present")
		config.Topic = "sr-topic-pb-test"
	} else {
		config.Topic = topic
	}

	numMessages, ok := os.LookupEnv("NUM_MESSAGES")
	if !ok {
		log.Println("NUM_MESSAGES is not present")
		config.NumMessages = 1000
	} else {
		var err error
		config.NumMessages, err = strconv.Atoi(numMessages)
		if err != nil {
			log.Println("NUM_MESSAGES is not a number")
			config.NumMessages = 1000
		}
	}

	schemaRegistry, ok := os.LookupEnv("SCHEMA_REGISTRY")
	if !ok {
		log.Println("SCHEMA_REGISTRY is not present")
		config.SchemaRegistry = "http://localhost:8081"
	} else {
		config.SchemaRegistry = schemaRegistry
	}

	return config
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
