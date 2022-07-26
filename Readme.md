# Kafka & Schema Registry - protobuf

Example of using Kafka and [Schema Registry](https://docs.confluent.io/platform/current/schema-registry/serdes-develop/index.html) with [Protobuf](https://developers.google.com/protocol-buffers/).

[Confluent Kafka Go client](github.com/confluentinc/confluent-kafka-go) **v1.9.1** includes support for Schema Registry and Protobuf.

## Setup 

```bash
brew install protobuf
brew install protoc-gen-go
```

## Generate protobuf code (pb.go) 
 
```bash
protoc proto/*.proto -I. --go_out=:.
``` 

Go Protofub tips: https://jbrandhorst.com/post/go-protobuf-tips/


## Kafka and Schema Registry

Run Kafka and Schema Registry locally:

```sh
docker-compose up -d 
```

### Produce messages

```sh 
go run -mod=readonly ./go-producer/main.go
```

### Consume messages 

```sh 
go run -mod=readonly ./go-consumer/main.go
```

#### Golang Workspaces 

https://go.dev/ref/mod#workspaces

https://go.dev/doc/tutorial/workspaces

