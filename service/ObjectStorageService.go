package service

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/define"
	"XcxcPan/common/helper"
	"XcxcPan/common/redisUtil"
	"XcxcPan/fileServerClient_gRPC"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
)

func TestVideo(c *gin.Context) {
	hashInt := redisUtil.GetHashInt(define.REDIS_CHUNK + "FbjBFbLoLy" + ":" + "UMZHbZcOfJOrAfheLGubHoCssmEqGwVW")
	var dataMap = map[int][]byte{}
	wg := sync.WaitGroup{}
	var lock sync.Mutex
	for chunkIndex, serverId := range hashInt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := fileServerClient_gRPC.GetClientById(serverId)
			rsp, err := client.DownloadChunk(context.Background(), &XcXcPanFileServer.DownloadChunkRequest{
				FileName: "FbjBFbLoLy" + "_" + "UMZHbZcOfJOrAfheLGubHoCssmEqGwVW" + "_" + strconv.Itoa(chunkIndex),
				Server:   int64(serverId),
			})
			//加锁保护map
			lock.Lock()
			dataMap[chunkIndex] = rsp.Data
			lock.Unlock()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	wg.Wait()
	data, err := helper.MergeChunks(dataMap)
	if err != nil {
		fmt.Println(err)
	}

	// 获取客户端的 Range 请求头
	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		// 如果没有 Range 请求头，返回整个视频
		c.Data(200, "video/mp4", data)
		return
	}
	// 解析 Range 请求头
	var start, end int64
	fileSize := int64(len(data))
	fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if end == 0 || end >= fileSize {
		end = fileSize - 1
	}

	// 返回部分数据
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Length", fmt.Sprintf("%d", end-start+1))
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Status(206) // Partial Content

	// 写入指定范围的数据
	c.Writer.Write(data[start : end+1])

}
