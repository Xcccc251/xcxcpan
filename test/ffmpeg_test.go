package test

import (
	"XcxcPan/service"
	"testing"
)

func TestCutVideo(t *testing.T) {
	err := service.CutFileForVideo("E:/xcxcpan_file/dir/FbjBFbLoLy/gekloGWjOMXOvYgBixojAdMAYCsXsfjK.mp4")
	if err != nil {
		t.Error(err)
	}
}

func TestImage(t *testing.T) {
	err := service.CreateThumbnailForVideo("E:/xcxcpan_file/dir/FbjBFbLoLy/gekloGWjOMXOvYgBixojAdMAYCsXsfjK.mp4")
	if err != nil {
		t.Error(err)
	}
}
