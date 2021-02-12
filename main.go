package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/api/uploadphoto", uploadPhoto)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Only request on /api/uploadphoto",
		})
	})

	router.Run(os.Getenv("PORT"))
}

func uploadPhoto(c *gin.Context) {
	type RequestJSON struct {
		Latitude  float32
		Longitude float32
	}
	// c.FormFile()
}
