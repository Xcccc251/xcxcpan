package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UserId        string `gorm:"column:id" json:"userId"`
	NickName      string `gorm:"column:nick_name" json:"nickName"`
	Email         string `gorm:"column:email" json:"email"`
	OpenId        string `gorm:"column:qq_open_id" json:"qqOpenId"`
	Avatar        string `gorm:"column:qq_avatar" json:"qqAvatar"`
	Password      string `gorm:"column:password" json:"password"`
	JoinTime      MyTime `gorm:"column:join_time" json:"joinTime"`
	LastLoginTime MyTime `gorm:"column:last_login_time" json:"lastLoginTime"`
	Status        int    `gorm:"column:status;default:1" json:"status"`
	UseSpace      int    `gorm:"column:use_space" json:"useSpace"`
	TotalSpace    int    `gorm:"column:total_space" json:"totalSpace"`
}
type UserLoginDto struct {
	gorm.Model
	UserId   string `gorm:"column:id" json:"userId"`
	NickName string `gorm:"column:nick_name" json:"nickName"`
	IsAdmin  bool   `gorm:"column:is_admin" json:"isAdmin"`
	Avatar   string `gorm:"column:qq_avatar" json:"qqAvatar"`
}
type UserSpaceDto struct {
	UseSpace   int `json:"useSpace"`
	TotalSpace int `json:"totalSpace"`
}

func (User) TableName() string {
	return "user_info"
}
func (UserLoginDto) TableName() string {
	return "user_info"
}
func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	u.JoinTime = MyTime(time.Now())
	return
}
