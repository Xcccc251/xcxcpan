package Kafka

import (
	"XcxcPan/common/define"
	"XcxcPan/common/fileUtils"
	"os"

	"XcxcPan/common/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
	"time"
)

var (
	consumerKafkaBroker = "1.94.166.62:9092" // Kafka 地址
	consumerGroupID     = "group1"           // 消费者组 ID
)

func StartConsumer_Del() error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{consumerKafkaBroker},
		Topic:          define.KAFKA_DEL_TOPIC,
		GroupID:        consumerGroupID,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: 0,
	})

	defer reader.Close()

	log.Printf("Kafka consumer started, listening on topic %s\n", define.KAFKA_DEL_TOPIC)

	for {
		time.Sleep(time.Minute * 5)
		msg, err := reader.ReadMessage(context.Background())
		if time.Now().Before(msg.Time) {
			continue
		}
		if err != nil {
			return err
		}
		log.Printf("Received message:value=%s", string(msg.Value))
		var kafkaDelMessage define.KafkaDelMessage
		delMessageJson := msg.Value
		var delIds []string
		json.Unmarshal(delMessageJson, kafkaDelMessage)
		message := kafkaDelMessage.Message
		startTime := kafkaDelMessage.StartTime
		json.Unmarshal(message, delIds)
		for _, fileId := range delIds {
			var file models.File
			models.Db.Model(new(models.File)).
				Where("id = ?", fileId).Find(&file)

			if file.DelFlag == define.USING || file.RecoveryTime != startTime {
				continue
			}

			models.Db.Model(new(models.File)).Where("id = ?", fileId).Delete(&models.File{})

			if file.ChunkPrefix == "" {
				continue
			}

			var count int64
			models.Db.Model(new(models.File)).
				Where("chunk_prefix = ?", file.ChunkPrefix).
				Where("id != ?", fileId).
				Where("del_flag != ?", define.DEL).Count(&count)
			if count > 0 {
				continue
			}
			splitChunkPrefix := strings.Split(file.ChunkPrefix, "_")
			fmt.Println("删除切片")
			fileUtils.DelFileChunks(splitChunkPrefix[1], splitChunkPrefix[0])

			if file.FileCategory == define.GetCategoryCodeByCategory(define.VIDEO) {
				path := define.FILE_DIR + "/" + splitChunkPrefix[0] + "/" + splitChunkPrefix[1]
				os.RemoveAll(path)
			} else if file.FileCategory == define.GetCategoryCodeByCategory(define.IMAGE) {
				path := define.FILE_DIR + "/" + splitChunkPrefix[0] + "/" + splitChunkPrefix[1]
				os.RemoveAll(path)
			}

			models.RDb.Del(context.Background(), define.REDIS_USER_SPACE+file.UserId)
			models.RDb.Del(context.Background(), define.REDIS_CHUNK+splitChunkPrefix[0]+":"+splitChunkPrefix[1])

		}

		if err := reader.CommitMessages(context.Background(), msg); err != nil {
			log.Printf("Failed to commit message: %v", err)

		}
	}
}
