package main

import (
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"os"
	"time"

	"github.com/mcolomerc/kafkasr/proto/model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
)

type Config struct {
	BootstrapServers string
	Topic            string
	SchemaRegistry   string
}

func main() {
	defer TimeTrack(time.Now(), "consumer")
	config := readEnvironment()
	fmt.Println(config)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  config.BootstrapServers,
		"group.id":           "protobuf-consumer",
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(config.SchemaRegistry))

	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	deser, err := protobuf.NewDeserializer(client, serde.ValueSerde, protobuf.NewDeserializerConfig())

	if err != nil {
		fmt.Printf("Failed to create deserializer: %s\n", err)
		os.Exit(1)
	}

	// Register the Protobuf type so that Deserialize can be called.
	// An alternative is to pass a pointer to an instance of the Protobuf type
	// to the DeserializeInto method. com.google.protobuf.DynamicMessage
	// deser.ProtoRegistry.RegisterMessage()
	deser.ProtoRegistry.RegisterMessage((&model.Person{}).ProtoReflect().Type())

	err = c.SubscribeTopics([]string{config.Topic}, nil)

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				value, err := deser.Deserialize(*e.TopicPartition.Topic, e.Value) // Deserialize the message.
				if err != nil {
					fmt.Printf("Failed to deserialize payload: %s\n", err)
				} else {
					fmt.Printf("%% Message on %s:\n%+v\n", e.TopicPartition, value)
				}
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()

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
