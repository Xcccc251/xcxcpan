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
	v1.POST("/updatePassword", middlewares.AuthUserCheck(), service.UpdatePassword)
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
		file.GET("/getImage/:userId/:fileId", service.GetImage)
		file.POST("/getFile/:fileId", middlewares.AuthUserCheck(), service.GetFile)
		file.GET("/getFile/:fileId", middlewares.AuthUserCheck(), service.GetFile)
		file.POST("/newFolder", middlewares.AuthUserCheck(), service.AddNewFolder)
		file.POST("/getFolderInfo", middlewares.AuthUserCheck(), service.GetFolderInfo)
		file.POST("/rename", middlewares.AuthUserCheck(), service.FileRename)
		file.POST("/loadAllFolder", middlewares.AuthUserCheck(), service.GetFolderList)
		file.POST("/changeFileFolder", middlewares.AuthUserCheck(), service.ChangeFileFolder)
		file.POST("createDownloadUrl/:fileId", middlewares.AuthUserCheck(), service.CreateDownloadUrl)
		file.GET("/download/:code", service.Download)
		file.POST("/delFile", middlewares.AuthUserCheck(), service.DelFileToRecycle)
		ts := file.Group("/ts")
		{
			ts.GET("/getVideoInfo/:target", middlewares.AuthUserCheck(), service.GetVideoInfo)
		}
	}
	recycle := v1.Group("/recycle")
	{
		recycle.POST("/loadRecycleList", middlewares.AuthUserCheck(), service.LoadRecycleList)
		recycle.POST("/recoverFile", middlewares.AuthUserCheck(), service.RecoverFile)
		recycle.POST("/delFile", middlewares.AuthUserCheck(), service.DelFile)
	}
	share := v1.Group("/share")
	{
		share.POST("/loadShareList", middlewares.AuthUserCheck(), service.LoadShareList)
		share.POST("/shareFile", middlewares.AuthUserCheck(), service.ShareFile)
		share.POST("/cancelShare", middlewares.AuthUserCheck(), service.CancelShare)

	}
	showShare := v1.Group("/showShare")
	{
		showShare.POST("/getShareInfo", service.GetShareInfo)
		showShare.POST("/getShareLoginInfo", service.GetShareLoginInfo)
		showShare.POST("/checkShareCode", service.CheckShareCode)
		showShare.POST("/loadFileList", service.LoadShareFileList)
		showShare.GET("/ts/getVideoInfo/:shareId/:target", service.GetShareVideoInfo)
		showShare.POST("/getFolderInfo", service.GetShareFolderInfo)
		showShare.POST("/saveShare", middlewares.AuthUserCheck(), service.SaveShare)
		showShare.GET("/getFile/:shareId/:fileId", service.GetShareFile)
		showShare.POST("/createDownloadUrl/:shareId/:fileId", service.CreateShareFileDownloadUrl)
		showShare.GET("/download/:code", service.Download4ShareFile)
	}

	return r
}
