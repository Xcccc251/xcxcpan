package service

import (
	"XcxcPan/common/define"
	"github.com/gin-gonic/gin"
)

func GetImage(c *gin.Context) {
	userId := c.Param("userId")
	fileId := c.Param("fileId")
	imagePath := define.FILE_DIR + "/" + userId + "/" + fileId + "/" + define.THUMBNAIL
	c.File(imagePath)
	return
}
