package middlewares

//
//import (
//	"XcxcVideo/common/define"
//	"XcxcVideo/common/helper"
//	"XcxcVideo/common/models"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"strconv"
//	"strings"
//)
//
//func AuthUserCheck() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		auth := c.GetHeader("authorization")
//		token := strings.TrimPrefix(auth, "Bearer ")
//		if token == "" {
//			fmt.Println("token is null")
//			c.Abort()
//			c.JSON(401, gin.H{"code": 401, "msg": "请先登录"})
//			return
//		}
//		userClaim, _ := helper.AnalysisToken(token)
//		userId := userClaim.UserId
//		userIdStr := strconv.Itoa(userId)
//
//		userResult, err := models.RDb.Get(context.Background(), define.USER_PREFIX+userIdStr).Result()
//		if err != nil {
//			fmt.Println("redis get user error:", err)
//			c.Abort()
//			c.JSON(401, gin.H{"code": 401, "msg": "请先登录"})
//			return
//		}
//
//		// 刷新过期时间
//		models.RDb.Expire(context.Background(), define.TOKEN_PREFIX+userIdStr, define.TOKEN_TTL)
//		models.RDb.Expire(context.Background(), define.USER_PREFIX+userIdStr, define.TOKEN_TTL)
//		var user models.UserVo
//		json.Unmarshal([]byte(userResult), &user)
//
//		c.Set("userId", userId)
//		c.Set("user", user)
//		c.Next()
//	}
//}
