package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/luxifa/beauty/config"
	"github.com/luxifa/beauty/handlers"
)

func main() {
    // 初始化数据库
    config.InitDB()

    // 创建Gin实例
    r := gin.Default()

    // 配置CORS中间件
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},                            // 允许所有来源，生产环境建议设置为具体域名
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // 允许的HTTP方法
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,                                     // 允许发送cookie
        MaxAge:           12 * time.Hour,                           // 预检请求的有效期
    }))

    // 创建处理器
    authHandler := handlers.NewAuthHandler()

    // 注册路由
    r.POST("/api/register", authHandler.Register)
    r.POST("/api/login", authHandler.Login)
    r.GET("/api/users/:id", authHandler.GetUserInfo)

    // 启动服务器
    log.Println("服务器启动在 http://localhost:8080")
    r.Run(":8080")
}