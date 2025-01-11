package minIO

import (
	"XcxcPan/common/define"
	"XcxcPan/common/helper"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"os"
	"path"
)

var minioClient = InitMinioClient()

func InitMinioClient() *minio.Client {
	client, err := minio.New(define.Server1_endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(define.Server1_accessKey, define.Server1_secretKey, ""),
	})
	if err != nil {

		log.Fatalf("Error creating MinIO client: %v", err)
	}
	return client
}

//	func UploadMP4(fileName string, file *os.File) (finalUrl string, err error) {
//		ext := path.Ext(fileName)
//		fileInfo, err := file.Stat()
//		if err != nil {
//			return "", err
//		}
//		objectName := helper.GetUUID() + ext
//		contentType := "video/mp4"
//		uploadInfo, err := minioClient.PutObject(context.Background(), bucketName, objectName, file, fileInfo.Size(),
//			minio.PutObjectOptions{ContentType: contentType})
//		if err != nil {
//			return "", err
//		}
//		log.Printf("Successfully uploaded %s of size %d\n", objectName, uploadInfo.Size)
//		return "http://127.0.0.1:9001/xcxcaudio/" + objectName, nil
//	}
//
//	func DelObject(fileName string) error {
//		err := minioClient.RemoveObject(context.Background(), bucketName, fileName, minio.RemoveObjectOptions{})
//		if err != nil {
//			log.Fatalln(err)
//			return err
//		}
//		return nil
//	}
func UploadImage(fileName string, file *os.File) (finalUrl string, err error) {
	ext := path.Ext(fileName)
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	objectName := helper.GetUUID() + ext
	contentType := "image/jpg"
	uploadInfo, err := minioClient.PutObject(context.Background(), define.Image_bucketName, objectName, file, fileInfo.Size(),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, uploadInfo.Size)
	return "http://127.0.0.1:9001/xcxcpanimage/" + objectName, nil
}

func DownloadImage(fileName string) (file *os.File, err error) {
	object, err := minioClient.GetObject(context.Background(), define.Image_bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	tempFile, err := os.CreateTemp("", "downloaded-image-*.jpg")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(tempFile, object)
	if err != nil {
		os.Remove(tempFile.Name())
		tempFile.Close()
		return nil, err
	}
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		os.Remove(tempFile.Name())
		tempFile.Close()
		return nil, err
	}
	return tempFile, nil
}

func CheckAvatarExists(fileName string) bool {
	_, err := minioClient.StatObject(context.Background(), define.Image_bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return false
	}
	return true
}
