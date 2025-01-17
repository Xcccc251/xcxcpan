package middlewares

import (
	"XcxcPan/common/define"
	"XcxcPan/common/models"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get(define.SESSION_USER_ID)
		fmt.Println("userId", userId)
		if userId == nil {
			c.JSON(200, gin.H{
				"code": 901,
				"info": "用户未登录",
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
