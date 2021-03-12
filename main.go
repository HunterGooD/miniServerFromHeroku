package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	photos = make(map[string][]byte)
	DB     *gorm.DB
)

type UserDB struct {
	gorm.Model
	Login    string `gorm:"size:30"`
	Password string `gorm:"size:60"`
}

func main() {

	if db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{}); err != nil {
		panic(err)
	} else {
		DB = db
	}
	DB.AutoMigrate(&UserDB{})
	for i := 0; i < 10; i++ {
		u := &UserDB{
			Login:    "user_" + strconv.Itoa(i),
			Password: "pass_" + strconv.Itoa(i),
		}
		r := DB.Model(&UserDB{}).Create(u)
		if r.Error != nil {
			panic(r.Error)
		}
	}

	rand.Seed(time.Now().Unix())

	router := gin.Default()

	router.Any("/api/upload", uploadPhoto)
	router.POST("/api/signin", signIn)
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
	log.Printf("/photo/%s", hash)
	time.Sleep(time.Minute * 10)
	delete(photos, hash)
}

// ValidPassword проверка валидности пароля
func ValidPassword(password string) bool {
	return CheckPasswordLever(password) == nil
}

//CheckPasswordLever Сложность пароля
func CheckPasswordLever(ps string) error {
	if len(ps) < 9 {
		return fmt.Errorf("password len is < 9")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("password need num :%v", err)
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return fmt.Errorf("password need a_z :%v", err)
	}
	if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
		return fmt.Errorf("password need A_Z :%v", err)
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		return fmt.Errorf("password need symbol :%v", err)
	}
	return nil
}

// ValidLogin проверка валидности логина
func ValidLogin(login string) bool {
	return regexp.MustCompile(`^([A-z]+([-_]?[A-z0-9]+){0,2}){4,32}$`).Match([]byte(login))
}

func signIn(c *gin.Context) {
	type user struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	u := new(user)
	c.BindJSON(u)
	var udb UserDB
	// без проверки хэша
	if err := DB.Model(&UserDB{}).Where(map[string]interface{}{"login": u.Login, "password": u.Password}).First(&udb).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Запись не найдена. Проверьте данные",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Ошибка приложения",
		})
		return
	}
	// пока бесполезно
	t := make([]byte, 5)
	rand.Read(t)
	token := createHash(t)
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
