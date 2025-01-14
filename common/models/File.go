package models

import (
	"gorm.io/gorm"
	"time"
)

type File struct {
	gorm.Model
	Id             string `gorm:"column:id" json:"fileId"`
	UserId         string `gorm:"column:user_id" json:"userId"`
	FileMd5        string `gorm:"column:file_md5" json:"fileMd5"`
	FilePid        string `gorm:"column:file_pid" json:"filePid"`
	FileSize       int    `gorm:"column:file_size" json:"fileSize"`
	FileName       string `gorm:"column:file_name" json:"fileName"`
	FileCover      string `gorm:"column:file_cover" json:"fileCover"`
	FilePath       string `gorm:"column:file_path" json:"filePath"`
	LastUpdateTime MyTime `gorm:"column:last_update_time" json:"lastUpdateTime"`
	FolderType     int    `gorm:"column:folder_type" json:"folderType"`
	FileCategory   int    `gorm:"column:file_category" json:"fileCategory"`
	FileType       int    `gorm:"column:file_type" json:"fileType"`
	Status         int    `gorm:"column:status" json:"status"`
	RecoveryTime   MyTime `gorm:"column:recovery_time" json:"recoveryTime"`
	DelFlag        int    `gorm:"column:del_flag;default:2" json:"delFlag"`
	ChunkPrefix    string `gorm:"column:chunk_prefix" json:"chunkPrefix"`
}
type FileVo struct {
	gorm.Model
	Id             string `gorm:"column:id" json:"fileId"`
	UserId         string `gorm:"column:user_id" json:"userId"`
	FileMd5        string `gorm:"column:file_md5" json:"fileMd5"`
	FilePid        string `gorm:"column:file_pid" json:"filePid"`
	FileSize       int64  `gorm:"column:file_size" json:"fileSize"`
	FileName       string `gorm:"column:file_name" json:"fileName"`
	FileCover      string `gorm:"column:file_cover" json:"fileCover"`
	FilePath       string `gorm:"column:file_path" json:"filePath"`
	LastUpdateTime MyTime `gorm:"column:last_update_time" json:"lastUpdateTime"`
	FolderType     int    `gorm:"column:folder_type" json:"folderType"`
	FileCategory   int    `gorm:"column:file_category" json:"fileCategory"`
	FileType       int    `gorm:"column:file_type" json:"fileType"`
	Status         int    `gorm:"column:status" json:"status"`
	RecoveryTime   MyTime `gorm:"column:recovery_time" json:"recoveryTime"`
	DelFlag        int    `gorm:"column:del_flag;default:2" json:"delFlag"`
	CreateTime     MyTime `gorm:"column:created_at" json:"createTime"`
	ChunkPrefix    string `gorm:"column:chunk_prefix" json:"chunkPrefix"`
}

type UploadResultDto struct {
	FileId string `json:"fileId"`
	Status string `json:"status"`
}

func (File) TableName() string {
	return "file_info"
}
func (FileVo) TableName() string {
	return "file_info"
}
func (f *File) BeforeUpdate(tx *gorm.DB) (err error) {
	f.LastUpdateTime = MyTime(time.Now())
	return
}
func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	f.LastUpdateTime = MyTime(time.Now())
	return
}
