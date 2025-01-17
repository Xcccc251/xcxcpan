package service

import (
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"XcxcPan/common/response"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"sync"
)

func LoadRecycleList(c *gin.Context) {
	pageNo, _ := strconv.Atoi(c.PostForm("pageNo"))
	pageSize, _ := strconv.Atoi(c.PostForm("pageSize"))
	if pageNo == 0 {
		pageNo = define.DEFAULT_PAGE_NO
	}
	if pageSize == 0 {
		pageSize = define.DEFAULT_PAGE_SIZE
	}

	var fileList []models.FileVo
	db := models.Db.Model(new(models.File))

	userId, _ := c.Get("userId")
	db.Where("user_id = ?", userId.(string))
	db.Where("del_flag = ?", define.RECYCLE)
	db.Order("recovery_time desc")
	//db.Offset((pageNo - 1) * pageSize).Limit(pageSize)
	//db.Find(&fileList)
	pageResult := models.QueryPageList(db, pageNo, pageSize, &fileList)
	response.ResponseOKWithData(c, pageResult)
	return
}

func RecoverFile(c *gin.Context) {
	fileIds := c.PostForm("fileIds")
	userId, _ := c.Get("userId")
	ids := strings.Split(fileIds, ",")
	finalIds := []string{}
	finalIds = append(finalIds, ids...)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	delIds := []string{}
	for _, v := range ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			childIds := FindChildrenIds(v, userId.(string))
			mu.Lock()
			delIds = append(delIds, childIds...)
			mu.Unlock()
		}()

	}
	wg.Wait()

	allIds := append(finalIds, delIds...)
	for _, v := range allIds {
		var file models.File
		models.Db.Model(new(models.File)).
			Where("id = ?", v).Find(&file)
		var count int64
		models.Db.Model(new(models.File)).
			Where("file_name = ?", file.FileName).
			Where("file_pid = ?", file.FilePid).
			Where("del_flag = ?", define.USING).Count(&count)
		if count != 0 {
			models.Db.Model(new(models.File)).Where("id = ?", v).Update("file_name", fileRename(file.FileName, userId.(string), file.FilePid))
		}
	}

	models.Db.Model(new(models.File)).
		Where("id in ?", delIds).
		Where("del_flag = ?", define.DEL).
		Where("user_id = ?", userId.(string)).
		Update("del_flag", define.USING)

	models.Db.Model(new(models.File)).
		Where("id in ?", finalIds).
		Where("user_id = ?", userId.(string)).
		Where("del_flag = ?", define.RECYCLE).
		Update("del_flag", define.USING)
	response.ResponseOK(c)
	return

}

func DelFile(c *gin.Context) {
	fileIds := c.PostForm("fileIds")
	userId, _ := c.Get("userId")
	ids := strings.Split(fileIds, ",")
	finalIds := []string{}
	finalIds = append(finalIds, ids...)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	delIds := []string{}
	for _, v := range ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			childIds := FindChildrenIds(v, userId.(string))
			mu.Lock()
			delIds = append(delIds, childIds...)
			mu.Unlock()
		}()

	}
	wg.Wait()

	finalIds = append(finalIds, delIds...)
	//todo 消息队列
	models.Db.Model(new(models.File)).
		Where("id in ?", finalIds).
		Where("user_id = ?", userId.(string)).
		Delete(&models.File{})

	models.RDb.Del(context.Background(), define.REDIS_USER_SPACE+userId.(string))
	response.ResponseOK(c)
	return
}
