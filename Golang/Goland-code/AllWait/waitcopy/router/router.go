package router

import (
	"github.com/gin-gonic/gin"
	"wait/controller"
)

func Start() {
	e := gin.Default()
	//e.Use(favicon.New("./favicon/.ico"))
	e.LoadHTMLGlob("templates/*")
	e.Static("/favicon", "./favicon")
	e.GET("/login", controller.GoLogin)
	e.POST("/login", controller.Login)

	e.POST("/register", controller.RegisterUser)
	e.GET("/register", controller.GoRegister)

	e.GET("/", controller.Index)
	e.GET("/quit", controller.Quitwait)
	// 博客操作
	e.GET("/post_index", controller.GetWaitIndex) // 博客列表
	e.POST("/post", controller.AddWait)           // 添加博客
	e.GET("/post", controller.GoAddWait)          // 跳转到添加博客页面
	e.GET("/detail", controller.WaitDetail)       // 跳转到博客详细

	e.Run()
}
