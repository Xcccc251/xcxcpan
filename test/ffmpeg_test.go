package test

import (
	"XcxcPan/service"
	"testing"
)

func TestCutVideo(t *testing.T) {
	err := service.CutFileForVideo("E:\\xcxcpan_file\\dir\\2025-01\\FbjBFbLoLy\\UMZHbZcOfJOrAfheLGubHoCssmEqGwVW.mp4")
	if err != nil {
		t.Error(err)
	}
}

func TestImage(t *testing.T) {
	path, err := service.CreateThumbnailForVideo("E:\\xcxcpan_file\\dir\\2025-01\\FbjBFbLoLy\\test\\input.mp4")
	if err != nil {
		t.Error(err)
	}
	t.Log(path)
}
