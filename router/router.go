package router

import (
	"XcxcPan/common/middlewares"
	"XcxcPan/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	// 初始化基于 Cookie 的存储引擎
	store := cookie.NewStore([]byte("xcxcpan_secret")) //session加密密钥
	r.Use(sessions.Sessions("xcxcpan_session", store))
	r.Use(middlewares.CorsMiddleWare())
	v1 := r.Group("/api")

	v1.GET("/checkCode", service.CheckCode)
	v1.POST("/sendEmailCode", service.SendEmailCode)
	v1.POST("/register", service.Register)
	v1.POST("/login", service.Login)
	v1.POST("/resetPwd", service.ResetPassword)
	v1.GET("/getAvatar/:userId", service.GetAvatar)

	return r
}
