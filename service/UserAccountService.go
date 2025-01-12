package service

import (
	"XcxcPan/common/define"
	"XcxcPan/common/helper"
	"XcxcPan/common/imageCode"
	"XcxcPan/common/minIO"
	"XcxcPan/common/models"
	"XcxcPan/common/response"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

func CheckCode(c *gin.Context) {
	typeOfCode, _ := strconv.Atoi(c.Query("type"))
	imgCode := imageCode.NewCreateImageCode()
	var buf bytes.Buffer
	imgCode.Write(&buf)
	fmt.Println(imgCode.Code)
	session := sessions.Default(c)
	if typeOfCode == 0 {
		session.Set(define.CHECK_CODE_KEY, imgCode.Code)
	} else {
		session.Set(define.CHECK_CODE_KEY_EMAIL, imgCode.Code)
	}

	err := session.Save()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "服务器错误",
		})
		return
	}
	c.Data(200, "image/png", buf.Bytes())
	return
}

func SendEmailCode(c *gin.Context) {
	email := c.PostForm("email")
	checkCode := c.PostForm("checkCode")
	//typeOfCode, _ := strconv.Atoi(c.PostForm("type"))
	if !helper.IsValidEmail(email) || email == "" {
		response.ResponseFailWithData(c, 0, "邮箱格式错误")
		return
	}
	session := sessions.Default(c)
	if checkCode == "" {
		response.ResponseFailWithData(c, 0, "验证码错误")
		return
	}
	code := session.Get(define.CHECK_CODE_KEY_EMAIL).(string)
	fmt.Println("code:", code)
	fmt.Println("checkCode:", checkCode)
	if code != checkCode {
		response.ResponseFailWithData(c, 0, "验证码错误")
		return
	}
	codeOfEmail := helper.GetRandNumber()
	helper.SendCode(email, codeOfEmail)
	models.RDb.Set(context.Background(), define.CODE_EMAIL+email, codeOfEmail, 15*60*time.Second)
	response.ResponseOK(c)
	return

}

func Register(c *gin.Context) {
	email := c.PostForm("email")
	emailCode := c.PostForm("emailCode")
	nickName := c.PostForm("nickName")
	checkCode := c.PostForm("checkCode")
	password := c.PostForm("password")
	session := sessions.Default(c)
	if checkCode == "" || checkCode != session.Get(define.CHECK_CODE_KEY).(string) {
		response.ResponseFailWithData(c, 0, "图片验证码错误")
		return
	}
	if !helper.IsValidEmail(email) || email == "" {
		response.ResponseFailWithData(c, 0, "邮箱格式错误")
		return
	}
	var count int64
	models.Db.Model(new(models.User)).Where("email = ?", email).Count(&count)
	if count > 0 {
		response.ResponseFailWithData(c, 0, "邮箱已存在")
		return
	}
	models.Db.Model(new(models.User)).Where("nick_name = ?", nickName).Count(&count)
	if count > 0 {
		response.ResponseFailWithData(c, 0, "昵称已存在")
		return
	}
	if emailCode == "" || nickName == "" || password == "" {
		response.ResponseFailWithData(c, 0, "参数错误")
		return
	}

	code, err := models.RDb.Get(context.Background(), define.CODE_EMAIL+email).Result()
	if err != nil || code != emailCode {
		response.ResponseFailWithData(c, 0, "邮箱验证码错误,请重试")
		return
	}
	password = helper.GetMd5(password)
	user := models.User{
		UserId:     helper.GetRandomStr(10),
		Email:      email,
		NickName:   nickName,
		Password:   password,
		Status:     define.USER_STATUS_ENABLE,
		TotalSpace: define.DEFAULT_TOTAL_SPACE,
	}
	err = models.Db.Create(&user).Error
	if err != nil {
		response.ResponseFailWithData(c, 0, "注册失败")
		return
	}
	models.RDb.Del(context.Background(), define.CODE_EMAIL+email)
	response.ResponseOK(c)
	return

}

func Login(c *gin.Context) {
	var userInfo models.User
	var count int64
	email := c.PostForm("email")
	password := c.PostForm("password")
	checkCode := c.PostForm("checkCode")
	session := sessions.Default(c)
	if checkCode == "" || checkCode != session.Get(define.CHECK_CODE_KEY).(string) {
		response.ResponseFailWithData(c, 0, "图片验证码错误")
		return
	}
	if email == "" || password == "" {
		response.ResponseFailWithData(c, 0, "参数错误")
		return
	}
	db := models.Db.Model(new(models.User)).Where("email = ?", email)
	db.Count(&count)
	if count == 0 {
		response.ResponseFailWithData(c, 0, "账号或密码错误")
		return
	}
	db.Find(&userInfo)
	if userInfo.Status == define.USER_STATUS_DISABLE {
		response.ResponseFailWithData(c, 0, "账号已被禁用")
		return
	}
	if userInfo.Password != password {
		response.ResponseFailWithData(c, 0, "账号或密码错误")
		return
	}
	db.Update("last_login_time", models.MyTime(time.Now()))
	userLoginDto := models.UserLoginDto{
		UserId:   userInfo.UserId,
		NickName: userInfo.NickName,
	}
	if userInfo.Email == define.ADMIN_EMAIL {
		userLoginDto.IsAdmin = true
	}
	userInfoJson, _ := json.Marshal(&userLoginDto)
	models.RDb.Set(context.Background(), define.REDIS_USER_INFO+userInfo.UserId, userInfoJson, define.EXPIRE_DAY)

	session.Set(define.SESSION_USER_ID, userInfo.UserId)
	err := session.Save()
	if err != nil {
		response.ResponseFailWithData(c, 0, "服务器错误")
		return
	}

	var userSpaceDto models.UserSpaceDto
	userSpaceDto.UseSpace = userInfo.UseSpace
	userSpaceDto.TotalSpace = userInfo.TotalSpace
	//todo 查询已使用的空间大小
	userSpaceJson, _ := json.Marshal(userSpaceDto)
	models.RDb.Set(context.Background(), define.REDIS_USER_SPACE+userInfo.UserId, userSpaceJson, define.EXPIRE_DAY)

	response.ResponseOKWithData(c, userLoginDto)
	return

}

func ResetPassword(c *gin.Context) {
	email := c.PostForm("email")
	emailCode := c.PostForm("emailCode")
	checkCode := c.PostForm("checkCode")
	password := c.PostForm("password")
	session := sessions.Default(c)
	if checkCode == "" || checkCode != session.Get(define.CHECK_CODE_KEY).(string) {
		response.ResponseFailWithData(c, 0, "图片验证码错误")
		return
	}
	if !helper.IsValidEmail(email) || email == "" {
		response.ResponseFailWithData(c, 0, "邮箱格式错误")
		return
	}
	var user models.User
	if models.Db.Model(new(models.User)).Where("email = ?", email).First(&user).RowsAffected == 0 {
		response.ResponseFailWithData(c, 0, "账号不存在")
		return
	}
	if emailCode == "" || password == "" {
		response.ResponseFailWithData(c, 0, "参数错误")
		return
	}

	code, err := models.RDb.Get(context.Background(), define.CODE_EMAIL+email).Result()
	if err != nil || code != emailCode {
		response.ResponseFailWithData(c, 0, "邮箱验证码错误,请重试")
		return
	}

	passwordMd5 := helper.GetMd5(password)
	models.Db.Model(new(models.User)).Where("email = ?", email).Update("password", passwordMd5)
	models.RDb.Del(context.Background(), define.CODE_EMAIL+email)
	response.ResponseOK(c)
	return
}

func GetAvatar(c *gin.Context) {
	userId := c.Param("userId")
	exists := minIO.CheckAvatarExists(userId + ".jpg")
	var file *os.File
	var err error
	if !exists {
		//如果minio中不存在，则获取默认头像，需要提前将默认头像上传至minio
		file, err = minIO.DownloadImage(define.DEFAULT_AVATAR_NAME)
		if err != nil {
			response.ResponseFailWithData(c, 0, "获取头像失败")
		}
	} else {
		file, err = minIO.DownloadImage(userId + ".jpg")
		if err != nil {
			response.ResponseFailWithData(c, 0, "获取头像失败")
		}
	}
	data, err := helper.FileToBytes(file)
	if err != nil {
		response.ResponseFailWithData(c, 0, "获取头像失败")
	}
	c.Data(200, "image/jpg", data)
	return
}

func GetUserInfo(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get(define.SESSION_USER_ID).(string)
	var userLoginDto models.UserLoginDto
	result, _ := models.RDb.Get(context.Background(), define.REDIS_USER_INFO+userId).Result()
	json.Unmarshal([]byte(result), &userLoginDto)
	response.ResponseOKWithData(c, userLoginDto)
	return

}

func GetUseSpace(c *gin.Context) {
	//session := sessions.Default(c)
	//userInfo := session.Get(define.USER_INFO).([]byte)
	//var userLoginDto models.UserLoginDto
	//json.Unmarshal(userInfo, &userLoginDto)
	session := sessions.Default(c)
	userId := session.Get(define.SESSION_USER_ID).(string)
	response.ResponseOKWithData(c, getUserUseSpace(userId))
	return
}

func getUserUseSpace(userId string) models.UserSpaceDto {
	result, _ := models.RDb.Get(context.Background(), define.REDIS_USER_SPACE+userId).Result()
	var userSpaceDto models.UserSpaceDto
	json.Unmarshal([]byte(result), &userSpaceDto)
	//todo 查询已使用的空间大小
	fmt.Println(define.REDIS_USER_SPACE + userId)
	fmt.Println(userSpaceDto)
	return userSpaceDto
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	userId, _ := c.Get("userId")
	session.Clear()
	err := session.Save()
	if err != nil {
		response.ResponseFailWithData(c, 0, "服务器错误")
		return
	}
	models.RDb.Del(context.Background(), define.REDIS_USER_INFO+userId.(string))
	models.RDb.Del(context.Background(), define.REDIS_USER_SPACE+userId.(string))
	response.ResponseOK(c)
	return
}

func UpdateUserAvatar(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get(define.SESSION_USER_ID).(string)
	avatar, _ := c.FormFile("avatar")
	file, _ := avatar.Open()
	finalFile, _ := helper.SaveMultipartFile(file)
	minIO.UploadUserAvatar(avatar.Filename, userId, finalFile)
	go func() {
		models.Db.Model(new(models.User)).Where("id = ?", userId).Update("qq_avatar", "")
	}()
	go func() {
		var userLoginDto models.UserLoginDto
		result, _ := models.RDb.Get(context.Background(), define.REDIS_USER_INFO+userId).Result()
		json.Unmarshal([]byte(result), &userLoginDto)
		userLoginDto.Avatar = ""
		userInfoJson, _ := json.Marshal(&userLoginDto)
		models.RDb.Set(context.Background(), define.REDIS_USER_INFO+userId, userInfoJson, define.EXPIRE_DAY)
	}()
	response.ResponseOK(c)
	return
}
