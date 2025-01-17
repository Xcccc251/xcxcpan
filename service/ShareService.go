package service

import (
	"XcxcPan/common/define"
	"XcxcPan/common/helper"
	"XcxcPan/common/models"
	"XcxcPan/common/response"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func LoadShareList(c *gin.Context) {
	pageNo, _ := strconv.Atoi(c.PostForm("pageNo"))
	pageSize, _ := strconv.Atoi(c.PostForm("pageSize"))
	if pageNo == 0 {
		pageNo = define.DEFAULT_PAGE_NO
	}
	if pageSize == 0 {
		pageSize = define.DEFAULT_PAGE_SIZE
	}

	var shareList []models.ShareVo
	db := models.Db.Model(new(models.ShareVo)).
		Table("file_share fs").
		Joins("left join file_info f on f.id = fs.file_id").
		Select("fs.*,f.file_name as file_name,f.folder_type as folder_type,f.file_type as file_type,f.status as status")
	userId, _ := c.Get("userId")
	db.Where("fs.user_id = ?", userId.(string))
	db.Order("fs.share_time desc")

	pageResult := models.QueryPageList(db, pageNo, pageSize, &shareList)
	response.ResponseOKWithData(c, pageResult)
	return
}
func ShareFile(c *gin.Context) {
	fileId := c.PostForm("fileId")
	validType, _ := strconv.Atoi(c.PostForm("validType"))
	code := c.PostForm("code")
	var share models.Share
	if validType != define.FOREVER {
		validTime := define.ShareValidTypeMap[validType]
		if validTime == 0 {
			response.ResponseFail(c)
			return
		}

		expireTime := models.MyTime(time.Now().Add(time.Duration(validTime*24) * time.Hour))
		share.ExpireTime = expireTime
	}
	if code == "" {
		share.Code = helper.GetRandomStr(5)
	} else {
		share.Code = code
	}
	share.ShareTime = models.MyTime(time.Now())
	share.ShareId = helper.GetRandomStr(32)
	share.FileId = fileId
	share.UserId = c.GetString("userId")
	share.ValidType = validType
	if err := share.AddShare(models.Db); err != nil {
		response.ResponseFail(c)
		return
	}
	response.ResponseOKWithData(c, share)
	return

}

func CancelShare(c *gin.Context) {
	Ids := c.PostForm("shareIds")
	userId := c.GetString("userId")
	shareIds := strings.Split(Ids, ",")
	models.Db.Model(new(models.Share)).
		Where("user_id = ?", userId).
		Where("share_id in ?", shareIds).
		Delete(&models.Share{})
	response.ResponseOK(c)
	return

}
