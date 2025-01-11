package middlewares

//
//import (
//	"XcxcPan/common/helper"
//	"XcxcPan/common/models"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"strconv"
//	"strings"
//)
//
//func AuthAdminCheck() gin.HandlerFunc {
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
//		userResult, err := models.RDb.Get(context.Background(), define.ADMIN_PREFIX+userIdStr).Result()
//		if err != nil {
//			fmt.Println("redis get user error:", err)
//			c.Abort()
//			c.JSON(401, gin.H{"code": 401, "msg": "请先登录"})
//			return
//		}
//		var user models.UserVo
//		json.Unmarshal([]byte(userResult), &user)
//		if user.Role == 0 {
//			c.Abort()
//			c.JSON(401, gin.H{"code": 403, "msg": "无权限"})
//			return
//		}
//		if user.State == 2 {
//			c.Abort()
//			c.JSON(401, gin.H{"code": 404, "msg": "账号已注销"})
//			return
//		}
//		if user.State == 1 {
//			c.Abort()
//			c.JSON(401, gin.H{"code": 403, "msg": "账号封禁中"})
//			return
//		}
//		// 刷新过期时间
//		models.RDb.Expire(context.Background(), define.TOKEN_PREFIX+userIdStr, define.TOKEN_TTL)
//		models.RDb.Expire(context.Background(), define.ADMIN_PREFIX+userIdStr, define.TOKEN_TTL)
//
//		c.Set("userId", userId)
//		c.Set("user", user)
//		c.Next()
//	}
//}
