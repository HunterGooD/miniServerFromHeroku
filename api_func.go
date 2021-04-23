package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (a *App) getStorageByID(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	id := c.Param("id")
	var storageDB StorageDB
	if err := a.DB.Model(&StorageDB{}).Where("id = ?", id).Preload("Objects.Photos").First(&storageDB).Error; err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"error": "Ошибка получения записей",
		})
		return
	}
	obj := make([]Object, len(storageDB.Objects))
	for i, o := range storageDB.Objects {
		phs := make([]Photo, len(o.Photos))
		for j, p := range o.Photos {
			phs[j] = Photo{
				ID:        p.ID,
				Path:      p.Path,
				Longitude: p.Longitude,
				Latitude:  p.Latitude,
				CreatedAt: p.CreatedAt,
			}
		}
		obj[i] = Object{
			ID:         o.ID,
			NameObject: o.NameObject,
			Photos:     phs,
		}
	}
	c.JSON(http.StatusOK, obj)
}

func (a *App) getStoragesInfo(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	res := make([]User, 0)
	rows, err := a.DB.Debug().Model(&UserDB{}).Rows()
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"error": "Ошибка получения записей",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ag UserDB
		a.DB.ScanRows(rows, &ag)
		var agent UserDB
		a.DB.Debug().Model(&UserDB{}).Where("id = ?", ag.ID).Preload("Storages").First(&agent)
		storages := make([]Storage, len(agent.Storages))
		for i, s := range agent.Storages {
			storages[i] = Storage{
				ID:          s.ID,
				NameStorage: s.NameStorage,
				Address:     s.Address,
			}
		}
		res = append(res, User{
			ID:       agent.ID,
			FIO:      agent.FIO,
			Storages: storages,
		})
	}
	c.JSON(http.StatusOK, res)
}

func (a *App) showPhotos(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := make([]string, len(photos))
	var index int
	for i := range photos {
		res[index] = i
		index++
	}
	c.JSON(http.StatusOK, res)
}

func (a *App) getAllInfo(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := make([]UserDB, 0)
	rows, err := a.DB.Debug().Model(&UserDB{}).Rows()
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"error": "Ошибка получения записей",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ag UserDB
		a.DB.ScanRows(rows, &ag)
		var addAgent UserDB
		a.DB.Debug().Model(&UserDB{}).Where("id = ?", ag.ID).Preload("Storages.Objects.Photos").Preload(clause.Associations).First(&addAgent)
		res = append(res, addAgent)
	}
	c.JSON(http.StatusOK, res)
}

func (a *App) getStoragesAgent(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	type AutoS struct {
		ID       int    `json:"id"`
		NameAuto string `json:"name_auto"`
	}

	type StorageS struct {
		ID          int     `json:"id"`
		NameStorage string  `json:"name_storage"`
		Address     string  `json:"address"`
		Autos       []AutoS `json:"autos"`
	}

	type Response struct {
		Storages []StorageS `json:"storages"`
	}

	userToken := c.GetHeader("Authorization") // проверить кто это и записать

	if err := a.DB.Model(&UserDB{}).Where("token = ?", userToken).First(nil).Error; err != nil {
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
	u := new(UserDB)
	// TODO: передалть на Preload
	err := a.DB.Debug().Model(&UserDB{}).Where("user_dbs.token = ?", userToken).Preload("Storages.Objects").First(u).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка приложения",
			"err":   err.Error(),
		})
		return
	}

	res := new(Response)

	for _, storage := range u.Storages {
		stor := new(StorageS)
		autos := make([]AutoS, len(storage.Objects))
		for _, a := range storage.Objects {
			autos = append(autos, AutoS{
				ID:       int(a.ID),
				NameAuto: a.NameObject,
			})
		}
		stor = &StorageS{
			ID:          int(storage.ID),
			NameStorage: storage.NameStorage,
			Address:     storage.Address,
			Autos:       autos,
		}
		res.Storages = append(res.Storages, *stor)
	}

	c.JSON(http.StatusOK, res)
}

func (a *App) uploadPhoto(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var url string
	userToken := c.GetHeader("Authorization") // проверить кто это и записать
	var agent UserDB
	if err := a.DB.Model(&UserDB{}).Where("token = ?", userToken).Preload("Storages.Objects.Photos").Preload(clause.Associations).First(&agent).Error; err != nil {
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
	id_object, err := strconv.Atoi(c.Request.FormValue("object_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка при опозновании склада",
		})
		return
	}
	id_storage, err := strconv.Atoi(c.Request.FormValue("storage_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка при опозновании склада",
		})
		return
	}
	if longitude == "" {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Пустое значение longitude",
		})
		return
	}
	if latitude == "" {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Пустое значение latitude",
		})
		return
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
	for i, s := range agent.Storages {
		if int(s.ID) == id_storage {
			for j, o := range s.Objects {
				if int(o.ID) == id_object {
					agent.Storages[i].Objects[j].Photos = append(agent.Storages[i].Objects[j].Photos, PhotoDB{
						Path:      hash,
						Longitude: longitude,
						Latitude:  latitude,
					})
					break
				}
			}
			break
		}
	}
	//TODO: раскомментировать
	// filePhoto, err := os.Create(hash + ".png")
	// if err != nil {
	// 	c.JSON(http.StatusOK, map[string]string{
	// 		"error": "Ошибка загрузки фотографии",
	// 	})
	// 	return
	// }
	// defer filePhoto.Close()

	// if _, err := filePhoto.Write(buffer); err != nil {
	// 	c.JSON(http.StatusOK, map[string]string{
	// 		"error": "Ошибка загрузки фотографии",
	// 	})
	// 	return
	// }

	if err := a.DB.Save(agent).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, map[string]string{
			"error": "Ошибка добавления в базу",
		})
		return
	}

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

func (a *App) showPhoto(c *gin.Context) {
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

func (a *App) deletePhoto(hash string) {
	log.Printf("/photo/%s", hash)
	time.Sleep(time.Minute * 10)
	delete(photos, hash)
}

// ValidPassword проверка валидности пароля
func (a *App) ValidPassword(password string) bool {
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
func (a *App) ValidLogin(login string) bool {
	return regexp.MustCompile(`^([A-z]+([-_]?[A-z0-9]+){0,2}){4,32}$`).Match([]byte(login))
}

func (a *App) signIn(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	type user struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	u := new(user)
	c.BindJSON(u)
	var udb UserDB
	// без проверки хэша
	if err := a.DB.Model(&UserDB{}).Where(map[string]interface{}{"login": u.Login, "password": u.Password}).First(&udb).Error; err != nil {
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
	if err := a.DB.Debug().Model(&UserDB{}).Where("login = ?", u.Login).Update("token", token).Error; err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"error": "Ошибка при обновлении данных",
			"err":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"fio":   udb.FIO,
	})
}
