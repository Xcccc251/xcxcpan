package test

import (
	"XcxcPan/service"
	"testing"
)

func TestMinio(t *testing.T) {
	//service.DelFileChunks("xkYOvWxyUwvIQsRCIRTpXxpIadBFaRxx", "FbjBFbLoLy")
	err := service.TransferFile("UMZHbZcOfJOrAfheLGubHoCssmEqGwVW")
	if err != nil {
		t.Error(err)
	}
}
