package helper

import (
	"bytes"
	"fmt"
	"sort"
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
