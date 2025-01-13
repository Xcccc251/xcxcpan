package Server_MinIO

import (
	"XcxcPan/Server/XcXcPanFileServer/define"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func InitMinioClient_Server1() *minio.Client {
	client, err := minio.New(define.Server1_endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(define.Server1_accessKey, define.Server1_secretKey, ""),
	})
	if err != nil {

		log.Fatalf("Error creating MinIO client1: %v", err)
	}
	return client
}

func InitMinioClient_Server2() *minio.Client {
	client, err := minio.New(define.Server2_endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(define.Server2_accessKey, define.Server2_secretKey, ""),
	})
	if err != nil {

		log.Fatalf("Error creating MinIO client2: %v", err)
	}
	return client
}
func DelChunk(chunk_id string, serverId int) error {
	var minioClient *minio.Client
	if serverId == 1 {
		minioClient = InitMinioClient_Server1()
	} else {
		minioClient = InitMinioClient_Server2()
	}

	err := minioClient.RemoveObject(context.Background(), define.Chunk_bucketName, chunk_id, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}
func UploadChunk(chunk_id string, file *os.File, serverId int) (err error) {
	var minioClient *minio.Client
	if serverId == 1 {
		minioClient = InitMinioClient_Server1()
	} else {
		minioClient = InitMinioClient_Server2()
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	objectName := chunk_id
	contentType := "application/octet-stream"
	uploadInfo, err := minioClient.PutObject(context.Background(), define.Chunk_bucketName, objectName, file, fileInfo.Size(),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, uploadInfo.Size)
	return nil
}

func CheckChunkExists(chunk_id string, serverId int) bool {
	var minioClient *minio.Client
	if serverId == 1 {
		minioClient = InitMinioClient_Server1()
	} else {
		minioClient = InitMinioClient_Server2()
	}
	_, err := minioClient.StatObject(context.Background(), define.Chunk_bucketName, chunk_id, minio.StatObjectOptions{})
	if err != nil {
		return false
	}
	return true
}
