package service

import (
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"XcxcPan/common/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetFileList(c *gin.Context) {
	categoryStr := c.PostForm("category")
	pageNo, _ := strconv.Atoi(c.PostForm("pageNo"))
	pageSize, _ := strconv.Atoi(c.PostForm("pageSize"))
	if pageNo == 0 {
		pageNo = define.DEFAULT_PAGE_NO
	}
	if pageSize == 0 {
		pageSize = define.DEFAULT_PAGE_SIZE
	}
	filePid, _ := strconv.Atoi(c.PostForm("filePid"))
	var category int
	var fileList []models.FileVo
	db := models.Db.Model(new(models.File))
	if define.ExistsCategory(categoryStr) {
		category = define.VIDEO_CATEGORY[categoryStr]
		db.Where("file_category = ?", category)
	}
	userId, _ := c.Get("userId")
	db.Where("user_id = ?", userId.(string))
	db.Where("del_flag = ?", define.USING)
	db.Where("file_pid = ?", filePid)
	db.Order("last_update_time desc")
	//db.Offset((pageNo - 1) * pageSize).Limit(pageSize)
	//db.Find(&fileList)
	fmt.Println("pageNo:", pageNo, "pageSize:", pageSize, "category:", category, "filePid:", filePid)
	response.ResponseOKWithData(c, models.QueryPageList(db, pageNo, pageSize, &fileList))
	return

}

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	fileName := c.PostForm("fileName")
	fileMd5 := c.PostForm("fileMd5")
	chunkIndex, _ := strconv.Atoi(c.PostForm("chunkIndex"))
	chunks, _ := strconv.Atoi(c.PostForm("chunks"))
	fileId := c.PostForm("fileId")
	filePid := c.PostForm("filePid")
	userId, _ := c.Get("userId")
	var userLoginDto models.UserLoginDto
	result, _ := models.RDb.Get(context.Background(), define.REDIS_USER_INFO+userId.(string)).Result()
	json.Unmarshal([]byte(result), &userLoginDto)

}
