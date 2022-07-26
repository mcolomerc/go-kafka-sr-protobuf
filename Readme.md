# Kafka & Schema Registry - protobuf

Example of using Kafka and [Schema Registry](https://docs.confluent.io/platform/current/schema-registry/serdes-develop/index.html) with [Protobuf](https://developers.google.com/protocol-buffers/).

[Confluent Kafka Go client](github.com/confluentinc/confluent-kafka-go) includes support for Schema Registry.
 
 * [**v1.9.1**](https://github.com/confluentinc/confluent-kafka-go/releases/tag/v1.9.1): *Schema Registry support for Avro Generic and Specific, Protocol Buffers and JSON Schema*

## Setup 

```bash
brew install protobuf
brew install protoc-gen-go
```

## Generate protobuf code (pb.go) 
 
```bash
protoc proto/*/*.proto -I. --go_out=:./proto
```  

Go Protofub tips: https://jbrandhorst.com/post/go-protobuf-tips/


## Kafka and Schema Registry

Run Zookeeper, Kafka, Confluent Control Center and Schema Registry locally:

```sh
docker-compose up -d 
```

### Produce messages

```sh 
go run -mod=readonly ./go-producer/main.go
```

Environment:

* BOOTSTRAP_SERVERS (default: localhost:9092)
* TOPIC (default: sr-topic-pb-test)
* NUM_MESSAGES (default: 1000)
* SCHEMA_REGISTRY (default: http://localhost:8081)

### Consume messages 

```sh 
go run -mod=readonly ./go-consumer/main.go
```

Environment:

* BOOTSTRAP_SERVERS (default: localhost:9092)
* TOPIC (default: sr-topic-pb-test)
* SCHEMA_REGISTRY (default: http://localhost:8081)

#### Golang Workspaces 

https://go.dev/ref/mod#workspaces

https://go.dev/doc/tutorial/workspaces

