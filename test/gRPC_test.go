package test

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/fileServerClient_gRPC"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"os"

	"testing"
)

func TestGRPC(t *testing.T) {
	client := fileServerClient_gRPC.ClientOfFileServer1
	file, err := os.Open("test_image_1.jpg")
	if err != nil {
		t.Error(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}
	rsp, err := client.UploadFile(context.Background(), &XcXcPanFileServer.UploadFileRequest{
		Data:     data,
		FileName: "test_image_1.png",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rsp.Message)
}
