package main

import (
	"XcxcPan/common/define"
	hashRing "XcxcPan/common/hash"
	"XcxcPan/router"
)

func main() {
	r := router.Router()
	//初始哈希环，添加两个服务节点
	hashRing.Hash.Add(define.SERVER_1, define.SERVER_2)
	r.Run(":7090")

}
