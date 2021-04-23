package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

var photos = make(map[string][]byte)

func main() {

	rand.Seed(time.Now().Unix())
	app := new(App)
	app.InitDB()
	if err := os.Mkdir("images/", 0700); err != nil {
		if os.IsNotExist(err) {
			panic(err)
		}
	}

	htmlFiles := packr.New("htmlFiles", "./web/dist")
	assetsFile := packr.New("assets", "./web/dist/assets")
	router := gin.Default()

	data, err := htmlFiles.FindString("index.tmpl")
	if err != nil {
		panic("html not found")
	}
	render := multitemplate.New()
	render.AddFromString("index.tmpl", data)

	webAppURL := "/web"
	router.HTMLRender = render

	router.Any("/api/upload", app.uploadPhoto)           // загрузка фото
	router.POST("/api/signin", app.signIn)               // авторизация в системе
	router.GET("/photo/:hash", app.showPhoto)            // получение фотографийd
	router.GET("/photos", app.showPhotos)                // получение фотографий
	router.GET("/api/storages", app.getStoragesInfo)     // показывает все склады
	router.GET("/api/allInfo", app.getAllInfo)           // показывает всю информацию
	router.GET("/api/storage/:id", app.getStorageByID)   // показывает объекты склада с его фотографиями
	router.GET("/api/getStorages", app.getStoragesAgent) // показывает все склады агента
	router.GET(webAppURL+"/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"webAppURL": webAppURL,
		})
	})
	router.StaticFS(webAppURL+"/assets", assetsFile)
	// router.GET("/agent/:id", nil)             // показывает скалды агентов
	// router.GET("/agent/:id/storages", nil)    // список вех складов агента
	// router.GET("/agent/:id/storage/:id", nil) // список автомобилей с последней датой обновления

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "-_- такого тут нет",
		})
	})
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port not set")
	}
	router.Run(":" + port)
}
