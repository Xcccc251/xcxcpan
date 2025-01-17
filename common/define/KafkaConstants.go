package define

import (
	"XcxcPan/common/models"
	"time"
)

var KAFKA_DEL_TOPIC = "xcxc_pan_del"

var KAFKA_DEL_DURATION = time.Hour * 24 * 10

type KafkaDelMessage struct {
	Message       []byte
	StartTime     models.MyTime
	DelayDuration time.Duration
}
