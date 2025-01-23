package test

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/define"
	"XcxcPan/common/helper"
	"XcxcPan/common/models"
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
	hashInt := redisUtil.GetChunkMap(define.REDIS_CHUNK + "FbjBFbLoLy" + ":" + "XZJYWUTetSPgvQgSwrVjlxEfxAdDvnNZ")
	var dataMap = map[int][]byte{}
	wg := sync.WaitGroup{}
	now := time.Now()
	for chunkIndex, serverId := range hashInt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(chunkIndex)
			client := fileServerClient_gRPC.GetClientById(serverId)
			ctx, _ := context.WithTimeout(context.Background(), time.Minute*30)
			rsp, err := client.DownloadChunk(ctx, &XcXcPanFileServer.DownloadChunkRequest{
				FileName: "FbjBFbLoLy" + "_" + "XZJYWUTetSPgvQgSwrVjlxEfxAdDvnNZ" + "_" + strconv.Itoa(chunkIndex),
				Server:   int64(serverId),
			})
			if err != nil {
				t.Error(err)
			}
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
	hashInt := redisUtil.GetChunkMap(define.REDIS_CHUNK + "FbjBFbLoLy" + ":" + "fdamHkpjhuaVMArDmabexOzKhDIKjUDR")
	var dataMap = map[int][]byte{}

	now := time.Now()
	for chunkIndex, serverId := range hashInt {
		client := fileServerClient_gRPC.GetClientById(serverId)
		rsp, err := client.DownloadChunk(context.Background(), &XcXcPanFileServer.DownloadChunkRequest{
			FileName: "FbjBFbLoLy" + "_" + "fdamHkpjhuaVMArDmabexOzKhDIKjUDR" + "_" + strconv.Itoa(chunkIndex),
			Server:   int64(serverId),
		})
		fmt.Println(len(rsp.Data))
		fmt.Println("2")
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

func TestDelChunks(t *testing.T) {
	chunkIdServerIdMap := redisUtil.GetChunkMap(define.REDIS_CHUNK + "FbjBFbLoLy" + ":" + "XZJYWUTetSPgvQgSwrVjlxEfxAdDvnNZ")
	for chunkId, serverId := range chunkIdServerIdMap {
		client := fileServerClient_gRPC.GetClientById(serverId)
		ctx, _ := context.WithTimeout(context.Background(), time.Minute*30)
		_, err := client.DelChunk(ctx, &XcXcPanFileServer.DelChunkRequest{
			FileName: "FbjBFbLoLy" + "_" + "XZJYWUTetSPgvQgSwrVjlxEfxAdDvnNZ" + "_" + strconv.Itoa(chunkId),
			Server:   int64(serverId),
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	models.RDb.Del(context.Background(), define.REDIS_CHUNK+"FbjBFbLoLy"+":"+"XZJYWUTetSPgvQgSwrVjlxEfxAdDvnNZ")
}
