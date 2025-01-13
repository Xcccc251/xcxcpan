package main

import (
	"XcxcPan/common/define"
	hashRing "XcxcPan/common/hash"
	"XcxcPan/router"
)

func main() {
	r := router.Router()
	hashRing.Hash.Add(define.SERVER_1, define.SERVER_2)
	r.Run(":7090")

}
