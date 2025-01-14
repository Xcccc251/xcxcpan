package service

import (
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"XcxcPan/common/response"
	"github.com/gin-gonic/gin"
)

func GetVideoInfo(c *gin.Context) {
	fileId := c.Param("fileId")
	userId, _ := c.Get("userId")
	var file models.File

	db := models.Db.Model(new(models.File)).
		Where("user_id = ?", userId.(string)).
		Where("id = ?", fileId).
		Where("file_type = ?", define.GetCategoryCodeByCategory(define.VIDEO))
	db.Find(&file)
	response.ResponseOKWithData(c, gin.H{
		"src": "http://127.0.0.1:7090/api/video/test",
	})
	return

}
