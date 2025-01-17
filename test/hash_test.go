package test

import (
	hashRing "XcxcPan/common/hash"
	"fmt"
	"testing"
)

func TestHashRing(t *testing.T) {
	hashRing.Hash.Add("node1", "node2")
	fmt.Println(hashRing.Hash.Get("FbjBFbLoLy_LaHjwIjuzGDhOxFsNlSVGVNiNRNJeCYq_0"))
}
