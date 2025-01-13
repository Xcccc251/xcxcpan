package test

import (
	hashRing "XcxcPan/common/hash"
	"fmt"
	"testing"
)

func TestHashRing(t *testing.T) {
	hashRing.Hash.Add("node1", "node2")
	fmt.Println(hashRing.Hash.Get("test"))
	fmt.Println(hashRing.Hash.Get("test2"))
	fmt.Println(hashRing.Hash.Get("test3"))
	fmt.Println(hashRing.Hash.Get("test4"))
	fmt.Println(hashRing.Hash.Get("test"))
	fmt.Println(hashRing.Hash.Get("test"))
}
