package Router

import (
	"main/Auth"
	"main/Operations"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowHeaders = append(config.AllowHeaders, "filename")
	config.ExposeHeaders = append(config.ExposeHeaders, "filename")

	router.Use(cors.New(config))
	v1 := router.Group("/api/v1")
	{
		v1.POST("/signup", Auth.Signup)
		v1.POST("/login", Auth.Login)
		v1.POST("/upload", Auth.AuthMiddleware, Operations.Upload)
		// v1.GET("/undo", undo)
		// v1.GET("/redo", redo)
		// v1.POST("/crop", Auth.AuthMiddleware, crop)
		v1.POST("/resize", Auth.AuthMiddleware)
		v1.POST("/rotate", Auth.AuthMiddleware, Operations.Rotate)
		v1.GET("/grayscale", Auth.AuthMiddleware, Operations.Grayscale)
		v1.GET("/blurinc", Auth.AuthMiddleware)
		v1.POST("/blurdec", Auth.AuthMiddleware)
		// v1.GET("/sharpinc", Auth.AuthMiddleware, sharpinc)
		v1.POST("/sharpdec", Auth.AuthMiddleware)
		v1.GET("/brightinc", Auth.AuthMiddleware, Operations.Bright_inc)
		v1.GET("/brightdec", Auth.AuthMiddleware, Operations.Bright_dec)
		v1.GET("/contrastinc", Auth.AuthMiddleware, Operations.Contrast_inc)
		v1.GET("/contrastdec", Auth.AuthMiddleware, Operations.Contrast_dec)
		v1.GET("/export", Auth.AuthMiddleware, Operations.Export)
		v1.GET("/save/:pname", Auth.AuthMiddleware, Operations.Save)
		// v1.GET("/getImage", Auth.AuthMiddleware, getImage)
		v1.GET("/getdetails", Auth.AuthMiddleware, Operations.GetDetails)
		v1.GET("/resize", Auth.AuthMiddleware, Operations.Resize)
	}
	return router
}
