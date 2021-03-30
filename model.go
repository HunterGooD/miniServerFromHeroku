package main

import (
	"time"

	"gorm.io/gorm"
)

// UserDB Таблица пользователей в БД
type UserDB struct {
	gorm.Model
	FIO      string      `gorm:"size:90"`
	Login    string      `gorm:"size:30;unique"`
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
	Path      string `gorm:"size:128`
	Longitude string
	Latitude  string
}

// User Таблица пользователей в БД
type User struct {
	ID       uint      `json:"id"`
	FIO      string    `json:"fio"`
	Storages []Storage `json:"storages"`
}

//Storage  /
type Storage struct {
	NameStorage string   `json:"name_storage"`
	Address     string   `json:"address"`
	Objects     []Object `json:"objects"`
}

// Object /
type Object struct {
	NameObject string    `json:"name_object"`
	Photos     []PhotoDB `json:"photos"`
}

// Photo /
type Photo struct {
	CreatedAt time.Time `json:"created_at"`
	Longitude string    `json:"longitude"`
	Latitude  string    `json:"latitude"`
	Path      string    `json:"path"`
}
