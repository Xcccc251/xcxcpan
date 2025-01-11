package helper

import (
	"XcxcPan/common/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func MyTimeToTimestamp(t models.MyTime) *timestamppb.Timestamp {
	return timestamppb.New(time.Time(t))
}
func TimestampToMyTime(t *timestamppb.Timestamp) models.MyTime {
	return models.MyTime(t.AsTime().Local())
}
