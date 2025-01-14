package service

import (
	"XcxcPan/common/define"
	"XcxcPan/common/response"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

func GetVideoInfo(c *gin.Context) {
	target := c.Param("target")
	userId, _ := c.Get("userId")
	ext := path.Ext(target)
	if ext != ".ts" {
		m3u8Path := define.FILE_DIR + "/" + userId.(string) + "/" + target + "/" + define.M3U8
		c.File(m3u8Path)
		return
	} else {
		splitTarget := strings.Split(target, "_")
		targetUserId := splitTarget[0]
		targetFileId := splitTarget[1]
		if targetUserId != userId.(string) {
			response.ResponseFail(c)
			return
		}

		tsPath := define.FILE_DIR + "/" + targetUserId + "/" + targetFileId + "/" + target
		c.File(tsPath)
		return
	}

	return

}
