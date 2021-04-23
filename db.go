package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB инициализация ДБ
func (a *App) InitDB() {

	// if db, err := gorm.Open(sqlite.Open("fotocontroll.db"), &gorm.Config{}); err != nil {
	// 	panic(err)
	// } else {
	// 	a.DB = db
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
		a.DB = db
	}

	a.DB.AutoMigrate(&UserDB{}, &StorageDB{}, &ObjectDB{}, &PhotoDB{})

	// для sqlite
	// if _, err := os.Stat("fotocontroll.db"); os.IsExist(err) {
	// 	log.Println("Тестовые записи не добавлялись")
	// 	return
	// }

	if err := a.DB.Model(&UserDB{}).First(nil).Where("login = ?", "user_1").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Добавление тестовых данных")
		}
	} else {
		log.Println("Тестовые записи не добавлялись")
		return
	}

	for i := 0; i < 10; i++ {
		var u UserDB
		u = UserDB{
			FIO:      "FIO_user_" + strconv.Itoa(i),
			Login:    "user_" + strconv.Itoa(i),
			Password: "pass_" + strconv.Itoa(i),
			Storages: []StorageDB{
				StorageDB{
					NameStorage: "Storage_" + strconv.Itoa(i),
					Address:     "city,street,dom_" + strconv.Itoa(i),
					Objects: []ObjectDB{
						ObjectDB{
							NameObject: "auto_" + strconv.Itoa(i),
						},
						ObjectDB{
							NameObject: "auto_1." + strconv.Itoa(i),
						},
					},
				},
				StorageDB{
					NameStorage: "Storage_" + strconv.Itoa(i) + "." + strconv.Itoa(i+1),
					Address:     "city,street,domik_1." + strconv.Itoa(i),
					Objects: []ObjectDB{
						ObjectDB{
							NameObject: "auto_2." + strconv.Itoa(i),
						},
						ObjectDB{
							NameObject: "auto_3." + strconv.Itoa(i),
						},
					},
				},
			},
		}
		err := a.DB.Debug().Where("login = ?", u.Login).FirstOrCreate(&u).Error
		if err != nil {
			panic(err)
		}
	}
	u := new(UserDB)
	s := new(StorageDB)
	a.DB.Model(&UserDB{}).Where("id = 1").First(u)
	a.DB.Model(&StorageDB{}).Where("id = 3").First(s)
	u.Storages = append(u.Storages, *s)
	a.DB.Save(u)
}
