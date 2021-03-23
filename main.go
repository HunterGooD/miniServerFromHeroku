package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	photos = make(map[string][]byte)
	//DB соединение с БД
	DB *gorm.DB
)

// UserDB Таблица пользователей в БД
type UserDB struct {
	gorm.Model
	Login    string      `gorm:"size:30"`
	Password string      `gorm:"size:60"`
	Token    string      `gorm:"size:64;index"`
	Storages []StorageDB `gorm:"foreignKey:UserID"`
}

//StorageDB  Таблица хранилищпользователя в БД
type StorageDB struct {
	gorm.Model
	NameStorage string   `gorm:"size:60"`
	Address     string   `gorm:"text"`
	Autos       []AutoDB `gorm:"foreignKey:StorageID"`
	UserID      int
}

// AutoDB Таблица автомобилей в хранилище
type AutoDB struct {
	gorm.Model
	NameAuto  string    `gorm:"size:50`
	Photos    []PhotoDB `gorm:"many2many:auto_photos"`
	StorageID int
}

// PhotoDB Фотографии автомобилей
type PhotoDB struct {
	gorm.Model
	Path string `gorm:"size:128`
	Date time.Time
}

// InitDB инициализация ДБ
func InitDB() {

	// var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
	// 	"root",
	// 	"gfhjkm",
	// 	"localhost",
	// 	"3306",
	// 	"fotocontroll",
	// )
	// if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
	// 	panic(err)
	// } else {
	// 	DB = db
	// }

	// if db, err := gorm.Open(sqlite.Open("forocontroll.db"), &gorm.Config{}); err != nil {
	// 	panic(err)
	// } else {
	// 	DB = db
	// }

	urlDB := os.Getenv("DATABASE_URL")
	if urlDB == "" {
		panic("Error url for database  not found")
	}

	spDB := strings.Split(urlDB, "://")[1]
	infU := strings.Split(spDB, ":")
	user := infU[0]
	pass := strings.Split(infU[1], "@")[0]
	infH := strings.Split(urlDB, "@")[1]
	infoHost := strings.Split(infH, ":")
	host := infoHost[0]
	port := strings.Split(infoHost[1], "/")[0]
	dbName := strings.Split(infoHost[1], "/")[1]

	var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Europe/Samara",
		host,
		user,
		pass,
		dbName,
		port,
	)

	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Printf("%v", err)
		panic(err)
	} else {
		DB = db
	}

	DB.AutoMigrate(&UserDB{}, &StorageDB{}, &AutoDB{}, &PhotoDB{})

	if _, err := os.Stat("forocontroll.db"); os.IsExist(err) {
		return
	}

	for i := 0; i < 10; i++ {
		u := &UserDB{
			Login:    "user_" + strconv.Itoa(i),
			Password: "pass_" + strconv.Itoa(i),
			Storages: []StorageDB{
				StorageDB{
					NameStorage: "Storage_" + strconv.Itoa(i),
					Address:     "city,street,dom_" + strconv.Itoa(i),
					Autos: []AutoDB{
						AutoDB{
							NameAuto: "auto_" + strconv.Itoa(i),
						},
						AutoDB{
							NameAuto: "auto_1." + strconv.Itoa(i),
						},
					},
				},
				StorageDB{
					NameStorage: "Storage_" + strconv.Itoa(i),
					Address:     "city,street,domik_1." + strconv.Itoa(i),
					Autos: []AutoDB{
						AutoDB{
							NameAuto: "auto_2." + strconv.Itoa(i),
						},
						AutoDB{
							NameAuto: "auto_3." + strconv.Itoa(i),
						},
					},
				},
			},
		}

		r := DB.Create(u)
		if r.Error != nil {
			panic(r.Error)
		}
	}
}

func main() {
	InitDB()

	router := gin.Default()

	router.Any("/api/upload", uploadPhoto)           // загрузка фото
	router.POST("/api/signin", signIn)               // авторизация в системе
	router.GET("/photo/:hash", showPhoto)            // получение фотографий
	router.GET("/agents", nil)                       // показывает всех агентов
	router.GET("/api/getStorages", getStoragesAgent) // показывает все склады агента
	// router.GET("/agent/:id", nil)             // показывает скалды агентов
	// router.GET("/agent/:id/storages", nil)    // список вех складов агента
	// router.GET("/agent/:id/storage/:id", nil) // список автомобилей с последней датой обновления

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

func getStoragesAgent(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	type AutoS struct {
		ID       int64  `json:"id"`
		NameAuto string `json:"name_auto"`
	}

	type StorageS struct {
		ID          int64   `json:"id"`
		NameStorage string  `json:"name_storage"`
		Address     string  `json:"address"`
		Autos       []AutoS `json:"autos"`
	}

	type Response struct {
		Storages []StorageS `json:"storages"`
	}

	userToken := c.GetHeader("Authorization") // проверить кто это и записать

	if err := DB.Model(&UserDB{}).Where("token = ?", userToken).First(nil).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Ошибка авторизации",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Ошибка приложения",
		})
		return
	}

	rows, err := DB.Debug().Table("storage_dbs").
		Select("storage_dbs.id as storage_id, storage_dbs.name_storage, storage_dbs.address, auto_dbs.id as auto_id, auto_dbs.name_auto").
		Joins("inner join user_dbs on storage_dbs.user_id = user_dbs.id").
		Joins("inner join auto_dbs on storage_dbs.id = auto_dbs.storage_id").
		Where("user_dbs.token = ?", userToken).Rows()

	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Ошибка приложения",
		})
		return
	}
	defer rows.Close()

	res := new(Response)
	var currentID int64

	for rows.Next() {
		var s = make(map[string]interface{})
		DB.ScanRows(rows, s)

		if s["storage_id"] != currentID {
			currentID = s["storage_id"].(int64)

			res.Storages = append(res.Storages, StorageS{
				ID:          currentID,
				NameStorage: s["name_storage"].(string),
				Address:     s["address"].(string),
				Autos: []AutoS{
					AutoS{
						ID:       s["auto_id"].(int64),
						NameAuto: s["name_auto"].(string),
					},
				},
			})
			continue
		}

		for i := range res.Storages {
			if res.Storages[i].ID == currentID {
				res.Storages[i].Autos = append(res.Storages[i].Autos, AutoS{
					ID:       s["auto_id"].(int64),
					NameAuto: s["name_auto"].(string),
				})
				break
			}
		}
	}
	c.JSON(http.StatusOK, res)
}

func uploadPhoto(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var url string
	userToken := c.GetHeader("Authorization") // проверить кто это и записать

	if err := DB.Model(&UserDB{}).Where("token = ?", userToken).First(nil).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Ошибка авторизации",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Ошибка приложения",
		})
		return
	}

	longitude := c.Request.FormValue("longitude")
	latitude := c.Request.FormValue("latitude")
	// id_auto := c.Request.FormValue("id_auto")
	// id_storage := c.Request.FormValue("id_storage")

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

	// TODO: загружать картинки в папку photos/ с уникальным именем. и записывать в базу
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
	h := sha256.New()
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
	if len(ps) < 6 {
		return fmt.Errorf("password len is < 6")
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
	c.Header("Content-Type", "application/json")
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
			c.JSON(http.StatusOK, map[string]interface{}{
				"error": "Запись не найдена. Проверьте данные",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Ошибка приложения",
		})
		return
	}

	t := make([]byte, 5)
	rand.Read(t)
	token := createHash(t)
	if err := DB.Debug().Model(&UserDB{}).Where("login = ?", u.Login).Update("token", token).Error; err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"error": "Ошибка при обновлении данных",
			"err":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
