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
	ID       uint
	FIO      string
	Login    string
	Storages []Storage
}

//Storage  /
type Storage struct {
	NameStorage string
	Address     string
	Objects     []Object
}

// Object /
type Object struct {
	NameAuto string
	Photos   []PhotoDB
}

// Photo /
type Photo struct {
	CreatedAt time.Time
	Longitude string
	Latitude  string
	Path      string
}
