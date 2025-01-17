package Kafka

import (
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

// Kafka
const (
	kafkaBroker = "1.94.166.62:9092"
)

// 生产者逻辑
func ProduceMessage(topic string, message []byte) error {
	writer := kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%d", time.Now().UnixNano())),
			Value: message,
		},
	)

	return err
}

func ProduceMessageWithTime(topic string, message []byte, startTime models.MyTime, delayDuration time.Duration) error {
	writer := kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()
	var kafkaDelMessage define.KafkaDelMessage
	kafkaDelMessage.Message = message
	kafkaDelMessage.StartTime = startTime
	kafkaDelMessage.DelayDuration = delayDuration
	kafkaDelMessageJson, _ := json.Marshal(&kafkaDelMessage)

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: kafkaDelMessageJson,
			Time:  time.Time(startTime).Add(delayDuration),
		},
	)

	return err
}
