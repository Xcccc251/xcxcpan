package service

import (
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

func GetVideoInfo(c *gin.Context) {
	target := c.Param("target")
	userId := c.GetString("userId")

	ext := path.Ext(target)
	if ext != ".ts" {
		var dbFile models.File
		models.Db.Model(new(models.File)).Where("id = ?", target).Where("user_id = ?", userId).Find(&dbFile)
		splitPrefix := strings.Split(dbFile.ChunkPrefix, "_")
		m3u8Path := define.FILE_DIR + "/" + splitPrefix[0] + "/" + splitPrefix[1] + "/" + define.M3U8
		c.File(m3u8Path)
		return
	} else {
		splitTarget := strings.Split(target, "_")
		targetUserId := splitTarget[0]
		targetFileId := splitTarget[1]

		tsPath := define.FILE_DIR + "/" + targetUserId + "/" + targetFileId + "/" + target
		c.File(tsPath)
		return
	}

	return

}
