package fileUtils

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"XcxcPan/common/redisUtil"
	"XcxcPan/fileServerClient_gRPC"
	"context"
	"fmt"
	"strconv"
)

func DelFileChunks(fileId string, userId string) {
	chunkIdServerIdMap := redisUtil.GetHashInt(define.REDIS_CHUNK + userId + ":" + fileId)
	for chunkId, serverId := range chunkIdServerIdMap {
		client := fileServerClient_gRPC.GetClientById(serverId)
		_, err := client.DelChunk(context.Background(), &XcXcPanFileServer.DelChunkRequest{
			FileName: userId + "_" + fileId + "_" + strconv.Itoa(chunkId),
			Server:   int64(serverId),
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	models.RDb.Del(context.Background(), define.REDIS_CHUNK+userId+":"+fileId)

}
