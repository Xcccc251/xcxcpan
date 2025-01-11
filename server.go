package main

import (
	"XcxcPan/router"
)

func main() {
	r := router.Router()

	r.Run(":7090")

}
