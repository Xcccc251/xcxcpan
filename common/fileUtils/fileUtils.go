package fileUtils

import (
	"XcxcPan/StorageGroup"
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"XcxcPan/common/redisUtil"
	"context"
	"strconv"
)

func DelFileChunks(fileId string, userId string) {
	chunkMap := redisUtil.GetChunkMap(define.REDIS_CHUNK + userId + ":" + fileId)
	for chunkId, serverNode := range chunkMap {
		StorageGroup.Server.GrpcGetters[serverNode].Del(userId + "_" + fileId + "_" + strconv.Itoa(chunkId))
	}
	models.RDb.Del(context.Background(), define.REDIS_CHUNK+userId+":"+fileId)

}
