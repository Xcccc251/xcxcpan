package service

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/define"
	"XcxcPan/common/fileServerClient_gRPC"
	hashRing "XcxcPan/common/hash"
	"XcxcPan/common/helper"
	"XcxcPan/common/models"
	"XcxcPan/common/redisUtil"
	"XcxcPan/common/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path"
	"strconv"
	"time"
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
	var uploadResultDto models.UploadResultDto
	file, _ := c.FormFile("file")
	fileName := c.PostForm("fileName")
	//ext := path.Ext(fileName)
	fileMd5 := c.PostForm("fileMd5")
	chunkIndex, _ := strconv.Atoi(c.PostForm("chunkIndex"))
	chunks, _ := strconv.Atoi(c.PostForm("chunks"))
	fileId := c.PostForm("fileId")
	filePid := c.PostForm("filePid")
	userId, _ := c.Get("userId")
	var userLoginDto models.UserLoginDto
	result, _ := models.RDb.Get(context.Background(), define.REDIS_USER_INFO+userId.(string)).Result()
	json.Unmarshal([]byte(result), &userLoginDto)

	if fileId == "" {
		fileId = helper.GetRandomStr(32)
	}
	var userUseSpace = getUserUseSpace(userId.(string))
	if chunkIndex == 0 {
		db := models.Db.Model(new(models.File)).Where("file_md5=?", fileMd5).Where("status=?", define.FILE_TRANSFER_SUCCESS)
		var count int64
		var dbFile models.File
		db.Count(&count)
		//秒传逻辑
		if count > 0 {
			db.First(&dbFile)
			if dbFile.FileSize+userUseSpace.UseSpace > userUseSpace.TotalSpace {
				response.ResponseFailWithData(c, 904, "上传失败,空间不足")
			}
			dbFile.Id = fileId
			dbFile.FilePid = filePid
			dbFile.UserId = userId.(string)
			dbFile.LastUpdateTime = models.MyTime(time.Now())
			dbFile.Status = define.FILE_TRANSFER_SUCCESS
			dbFile.DelFlag = define.USING
			dbFile.FileMd5 = fileMd5
			dbFile.FileName = fileRename(fileName, userId.(string), filePid)
			models.Db.Model(new(models.File)).Create(&dbFile)
			uploadResultDto.FileId = fileId
			uploadResultDto.Status = define.UPLOAD_SECONDS
			updateUserUseSpace(userUseSpace, dbFile.FileSize, userId.(string))
			response.ResponseOKWithData(c, uploadResultDto)
		}

	}
	//正常传递切片
	//判断磁盘空间
	tempSize := getTempFileSize(fileId, userId.(string))
	if userUseSpace.UseSpace+tempSize+int(file.Size) > userUseSpace.TotalSpace {
		delFileChunks(fileId, userId.(string))
		response.ResponseFailWithData(c, 904, "上传失败,空间不足")
		return
	}
	chunk_id := userId.(string) + "_" + fileId + "_" + strconv.Itoa(chunkIndex)
	server_id := define.GetServerId(hashRing.Hash.Get(chunk_id))
	err := uploadChunk(chunk_id, server_id, file)
	if err != nil {
		delFileChunks(fileId, userId.(string))
		fmt.Println(err)
		response.ResponseFailWithData(c, 0, "上传失败")
		return
	}
	go func() {
		//todo kafka
		//异步存数据库
		var chunk models.Chunk
		chunk.FileId = fileId
		chunk.ChunkId = chunk_id
		chunk.ServerId = server_id
		models.Db.Model(new(models.Chunk)).Create(&chunk)
		//redis index server键值对
		redisUtil.SetHash(define.REDIS_CHUNK+userId.(string)+":"+fileId, chunkIndex, server_id)

	}()
	saveTempFileSize(fileId, int(file.Size), userId.(string))
	uploadResultDto.FileId = fileId
	uploadResultDto.Status = define.UPLOADING
	if chunkIndex == chunks-1 {
		uploadResultDto.Status = define.UPLOAD_FINISH
		response.ResponseOKWithData(c, uploadResultDto)
		return

	}

	response.ResponseOKWithData(c, uploadResultDto)
	return

}

func delFileChunks(fileId string, userId string) {
	chunkIdServerIdMap := redisUtil.GetHashInt(define.REDIS_CHUNK + userId + ":" + fileId)
	for chunkId, serverId := range chunkIdServerIdMap {
		client := fileServerClient_gRPC.GetClientById(serverId)
		_, err := client.DelChunk(context.Background(), &XcXcPanFileServer.DelChunkRequest{
			FileName: userId + "_" + fileId + "_" + strconv.Itoa(chunkId),
			Server:   int64(serverId),
		})
		if err != nil {
			fmt.Println(err)
		}
	}

}
func uploadChunk(chunk_id string, server_id int, file *multipart.FileHeader) error {
	fileOpen, err := file.Open()
	defer fileOpen.Close()
	if err != nil {
		return err
	}
	osFile, err := helper.SaveMultipartFile(fileOpen)
	defer osFile.Close()
	if err != nil {
		return err
	}
	data, err := helper.FileToBytes(osFile)
	if err != nil {
		return err
	}
	client := fileServerClient_gRPC.GetClientById(server_id)
	_, err = client.UploadChunk(context.Background(), &XcXcPanFileServer.UploadChunkRequest{
		Data:     data,
		FileName: chunk_id,
		Server:   int64(define.GetServerId(hashRing.Hash.Get(chunk_id))),
	})
	if err != nil {
		return err
	}
	return nil

}

func saveTempFileSize(fileId string, fileSize int, userId string) {
	key := define.REDIS_TEMP_FILE + userId + ":" + fileId
	var size int
	result, _ := models.RDb.Get(context.Background(), key).Result()
	if result != "" {
		size, _ = strconv.Atoi(result)
	}
	models.RDb.Set(context.Background(), key, size+fileSize, define.EXPIRE_HOUR)
}

func getTempFileSize(fileId string, userId string) int {
	key := define.REDIS_TEMP_FILE + userId + ":" + fileId
	var size int
	result, _ := models.RDb.Get(context.Background(), key).Result()
	if result != "" {
		size, _ = strconv.Atoi(result)
	}
	return size
}
func fileRename(fileName string, userId string, filePid string) string {
	db := models.Db.Model(new(models.File)).
		Where("file_name = ?", fileName).
		Where("user_id = ?", userId).
		Where("file_pid = ?", filePid)
	var count int64
	db.Count(&count)
	if count > 0 {
		ext := path.Ext(fileName)
		fileName = fileName[:len(fileName)-len(ext)] + "_" + helper.GetRandomStr(5) + ext
	}
	return fileName
}

func updateUserUseSpace(userSpaceDto models.UserSpaceDto, totalSize int, userId string) {
	userSpaceDto.UseSpace += totalSize

	models.Db.Model(new(models.User)).Where("id = ?", userId).Update("use_space", userSpaceDto.UseSpace)

	userSpaceJson, _ := json.Marshal(userSpaceDto)
	models.RDb.Set(context.Background(), define.REDIS_USER_SPACE+userId, userSpaceJson, define.EXPIRE_DAY)
}
func getUseSpaceById(userId string) int {
	var useSpace int
	models.Db.Model(new(models.File)).
		Select("ifnull(sum(file_size),0) as total_size").
		Where("user_id = ?", userId).
		Find(&useSpace)
	return useSpace

}
