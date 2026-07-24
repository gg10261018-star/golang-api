package main

import (
	"github.com/example/golang-api/app/points"
	"github.com/example/golang-api/app/user"
	"github.com/example/golang-api/app/weather"
	"github.com/example/golang-api/internal/middleware"
	"github.com/example/golang-api/pkg/connect"
	"github.com/example/golang-api/pkg/env"
	"github.com/gin-gonic/gin"
)

func main() {
	env.Init()

	// 預先建立 DB 連線
	if _, err := connect.NewDB("APP"); err != nil {
		panic(err)
	}

	// 天氣排程：啟動時立即執行一次，之後每 24 小時更新
	//weather.StartScheduler()

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.AssignDBSession)
	engine.Use(middleware.CloseDBSession)

	router := engine.Group("/api/v1")

	// 不需要 JWT
	user.Router(router)
	weather.Router(router)

	// 需要 JWT
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth(env.Getenv("JWT_SECRET")))
	{
		user.RouterProtected(protected)
		weather.RouterProtected(protected)
		points.RouterProtected(protected)
	}

	if err := engine.Run(":" + env.Getenv("PORT", "8080")); err != nil {
		panic(err)
	}
}
