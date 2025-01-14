package helper

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/define"
	"XcxcPan/common/fileServerClient_gRPC"
	"XcxcPan/common/redisUtil"
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
)

func MergeChunks(sliceMap map[int][]byte) ([]byte, error) {
	// 获取切片的顺序索引
	var keys []int
	for k := range sliceMap {
		keys = append(keys, k)
	}

	// 按照索引排序
	sort.Ints(keys)

	// 创建一个 buffer 存储合并后的数据
	var buffer bytes.Buffer
	for _, k := range keys {
		data, ok := sliceMap[k]
		if !ok {
			return nil, fmt.Errorf("missing slice at index %d", k)
		}
		buffer.Write(data)
	}

	return buffer.Bytes(), nil
}

func GetSliceMap(userId string, fileId string) map[int][]byte {
	hashInt := redisUtil.GetHashInt(define.REDIS_CHUNK + userId + ":" + fileId)
	var dataMap = map[int][]byte{}
	wg := sync.WaitGroup{}
	var lock sync.Mutex
	for chunkIndex, serverId := range hashInt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := fileServerClient_gRPC.GetClientById(serverId)
			rsp, err := client.DownloadChunk(context.Background(), &XcXcPanFileServer.DownloadChunkRequest{
				FileName: userId + "_" + fileId + "_" + strconv.Itoa(chunkIndex),
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
	return dataMap
}

func MergeChunksToFile(sliceMap map[int][]byte) (*os.File, error) {
	// 获取切片的顺序索引
	var keys []int
	for k := range sliceMap {
		keys = append(keys, k)
	}

	// 按照索引排序
	sort.Ints(keys)

	// 创建一个 buffer 存储合并后的数据
	var buffer bytes.Buffer
	for _, k := range keys {
		data, ok := sliceMap[k]
		if !ok {
			return nil, fmt.Errorf("missing slice at index %d", k)
		}
		buffer.Write(data)
	}
	file, err := BytesToFile(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	file.Seek(0, 0)
	return file, nil

}
