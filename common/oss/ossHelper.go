package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

var client *oss.Client
var bucketName = "hmleadnewsforxcccc"

func handleError(err error) {
	log.Fatalf("Error: %v", err)
}

func init() {

	var endpoint = "https://oss-cn-shanghai.aliyuncs.com"

	if endpoint == "" || bucketName == "" {
		log.Fatal("Please set yourEndpoint and bucketName.")
	}
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		handleError(err)
	}
	clientOptions := []oss.ClientOption{oss.SetCredentialsProvider(&provider)}
	clientOptions = append(clientOptions, oss.Region("cn-shanghai"))
	clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))
	client, err = oss.New(endpoint, "", "", clientOptions...)

	if err != nil {
		handleError(err)
	}
	// 输出客户端信息。
	log.Printf("Client: %#v\n", client)
}

func DownloadFile(objectName, downloadedFileName string) error {
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 下载文件。
	err = bucket.GetObjectToFile(objectName, downloadedFileName)
	if err != nil {
		return err
	}

	// 文件下载成功后，记录日志。
	log.Printf("File downloaded successfully to %s", downloadedFileName)
	return nil
}

func UploadFile(objectName string, localFileName io.Reader) (string, error) {
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	err = bucket.PutObject(objectName, localFileName)
	if err != nil {
		return "", err
	}
	// 上传文件。

	// 文件上传成功后，记录日志。
	log.Printf("File uploaded successfully to %s/%s", bucketName, objectName)
	path := "https://" + bucketName + "." + "oss-cn-shanghai.aliyuncs.com" + "/" + objectName
	return path, nil
}

func DelFile(objectName string) {
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return
	}
	err = bucket.DeleteObject(objectName)
	if err != nil {
		return
	}
}
