package router

import (
	"github.com/gin-gonic/gin"
)

func Start() {
	//创建服务
	e := gin.Default()

	//获取首页界面
	//e.GET("/", controller.Index)

	e.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	//服务器端口，默认8080
	e.Run(":8080")
}
