package middlewares

//import (
//	"SecKill/common/response"
//	"github.com/gin-gonic/gin"
//)
//
//func CORSMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 设置允许跨域的源
//		c.Header("Access-Control-Allow-Origin", "*")
//		// 允许的请求方法
//		c.Header("Access-Control-Allow-Methods", "*")
//		// 允许的请求头
//		c.Header("Access-Control-Allow-Headers", "*")
//		// 浏览器可以访问的响应头
//		c.Header("Access-Control-Expose-Headers", "*")
//		// 是否允许携带身份验证信息
//		c.Header("Access-Control-Allow-Credentials", "true")
//		// 缓存预检请求的时间
//		c.Header("Access-Control-Max-Age", "43200") // 12小时
//
//		// 如果是预检请求（OPTIONS），直接返回 200
//		if c.Request.Method == "OPTIONS" {
//			response.ResponseOK(c, nil)
//			return
//		}
//
//		// 继续处理请求
//		c.Next()
//	}
//}
