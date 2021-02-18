package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var photos = make(map[string][]byte)

func main() {
	rand.Seed(time.Now().Unix())

	router := gin.Default()

	router.Any("/api/upload", uploadPhoto)
	router.GET("/photo/:hash", showPhoto)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Only request on /api/upload",
		})
	})
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port not set")
	}
	router.Run(":" + port)
}

func uploadPhoto(c *gin.Context) {
	var url string

	longitude := c.Request.FormValue("longitude")
	latitude := c.Request.FormValue("latitude")
	if longitude == "" {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Пустое значение longitude",
		})
	}
	if latitude == "" {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Пустое значение latitude",
		})
	}

	photo, photoHeader, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer photo.Close()

	buffer := make([]byte, int(photoHeader.Size))
	if _, err := photo.Read(buffer); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	hash := createHash([]byte(strconv.Itoa(rand.Int())))
	go deletePhoto(hash)

	photos[hash] = buffer
	url = "/photo/" + hash
	c.JSON(http.StatusOK, map[string]interface{}{
		"response": map[string]string{
			"status":       "ok",
			"photoTempURL": url,
			"longitude":    longitude,
			"latitude":     latitude,
		},
	})
	// c.FormFile()
}

func showPhoto(c *gin.Context) {
	hash := c.Param("hash")
	if photo, ok := photos[hash]; ok {
		c.DataFromReader(http.StatusOK, int64(len(photo)), "image/png", bytes.NewReader(photo), nil)
	} else {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "photos not exist",
		})
	}

}

func createHash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func deletePhoto(hash string) {
	time.Sleep(time.Minute * 10)
	delete(photos, hash)
}
