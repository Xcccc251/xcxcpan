package middlewares

import (
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"context"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get(define.SESSION_USER_ID)
		if userId == nil {
			c.JSON(901, gin.H{
				"code": 901,
				"info": "用户未登录",
			})
			c.Abort()
			return
		}
		var userLoginDto models.UserLoginDto
		result, _ := models.RDb.Get(context.Background(), define.REDIS_USER_INFO+userId.(string)).Result()
		json.Unmarshal([]byte(result), &userLoginDto)
		if !userLoginDto.IsAdmin {
			c.JSON(404, gin.H{
				"code": 404,
				"info": "无权限",
			})
			c.Abort()
			return
		}

		models.RDb.Expire(context.Background(), define.REDIS_USER_INFO+userId.(string), define.EXPIRE_DAY)
		models.RDb.Expire(context.Background(), define.REDIS_USER_SPACE+userId.(string), define.EXPIRE_DAY)

		c.Set("userId", userId.(string))
		c.Next()
	}
}
