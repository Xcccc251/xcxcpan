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
	// 初始化 Cookie 的存储引擎
	store := cookie.NewStore([]byte("xcxcpan_secret")) //session加密密钥
	r.Use(sessions.Sessions("xcxcpan_session", store))
	r.Use(middlewares.CorsMiddleWare()) //跨域
	v1 := r.Group("/api")

	v1.GET("/checkCode", service.CheckCode)
	v1.POST("/sendEmailCode", service.SendEmailCode)
	v1.POST("/register", service.Register)
	v1.POST("/login", service.Login)
	v1.POST("/resetPwd", service.ResetPassword)
	v1.GET("/getAvatar/:userId", service.GetAvatar)
	v1.GET("/getUserInfo", middlewares.AuthUserCheck(), service.GetUserInfo)
	v1.POST("/getUseSpace", middlewares.AuthUserCheck(), service.GetUseSpace)
	v1.POST("/logout", middlewares.AuthUserCheck(), service.Logout)
	v1.POST("/updateUserAvatar", middlewares.AuthUserCheck(), service.UpdateUserAvatar)
	v1.POST("/qqlogin", service.QQLogin)
	file := v1.Group("/file")
	{
		file.POST("/loadDataList", middlewares.AuthUserCheck(), service.GetFileList)
		file.POST("/uploadFile", middlewares.AuthUserCheck(), service.UploadFile)
	}

	return r
}
