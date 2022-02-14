package http_server

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type HandlerConfig struct {
	Port string
}

type Incrementer interface {
	IncValueHandle(ctx *gin.Context)
}

type HashGiver interface {
	GetHMACHandle(ctx *gin.Context)
}

type Multiplier interface {
	MultiplyHandle(ctx *gin.Context)
}

type HandlerSet struct {
	Incrementer
	HashGiver
	Multiplier
}

func NewHandlerSet(incH Incrementer, hmacH HashGiver, mulH Multiplier) HandlerSet {
	return HandlerSet{
		Incrementer: incH,
		HashGiver:   hmacH,
		Multiplier:  mulH,
	}
}

func NewRouter(h HandlerSet) (router *gin.Engine) {
	router = gin.Default()

	router.POST("/test1", h.Incrementer.IncValueHandle)
	router.POST("/test2", h.HashGiver.GetHMACHandle)
	router.POST("/test3", h.Multiplier.MultiplyHandle)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//router.Run(":" + conf.Port)

	return router
}
