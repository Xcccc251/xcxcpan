package middlewares

//import (
//	"XcxcVideo/common/define"
//	"XcxcVideo/common/helper"
//	"XcxcVideo/common/models"
//	"context"
//	"github.com/gin-gonic/gin"
//	"strconv"
//	"strings"
//)
//
//func AuthToken() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		auth := c.GetHeader("authorization")
//		if auth == "" {
//			c.Next()
//			return
//		}
//		token := strings.TrimPrefix(auth, "Bearer ")
//
//		userClaim, _ := helper.AnalysisToken(token)
//		userId := userClaim.UserId
//		userIdStr := strconv.Itoa(userId)
//		models.RDb.Expire(context.Background(), define.TOKEN_PREFIX+userIdStr, define.TOKEN_TTL)
//		models.RDb.Expire(context.Background(), define.USER_PREFIX+userIdStr, define.TOKEN_TTL)
//		c.Next()
//
//	}
//}
