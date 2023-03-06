package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"reentry-distributed-lock/example/saleapi"
)

func main() {
	server := gin.Default()

	server.GET("/buy", saleapi.Sale)

	err := server.Run(":8080")
	if err != nil {
		fmt.Println("服务启动失败：" + err.Error())
		return
	}
}
