package router

import (
	"github.com/gin-gonic/gin"
	"waittest01/controller"
)

//单打开waittest01不行，要打开总文件

func Start() {

	// 创建一个默认的路由引擎
	e := gin.Default()
	e.LoadHTMLGlob("templates/*")

	e.GET("/index", controller.Index)

	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	e.Run(":8066")
}
