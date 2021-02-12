package main

import (
	"log"
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
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port not set")
	}
	router.Run(":" + port)
}

func uploadPhoto(c *gin.Context) {
	type RequestJSON struct {
		Latitude  float32
		Longitude float32
	}

	c.JSON(http.StatusOK, map[string]string{
		"response": "OK",
	})
	// c.FormFile()
}
