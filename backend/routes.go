package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	router.Use(cors.New(config))
	v1 := router.Group("/api/v1")
	{
		v1.POST("/signup", signup)
		v1.POST("/login", login)
		v1.POST("/upload", authMiddleware, upload)
		v1.GET("/undo", undo)
		v1.GET("/redo", redo)
		v1.POST("/crop", authMiddleware, crop)
		v1.POST("/resize", authMiddleware)
		v1.POST("/rotate", authMiddleware, rotate)
		v1.POST("/color", authMiddleware)
		v1.POST("/filter", authMiddleware)
		v1.GET("/export", authMiddleware, export)
		v1.GET("/save", authMiddleware, save)
		v1.GET("/getImage", authMiddleware, getImage)
	}
	return router
}
