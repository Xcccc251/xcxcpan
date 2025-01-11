package test

import (
	"XcxcPan/common/minIO"
	"fmt"
	"testing"
)

func TestMinio(t *testing.T) {
	fmt.Println(minIO.CheckAvatarExists("user1.jpg"))

}
