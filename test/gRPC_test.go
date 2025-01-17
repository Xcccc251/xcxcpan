package test

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/define"
	"XcxcPan/common/helper"
	"XcxcPan/common/redisUtil"
	"XcxcPan/fileServerClient_gRPC"
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestGRPC(t *testing.T) {
	//client := fileServerClient_gRPC.ClientOfFileServer1
	//file, err := os.Open("test_image_1.jpg")
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//data, err := io.ReadAll(file)
	//if err != nil {
	//	t.Error(err)
	//}
	//rsp, err := client.UploadFile(context.Background(), &XcXcPanFileServer.UploadFileRequest{
	//	Data:     data,
	//	FileName: "test_image_1.png",
	//})
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println(rsp.Message)
}

func TestDownloadChunk(t *testing.T) {
	hashInt := redisUtil.GetHashInt(define.REDIS_CHUNK + "FbjBFbLoLy" + ":" + "ZQbrDrLIJQBIPZjxIpDAqEGxzkaGCmqn")
	var dataMap = map[int][]byte{}
	wg := sync.WaitGroup{}
	now := time.Now()
	for chunkIndex, serverId := range hashInt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := fileServerClient_gRPC.GetClientById(serverId)
			rsp, err := client.DownloadChunk(context.Background(), &XcXcPanFileServer.DownloadChunkRequest{
				FileName: "FbjBFbLoLy" + "_" + "ZQbrDrLIJQBIPZjxIpDAqEGxzkaGCmqn" + "_" + strconv.Itoa(chunkIndex),
				Server:   int64(serverId),
			})
			dataMap[chunkIndex] = rsp.Data
			if err != nil {
				t.Error(err)
			}

		}()
	}
	wg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(now).Milliseconds())
	outputFile := "test.mp4"
	data, err := helper.MergeChunks(dataMap)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile(outputFile, data, 0666)

}

func TestDownloadChunk2(t *testing.T) {
	hashInt := redisUtil.GetHashInt(define.REDIS_CHUNK + "FbjBFbLoLy" + ":" + "ZQbrDrLIJQBIPZjxIpDAqEGxzkaGCmqn")
	var dataMap = map[int][]byte{}

	now := time.Now()
	for chunkIndex, serverId := range hashInt {
		client := fileServerClient_gRPC.GetClientById(serverId)
		rsp, err := client.DownloadChunk(context.Background(), &XcXcPanFileServer.DownloadChunkRequest{
			FileName: "FbjBFbLoLy" + "_" + "ZQbrDrLIJQBIPZjxIpDAqEGxzkaGCmqn" + "_" + strconv.Itoa(chunkIndex),
			Server:   int64(serverId),
		})
		dataMap[chunkIndex] = rsp.Data
		if err != nil {
			t.Error(err)
		}
	}

	end := time.Now()
	fmt.Println(end.Sub(now).Milliseconds())
	outputFile := "test2.mp4"
	data, err := helper.MergeChunks(dataMap)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile(outputFile, data, 0666)

}
