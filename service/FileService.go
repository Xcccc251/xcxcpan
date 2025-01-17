package service

import (
	"XcxcPan/Kafka"
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/common/define"
	hashRing "XcxcPan/common/hash"
	"XcxcPan/common/helper"
	"XcxcPan/common/models"
	"XcxcPan/common/redisUtil"
	"XcxcPan/common/response"
	"XcxcPan/fileServerClient_gRPC"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetFileList(c *gin.Context) {
	categoryStr := c.PostForm("category")
	fileName := c.PostForm("fileNameFuzzy")
	pageNo, _ := strconv.Atoi(c.PostForm("pageNo"))
	filePid := c.PostForm("filePid")
	if filePid == "" {
		filePid = "0"
	}

	pageSize, _ := strconv.Atoi(c.PostForm("pageSize"))
	if pageNo == 0 {
		pageNo = define.DEFAULT_PAGE_NO
	}
	if pageSize == 0 {
		pageSize = define.DEFAULT_PAGE_SIZE
	}
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
	db.Where("file_name like ?", "%"+fileName+"%")
	db.Order("last_update_time desc")
	//db.Offset((pageNo - 1) * pageSize).Limit(pageSize)
	//db.Find(&fileList)
	fmt.Println("pageNo:", pageNo, "pageSize:", pageSize, "category:", category, "filePid:", filePid)
	pageResult := models.QueryPageList(db, pageNo, pageSize, &fileList)
	response.ResponseOKWithData(c, pageResult)

	return

}

func UploadFile(c *gin.Context) {
	var uploadResultDto models.UploadResultDto
	cover := ""
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
			return
		}

	}
	//正常传递切片
	//判断磁盘空间
	tempSize := getTempFileSize(fileId, userId.(string))
	if userUseSpace.UseSpace+tempSize+int(file.Size) > userUseSpace.TotalSpace {
		DelFileChunks(fileId, userId.(string))
		response.ResponseFailWithData(c, 904, "上传失败,空间不足")
		return
	}
	chunk_id := userId.(string) + "_" + fileId + "_" + strconv.Itoa(chunkIndex)
	server_id := define.GetServerId(hashRing.Hash.Get(chunk_id))
	fmt.Println("chunk_id:", chunk_id, "server:", hashRing.Hash.Get(chunk_id))
	err := uploadChunk(chunk_id, server_id, file)
	if err != nil {
		DelFileChunks(fileId, userId.(string))

		fmt.Println(err)
		response.ResponseFailWithData(c, 0, "上传失败")
		return
	}
	go func() {
		//异步存数据库(切片)
		var chunk models.Chunk
		chunk.FileId = fileId
		chunk.ChunkId = chunk_id
		chunk.ServerId = server_id
		models.Db.Model(new(models.Chunk)).Create(&chunk)
		//redis index server键值对
		redisUtil.SetHash(define.REDIS_CHUNK+userId.(string)+":"+fileId, chunkIndex, server_id)

	}()
	saveTempFileSize(fileId, int(file.Size), userId.(string))

	if chunkIndex == chunks-1 {
		tx := models.Db.Begin()

		fileSuffix := path.Ext(fileName)
		//realFileName := helper.GetUUID() + fileSuffix
		fileName = fileRename(fileName, userId.(string), filePid)
		var newFile models.File
		newFile.Id = fileId
		newFile.FilePid = filePid
		newFile.FileMd5 = fileMd5
		newFile.UserId = userId.(string)
		newFile.FileName = fileName
		newFile.FilePath = define.FILE_DIR + "/" + userId.(string) + "/" + fileId
		newFile.ChunkPrefix = userId.(string) + "_" + fileId
		newFile.FileCategory = define.GetCategoryCodeBySuffix(fileSuffix)
		newFile.FileType = define.GetTypeCodeBySuffix(fileSuffix)
		if define.GetCategoryCodeBySuffix(fileSuffix) == define.VIDEO_CATEGORY[define.VIDEO] || define.GetCategoryCodeBySuffix(fileSuffix) == define.VIDEO_CATEGORY[define.IMAGE] {
			newFile.Status = define.FILE_TRANSFER
		} else {
			newFile.Status = define.FILE_TRANSFER_SUCCESS
		}
		newFile.FolderType = define.FILE_TYPE
		newFile.DelFlag = define.USING

		totalSize := getTempFileSize(fileId, userId.(string))
		newFile.FileSize = totalSize
		if err = tx.Model(new(models.File)).Create(&newFile).Error; err != nil {
			DelFileChunks(fileId, userId.(string))
			tx.Rollback()
			response.ResponseFailWithData(c, 0, "上传失败")
			return
		}
		if err = tx.Model(new(models.User)).Where("id = ?", userId.(string)).Update("use_space", gorm.Expr("use_space + ?", totalSize)).Error; err != nil {
			DelFileChunks(fileId, userId.(string))
			tx.Rollback()
			response.ResponseFailWithData(c, 0, "上传失败")
			return
		}

		tx.Commit()
		if define.GetCategoryCodeBySuffix(fileSuffix) == define.VIDEO_CATEGORY[define.VIDEO] || define.GetCategoryCodeBySuffix(fileSuffix) == define.VIDEO_CATEGORY[define.IMAGE] {
			go func() {
				err2 := TransferFile(fileId)
				if err2 != nil {
					//重试
				}
				if err2 == nil {
					if define.GetCategoryCodeBySuffix(fileSuffix) == define.VIDEO_CATEGORY[define.VIDEO] {
						err2 := CreateThumbnailForVideo(define.FILE_DIR + "/" + userId.(string) + "/" + fileId + fileSuffix)
						if err2 != nil {

						} else {
							cover = userId.(string) + "/" + fileId
						}

						CutFileForVideo(define.FILE_DIR + "/" + userId.(string) + "/" + fileId + fileSuffix)
					} else if define.GetCategoryCodeBySuffix(fileSuffix) == define.VIDEO_CATEGORY[define.IMAGE] {
						err2 := CreateThumbnailForImage(define.FILE_DIR + "/" + userId.(string) + "/" + fileId + fileSuffix)
						if err2 != nil {

						} else {
							cover = userId.(string) + "/" + fileId
						}
					}

					models.Db.Model(new(models.File)).
						Where("id = ?", fileId).
						Where("status = ?", define.FILE_TRANSFER).
						Clauses(clause.Locking{Strength: "UPDATE"}).
						Updates(map[string]interface{}{
							"status":     define.FILE_TRANSFER_SUCCESS,
							"file_cover": cover,
						})
				}
			}()
		}

		models.RDb.Del(context.Background(), define.REDIS_USER_SPACE+userId.(string))
		//异步转码

		uploadResultDto.FileId = fileId
		uploadResultDto.Status = define.UPLOAD_FINISH
		response.ResponseOKWithData(c, uploadResultDto)
		return

	}

	uploadResultDto.FileId = fileId
	uploadResultDto.Status = define.UPLOADING
	response.ResponseOKWithData(c, uploadResultDto)
	return

}

func FileRename(c *gin.Context) {
	fileId := c.PostForm("fileId")
	fileName := c.PostForm("fileName")
	userId, _ := c.Get("userId")
	var count int64
	db := models.Db.Model(new(models.File)).
		Where("id = ?", fileId).
		Where("user_id = ?", userId)
	db.Count(&count)
	if count == 0 {
		response.ResponseFail(c)
		return
	}
	var file models.File
	db.Find(&file)

	if file.FolderType == define.FOLDER_TYPE {
		if !CheckFolderNameIsValid(file.FilePid, fileName, userId.(string)) {
			response.ResponseFailWithData(c, 0, "存在同名文件夹")
			return
		} else {
			file.FileName = fileName
			file.LastUpdateTime = models.MyTime(time.Now())
			db.Update("last_update_time", file.LastUpdateTime)
			db.Update("file_name", fileName)
		}
	} else {
		ext := path.Ext(file.FileName)
		if !CheckFileNameIsValid(file.FilePid, fileName+ext, userId.(string)) {
			response.ResponseFailWithData(c, 0, "存在同名文件")
			return
		} else {
			file.FileName = fileName + ext
			file.LastUpdateTime = models.MyTime(time.Now())
			db.Update("last_update_time", file.LastUpdateTime)
			db.Update("file_name", fileName)
		}
	}
	response.ResponseOKWithData(c, file)
	return
}
func GetFolderList(c *gin.Context) {
	filePid := c.PostForm("filePid")
	currentFileIds := c.PostForm("currentFileIds")
	FileIds := strings.Split(currentFileIds, ",")
	userId, _ := c.Get("userId")
	var folderList []models.File
	models.Db.Model(new(models.File)).
		Where("file_pid = ?", filePid).
		Where("user_id = ?", userId.(string)).
		Where("folder_type = ?", define.FOLDER_TYPE).
		Where("id not in (?)", FileIds).
		Find(&folderList)
	response.ResponseOKWithData(c, folderList)
	return
}

func TransferFile(fileId string) error {
	var file models.File
	models.Db.Model(new(models.File)).Where("id = ?", fileId).First(&file)
	ext := path.Ext(file.FileName)
	targetPath := define.FILE_DIR + "/" + file.UserId
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return err
	}
	targetFile, err := os.Create(targetPath + "/" + fileId + ext)
	if err != nil {
		return err
	}
	tempFile, err := helper.MergeChunksToFile(helper.GetSliceMap(file.UserId, fileId))
	if err != nil {
		return err
	}
	_, err = io.Copy(targetFile, tempFile)
	if err != nil {
		return err
	}
	tempFile.Close()
	targetFile.Close()
	return nil
}

func DelFileChunks(fileId string, userId string) {
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
	models.RDb.Del(context.Background(), define.REDIS_CHUNK+userId+":"+fileId)

}

//文件预览

func GetFile(c *gin.Context) {
	userId, _ := c.Get("userId")
	fileId := c.Param("fileId")

	var dbFile models.File
	models.Db.Model(new(models.File)).Where("id = ?", fileId).Where("user_id = ?", userId).Find(&dbFile)
	splitPrefix := strings.Split(dbFile.ChunkPrefix, "_")
	data, err := DownloadFileToBytes(splitPrefix[0], splitPrefix[1])
	if err != nil {
		response.ResponseFailWithData(c, 0, "下载失败")
		return
	}
	c.Data(200, "application/octet-stream", data)
	return
}

func AddNewFolder(c *gin.Context) {
	filePid := c.PostForm("filePid")
	fileName := c.PostForm("fileName")
	userId, _ := c.Get("userId")
	if !CheckFolderNameIsValid(filePid, fileName, userId.(string)) {
		response.ResponseFailWithData(c, 0, "文件夹名重复")
		return
	} else {
		var newFile models.File
		newFile.Id = helper.GetRandomStr(32)
		newFile.FilePid = filePid
		newFile.UserId = userId.(string)
		newFile.FileName = fileName
		newFile.FolderType = define.FOLDER_TYPE
		newFile.Status = define.FILE_TRANSFER_SUCCESS
		newFile.DelFlag = define.USING
		models.Db.Model(new(models.File)).Create(&newFile)
		response.ResponseOKWithData(c, newFile)
		return
	}

}

func GetFolderInfo(c *gin.Context) {
	path := c.PostForm("path")
	userId, _ := c.Get("userId")
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

func ChangeFileFolder(c *gin.Context) {
	userId, _ := c.Get("userId")
	fileIds := c.PostForm("fileIds")
	filePid := c.PostForm("filePid")
	ids := strings.Split(fileIds, ",")
	for _, v := range ids {
		if v == filePid {
			response.ResponseFail(c)
			return
		}
		var thisFile models.File
		var count int64
		models.Db.Model(new(models.File)).
			Where("id = ?", v).
			Where("user_id = ?", userId).Find(&thisFile)
		models.Db.Model(new(models.File)).
			Where("user_id = ?", userId).
			Where("file_pid = ?", filePid).
			Where("file_name = ?", thisFile.FileName).Count(&count)

		if count > 0 {
			response.ResponseFailWithData(c, 0, "目标文件夹存在同名文件或文件夹")
			return
		}

	}
	models.Db.Model(new(models.File)).
		Where("user_id = ?", userId).
		Where("id in ?", ids).
		Update("file_pid", filePid)
	response.ResponseOK(c)
	return
}

func CreateDownloadUrl(c *gin.Context) {
	fileId := c.Param("fileId")
	userId, _ := c.Get("userId")
	var file models.File
	var count int64
	db := models.Db.Model(new(models.File)).
		Where("id = ?", fileId).
		Where("user_id = ?", userId)
	db.Count(&count)
	db.Find(&file)
	if count == 0 {
		response.ResponseFailWithData(c, 600, "文件不存在")
		return
	} else if file.FolderType == define.FOLDER_TYPE {
		response.ResponseFailWithData(c, 600, "文件夹不能下载")
		return
	}
	var dbFile models.File
	models.Db.Model(new(models.File)).Where("id = ?", fileId).Where("user_id = ?", userId).Find(&dbFile)
	splitPrefix := strings.Split(dbFile.ChunkPrefix, "_")
	code := helper.GetRandomStr(32)
	var downloadDto models.DownloadDto
	downloadDto.DownloadCode = code
	downloadDto.FileId = splitPrefix[1]
	downloadDto.UserId = splitPrefix[0]
	downloadDto.FileName = file.FileName
	downloadJson, _ := json.Marshal(&downloadDto)

	models.RDb.Set(context.Background(), define.REDIS_DOWNLOAD_CODE+":"+code, downloadJson, 5*time.Minute)

	response.ResponseOKWithData(c, code)
	return

}

func Download(c *gin.Context) {
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

func DelFileToRecycle(c *gin.Context) {
	userId, _ := c.Get("userId")
	fileIds := c.PostForm("fileIds")
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

	recoveryTime := models.MyTime(time.Now())

	go func() {
		allIds := append(finalIds, delIds...)
		delMessageJson, _ := json.Marshal(allIds)
		Kafka.ProduceMessageWithTime(define.KAFKA_DEL_TOPIC, delMessageJson, recoveryTime, define.KAFKA_DEL_DURATION)
	}()

	models.Db.Model(new(models.File)).
		Where("id in ?", delIds).
		Update("del_flag", define.DEL).
		Update("recovery_time", recoveryTime)

	models.Db.Model(new(models.File)).
		Where("id in ?", finalIds).
		Update("del_flag", define.RECYCLE).
		Update("recovery_time", recoveryTime)
	response.ResponseOK(c)
	return

}

func FindChildrenIds(fileId string, userId string) []string {
	var file models.File
	var allIds []string
	models.Db.Model(new(models.File)).
		Where("file_pid = ? and user_id = ?", fileId, userId).Find(&file)
	if file.FolderType == define.FOLDER_TYPE {
		var childIds []string
		models.Db.Model(new(models.File)).
			Select("id").
			Where("file_pid = ? and user_id = ?", fileId, userId).
			Find(&childIds)
		for _, v := range childIds {
			childIds = append(childIds, FindChildrenIds(v, userId)...)
		}
		allIds = append(allIds, childIds...)
	}
	return allIds
}

func CheckFolderNameIsValid(filePid string, fileName string, userId string) bool {
	var count int64
	models.Db.Model(new(models.File)).
		Where("file_pid = ? and file_name = ? and user_id = ?", filePid, fileName, userId).
		Where("del_flag = ?", define.USING).
		Count(&count)
	if count > 0 {
		return false
	} else {
		return true
	}
}

func CheckFileNameIsValid(filePid string, fileName string, userId string) bool {
	var count int64
	models.Db.Model(new(models.File)).
		Where("file_pid = ? and file_name = ? and user_id = ?", filePid, fileName, userId).
		Where("del_flag = ?", define.USING).
		Count(&count)
	if count > 0 {
		return false
	} else {
		return true
	}
}
func DownloadFileToBytes(userId string, fileId string) (data []byte, err error) {
	hashInt := redisUtil.GetHashInt(define.REDIS_CHUNK + userId + ":" + fileId)
	//dataMap := sync.Map{} 并发map
	dataMap := map[int][]byte{}
	wg := sync.WaitGroup{}
	var lock sync.Mutex
	for chunkIndex, serverId := range hashInt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := fileServerClient_gRPC.GetClientById(serverId)
			rsp, err := client.DownloadChunk(context.Background(), &XcXcPanFileServer.DownloadChunkRequest{
				FileName: userId + "_" + fileId + "_" + strconv.Itoa(chunkIndex),
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
	data, err = helper.MergeChunks(dataMap)
	if err != nil {
		return nil, err
	}
	return data, nil

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
		Where("file_pid = ?", filePid).
		Where("del_flag = ?", define.USING)
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
