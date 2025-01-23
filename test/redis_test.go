package test

import (
	"XcxcPan/common/redisUtil"
	"testing"
)

func TestHset(t *testing.T) {
	redisUtil.SetHash("xcxcpan:test", 5, 5)

	hashInt := redisUtil.GetChunkMap("xcxcpan:test")
	for k, v := range hashInt {
		t.Log(k, v)
	}

}
