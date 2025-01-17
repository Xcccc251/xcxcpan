package models

import (
	"gorm.io/gorm"
	"time"
)

type Share struct {
	gorm.Model
	ShareId    string `gorm:"column:share_id" json:"shareId"`
	FileId     string `gorm:"column:file_id" json:"fileId"`
	UserId     string `gorm:"column:user_id" json:"userId"`
	ValidType  int    `gorm:"column:valid_type" json:"validType"`
	ExpireTime MyTime `gorm:"column:expire_time" json:"expireTime"`
	ShareTime  MyTime `gorm:"column:share_time" json:"shareTime"`
	Code       string `gorm:"column:code" json:"code"`
	ShowCount  int    `gorm:"column:show_count" json:"showCount"`
}
type ShareVo struct {
	gorm.Model
	ShareId    string `gorm:"column:share_id" json:"shareId"`
	FileId     string `gorm:"column:file_id" json:"fileId"`
	UserId     string `gorm:"column:user_id" json:"userId"`
	ValidType  int    `gorm:"column:valid_type" json:"validType"`
	ExpireTime MyTime `gorm:"column:expire_time" json:"expireTime"`
	ShareTime  MyTime `gorm:"column:share_time" json:"shareTime"`
	Code       string `gorm:"column:code" json:"code"`
	FileName   string `gorm:"column:file_name" json:"fileName"`
	FileCover  string `gorm:"column:file_cover" json:"fileCover"`
	FileType   int    `gorm:"column:file_type" json:"fileType"`
	Status     int    `gorm:"column:status" json:"status"`
	FolderType int    `gorm:"column:folder_type" json:"folderType"`
	ShowCount  int    `gorm:"column:show_count" json:"showCount"`
}
type ShareInfoVo struct {
	gorm.Model
	CurrentUser bool   `json:"currentUser"`
	FileId      string `gorm:"column:file_id" json:"fileId"`
	UserId      string `gorm:"column:user_id" json:"userId"`
	ExpireTime  MyTime `gorm:"column:expire_time" json:"expireTime"`
	ShareTime   MyTime `gorm:"column:share_time" json:"shareTime"`
	FileName    string `gorm:"column:file_name" json:"fileName"`
	NickName    string `gorm:"column:nick_name" json:"nickName"`
	Avatar      string `gorm:"column:avatar" json:"avatar"`
}
type SessionShareDto struct {
	ShareId     string `json:"shareId"`
	ShareUserId string `json:"shareUserId"`
	ExpireTime  MyTime `json:"expireTime"`
	FileId      string `json:"fileId"`
}

func (Share) TableName() string {
	return "file_share"
}
func (ShareVo) TableName() string {
	return "file_share"
}

func (s *Share) BeforeCreate(tx *gorm.DB) (err error) {
	s.ShareTime = MyTime(time.Now())
	return
}

func (s *Share) AddShare(tx *gorm.DB) (err error) {

	return tx.Create(s).Error
}
