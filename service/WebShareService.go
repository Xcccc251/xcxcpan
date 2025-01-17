package service

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/define"
	"XcxcPan/common/helper"
	"XcxcPan/common/models"
	"XcxcPan/common/redisUtil"
	"XcxcPan/common/response"
	"XcxcPan/fileServerClient_gRPC"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetShareLoginInfo(c *gin.Context) {
	session := sessions.Default(c)
	shareId := c.PostForm("shareId")
	userId := session.Get(define.SESSION_USER_ID)
	shareDtoJson := session.Get(define.SESSION_SHARE_INFO + shareId)
	if shareDtoJson == nil {
		response.ResponseOKWithData(c, nil)
		return
	}
	var sessionShareDto models.SessionShareDto
	json.Unmarshal(shareDtoJson.([]byte), &sessionShareDto)

	var count int64
	db := models.Db.Model(new(models.Share)).Where("share_id = ?", shareId)
	db.Count(&count)
	if count == 0 {
		response.ResponseFailWithData(c, 902, "分享链接不存在或已失效")
		return
	}
	var share models.Share
	db.Find(&share)
	if share.ValidType != define.FOREVER && time.Time(share.ExpireTime).Before(time.Now()) {
		response.ResponseFailWithData(c, 902, "分享链接不存在或已失效")
		return
	}
	var shareInfoVo models.ShareInfoVo
	copier.Copy(&shareInfoVo, &share)
	var file models.File
	models.Db.Model(new(models.File)).Where("id = ?", share.FileId).Find(&file)
	if file.DelFlag != define.USING {
		response.ResponseFailWithData(c, 902, "分享链接不存在或已失效")
		return
	}
	shareInfoVo.FileName = file.FileName

	var user models.User
	models.Db.Model(new(models.User)).Where("id = ?", share.UserId).Find(&user)

	shareInfoVo.NickName = user.NickName
	shareInfoVo.Avatar = user.Avatar
	shareInfoVo.UserId = user.UserId
	if userId != nil {
		shareInfoVo.CurrentUser = user.UserId == userId.(string)
	}

	response.ResponseOKWithData(c, shareInfoVo)
	return

}

func GetShareInfo(c *gin.Context) {
	shareId := c.PostForm("shareId")
	var count int64
	db := models.Db.Model(new(models.Share)).Where("share_id = ?", shareId)
	db.Count(&count)
	if count == 0 {
		response.ResponseFailWithData(c, 902, "分享链接不存在或已失效")
		return
	}
	var share models.Share
	db.Find(&share)
	if share.ValidType != define.FOREVER && time.Time(share.ExpireTime).Before(time.Now()) {
		response.ResponseFailWithData(c, 902, "分享链接不存在或已失效")
		return
	}
	var shareInfoVo models.ShareInfoVo
	copier.Copy(&shareInfoVo, &share)
	var file models.File
	models.Db.Model(new(models.File)).Where("id = ?", share.FileId).Find(&file)
	if file.DelFlag != define.USING {
		response.ResponseFailWithData(c, 902, "分享链接不存在或已失效")
		return
	}
	shareInfoVo.FileName = file.FileName

	var user models.User
	models.Db.Model(new(models.User)).Where("id = ?", share.UserId).Find(&user)

	shareInfoVo.NickName = user.NickName
	shareInfoVo.Avatar = user.Avatar
	shareInfoVo.UserId = user.UserId
	response.ResponseOKWithData(c, shareInfoVo)
	return
}

func CheckShareCode(c *gin.Context) {
	shareId := c.PostForm("shareId")
	code := c.PostForm("code")
	var share models.Share
	var count int64
	db := models.Db.Model(new(models.Share)).Where("share_id = ?", shareId)
	db.Count(&count)
	if count == 0 {
		response.ResponseFailWithData(c, 902, "分享链接不存在或已失效")
		return
	}
	db.Find(&share)
	if share.ValidType != define.FOREVER && time.Time(share.ExpireTime).Before(time.Now()) {
		response.ResponseFailWithData(c, 902, "分享链接不存在或已失效")
		return
	}
	if share.Code != code {
		response.ResponseFailWithData(c, 902, "提取码错误")
		return
	}
	db.Update("show_count", gorm.Expr("show_count + ?", 1))
	var sessionShareDto models.SessionShareDto
	sessionShareDto.ShareId = shareId
	sessionShareDto.ExpireTime = share.ExpireTime
	sessionShareDto.FileId = share.FileId
	sessionShareDto.ShareUserId = share.UserId
	sessionShareDtoJson, _ := json.Marshal(sessionShareDto)
	session := sessions.Default(c)
	session.Set(define.SESSION_SHARE_INFO+shareId, sessionShareDtoJson)
	if err := session.Save(); err != nil {
		response.ResponseFail(c)
		return
	}
	response.ResponseOKWithData(c, sessionShareDto)
	return
}

func LoadShareFileList(c *gin.Context) {
	shareId := c.PostForm("shareId")
	pageNo, _ := strconv.Atoi(c.PostForm("pageNo"))
	filePid := c.PostForm("filePid")
	pageSize, _ := strconv.Atoi(c.PostForm("pageSize"))
	if pageNo == 0 {
		pageNo = define.DEFAULT_PAGE_NO
	}
	if pageSize == 0 {
		pageSize = define.DEFAULT_PAGE_SIZE
	}
	session := sessions.Default(c)
	shareDtoJson := session.Get(define.SESSION_SHARE_INFO + shareId)
	if shareDtoJson == nil {
		response.ResponseOKWithData(c, nil)
		return
	}
	var sessionShareDto models.SessionShareDto
	json.Unmarshal(shareDtoJson.([]byte), &sessionShareDto)
	var share models.Share
	models.Db.Model(new(models.Share)).Where("share_id = ?", shareId).Find(&share)

	var fileList []models.FileVo
	db := models.Db.Model(new(models.File))
	if filePid == "0" {
		db.Where("id = ?", share.FileId)
	} else {
		db.Where("file_pid = ?", filePid)
	}

	db.Where("del_flag = ?", define.USING)
	db.Order("last_update_time desc")
	//db.Offset((pageNo - 1) * pageSize).Limit(pageSize)
	//db.Find(&fileList)

	pageResult := models.QueryPageList(db, pageNo, pageSize, &fileList)
	response.ResponseOKWithData(c, pageResult)

	return

}

func GetShareVideoInfo(c *gin.Context) {
	target := c.Param("target")
	shareId := c.Param("shareId")
	session := sessions.Default(c)
	shareDtoJson := session.Get(define.SESSION_SHARE_INFO + shareId)
	if shareDtoJson == nil {
		fmt.Println("nil")
		response.ResponseOKWithData(c, nil)
		return
	}
	var sessionShareDto models.SessionShareDto
	json.Unmarshal(shareDtoJson.([]byte), &sessionShareDto)
	userId := sessionShareDto.ShareUserId

	ext := path.Ext(target)
	if ext != ".ts" {
		var share models.Share
		models.Db.Model(new(models.Share)).Where("share_id = ?", shareId).Find(&share)
		if share.FileId != target {
			response.ResponseFail(c)
			return
		}
		m3u8Path := define.FILE_DIR + "/" + userId + "/" + target + "/" + define.M3U8
		c.File(m3u8Path)
		return
	} else {
		splitTarget := strings.Split(target, "_")
		targetUserId := splitTarget[0]
		targetFileId := splitTarget[1]
		if targetUserId != userId {
			response.ResponseFail(c)
			return
		}

		tsPath := define.FILE_DIR + "/" + targetUserId + "/" + targetFileId + "/" + target
		c.File(tsPath)
		return
	}

	return
}

func GetShareFolderInfo(c *gin.Context) {
	shareId := c.PostForm("shareId")
	path := c.PostForm("path")
	sessionShareDto := getSessionShareDto(shareId, c)
	userId := sessionShareDto.ShareUserId
	pathArray := strings.Split(path, "/")
	var folderList []models.File
	models.Db.Model(new(models.File)).
		Where("user_id = ?", userId).
		Where("folder_type = ?", define.FOLDER_TYPE).
		Where("id in ?", pathArray).
		Order("field(id,'" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(pathArray)), "','"), "[]") + "')").
		Find(&folderList)
	response.ResponseOKWithData(c, folderList)
	return

}

func getSessionShareDto(shareId string, c *gin.Context) models.SessionShareDto {
	session := sessions.Default(c)
	shareDtoJson := session.Get(define.SESSION_SHARE_INFO + shareId)
	if shareDtoJson == nil {
		response.ResponseOKWithData(c, nil)
	}
	var sessionShareDto models.SessionShareDto
	json.Unmarshal(shareDtoJson.([]byte), &sessionShareDto)
	return sessionShareDto
}

func SaveShare(c *gin.Context) {
	shareId := c.PostForm("shareId")
	shareFileIds := c.PostForm("shareFileIds")
	myFolderId := c.PostForm("myFolderId")
	fileIds := strings.Split(shareFileIds, ",")
	shareDto := getSessionShareDto(shareId, c)
	userId := c.GetString("userId")
	if userId == shareDto.ShareUserId {
		response.ResponseFailWithData(c, 0, "不能分享给自己")
		return
	}

	var fileList []models.File
	models.Db.Model(new(models.File)).
		Where("id in ?", fileIds).
		Where("user_id = ?", shareDto.ShareUserId).Find(&fileList)

	for _, file := range fileList {
		if file.FolderType != define.FOLDER_TYPE {
			file.UserId = userId
			file.FilePid = myFolderId
			file.Id = helper.GetRandomStr(32)
			models.Db.Model(new(models.File)).Create(&file)
		} else {
			newId := helper.GetRandomStr(32)
			childrenFiles := FindChildrenFilesWithNewFilePid(file.Id, shareDto.ShareUserId, newId, userId)
			file.UserId = userId
			file.FilePid = myFolderId
			file.Id = newId
			models.Db.Model(new(models.File)).Create(&file)
			models.Db.Model(new(models.File)).Create(&childrenFiles)
		}

	}
	response.ResponseOK(c)
	return

}

func FindChildrenFilesWithNewFilePid(fileId string, userId string, newFilePid string, newUserId string) []models.File {
	var allChildrenFiles []models.File
	var childIds []string
	models.Db.Model(new(models.File)).
		Select("id").
		Where("file_pid = ? and user_id = ?", fileId, userId).
		Find(&childIds)
	for _, v := range childIds {
		var childFile models.File
		models.Db.Model(new(models.File)).
			Where("id = ?", v).
			Find(&childFile)
		childFile.UserId = newUserId
		childFile.Id = helper.GetRandomStr(32)
		childFile.FilePid = newFilePid
		allChildrenFiles = append(allChildrenFiles, childFile)
		allChildrenFiles = append(allChildrenFiles, FindChildrenFilesWithNewFilePid(v, userId, childFile.Id, newUserId)...)
	}
	return allChildrenFiles
}

func GetShareFile(c *gin.Context) {
	fileId := c.Param("fileId")
	shareId := c.Param("shareId")
	var count int64
	models.Db.Model(new(models.Share)).
		Where("share_id = ?", shareId).
		Where("file_id = ?", fileId).Count(&count)
	if count == 0 {
		response.ResponseFail(c)
		return
	}
	shareDto := getSessionShareDto(shareId, c)
	if shareDto.FileId != fileId {
		response.ResponseFail(c)
		return
	}

	var dbFile models.File
	models.Db.Model(new(models.File)).
		Where("id = ?", fileId).
		Where("user_id = ?", shareDto.ShareUserId).Find(&dbFile)
	splitPrefix := strings.Split(dbFile.ChunkPrefix, "_")
	data, err := DownloadFileToBytes(splitPrefix[0], splitPrefix[1])
	if err != nil {
		response.ResponseFailWithData(c, 0, "下载失败")
		return
	}
	c.Data(200, "application/octet-stream", data)
	return
}

func CreateShareFileDownloadUrl(c *gin.Context) {
	fileId := c.Param("fileId")
	shareId := c.Param("shareId")
	shareDto := getSessionShareDto(shareId, c)
	var count int64
	db := models.Db.Model(new(models.Share)).
		Where("share_id = ?", shareId).
		Where("file_id = ?", fileId)
	db.Count(&count)
	var dbFile models.File
	models.Db.Model(new(models.File)).Where("id = ?", fileId).Where("user_id = ?", shareDto.ShareUserId).Find(&dbFile)
	if count == 0 {
		response.ResponseFailWithData(c, 600, "文件不存在")
		return
	} else if dbFile.FolderType == define.FOLDER_TYPE {
		response.ResponseFailWithData(c, 600, "文件夹不能下载")
		return
	}

	splitPrefix := strings.Split(dbFile.ChunkPrefix, "_")
	code := helper.GetRandomStr(32)
	var downloadDto models.DownloadDto
	downloadDto.DownloadCode = code
	downloadDto.FileId = splitPrefix[1]
	downloadDto.UserId = splitPrefix[0]
	downloadDto.FileName = dbFile.FileName
	downloadJson, _ := json.Marshal(&downloadDto)

	models.RDb.Set(context.Background(), define.REDIS_DOWNLOAD_CODE+":"+code, downloadJson, 5*time.Minute)

	response.ResponseOKWithData(c, code)
	return

}

func Download4ShareFile(c *gin.Context) {
	code := c.Param("code")
	var downloadDto models.DownloadDto
	result, _ := models.RDb.Get(context.Background(), define.REDIS_DOWNLOAD_CODE+":"+code).Result()
	json.Unmarshal([]byte(result), &downloadDto)

	hashInt := redisUtil.GetHashInt(define.REDIS_CHUNK + downloadDto.UserId + ":" + downloadDto.FileId)
	var dataMap = map[int][]byte{}
	wg := sync.WaitGroup{}
	var lock sync.Mutex
	for chunkIndex, serverId := range hashInt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := fileServerClient_gRPC.GetClientById(serverId)
			rsp, err := client.DownloadChunk(context.Background(), &XcXcPanFileServer.DownloadChunkRequest{
				FileName: downloadDto.UserId + "_" + downloadDto.FileId + "_" + strconv.Itoa(chunkIndex),
				Server:   int64(serverId),
			})
			//加锁保护map
			lock.Lock()
			dataMap[chunkIndex] = rsp.Data
			lock.Unlock()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	wg.Wait()
	data, err := helper.MergeChunks(dataMap)
	if err != nil {
		response.ResponseFailWithData(c, 0, "下载失败")
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+downloadDto.FileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(200, "application/octet-stream", data)
	return

}
