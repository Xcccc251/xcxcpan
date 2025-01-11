package oss

import (
	"XcxcVideo/common/helper"
	"mime/multipart"
	"path"
)

func UploadImage(file *multipart.FileHeader) (FinalPath string, err error) {
	filename := file.Filename
	ext := path.Ext(filename)
	objectName := helper.GetUUID() + ext
	filePath, err := file.Open()
	if err != nil {
		return "", err
	}
	uploadFilePath, err := UploadFile(objectName, filePath)
	if err != nil {
		return "", err
	}
	return uploadFilePath, nil

}
