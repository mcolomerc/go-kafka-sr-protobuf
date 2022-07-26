module github.com/mcolomerc/kafkasr/producer

go 1.18

require (
	github.com/confluentinc/confluent-kafka-go v1.9.1
	google.golang.org/protobuf v1.28.0 // indirect
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/genproto v0.0.0-20220503193339-ba3ae3f07e29 // indirect
)

require (
	github.com/brianvoe/gofakeit/v6 v6.17.0
	github.com/mcolomerc/kafkasr/proto v0.0.0-00010101000000-000000000000
)

require github.com/jhump/protoreflect v1.12.0 // indirect

replace github.com/mcolomerc/kafkasr/proto v0.0.0-00010101000000-000000000000 => ../proto
