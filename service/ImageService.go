package service

import (
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"github.com/gin-gonic/gin"
)

func GetImage(c *gin.Context) {
	var file models.File
	fileId := c.Param("fileId")
	models.Db.Model(new(models.File)).Where("id = ?", fileId).Find(&file)
	imagePath := define.FILE_DIR + "/" + file.UserId + "/" + fileId + "/" + define.THUMBNAIL
	c.File(imagePath)
	return
}
