module github.com/mcolomerc/kafkasr/consumer

go 1.18

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jhump/protoreflect v1.12.0 // indirect
	google.golang.org/genproto v0.0.0-20220725144611-272f38e5d71b // indirect
)

require (
	github.com/confluentinc/confluent-kafka-go v1.9.1
	google.golang.org/protobuf v1.28.0 // indirect
)

require github.com/mcolomerc/kafkasr/proto v0.0.0-00010101000000-000000000000

replace github.com/mcolomerc/kafkasr/proto v0.0.0-00010101000000-000000000000 => ../proto
