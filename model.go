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
	Storages []StorageDB `gorm:"many2many:user_from_storage"`
	Photo    []PhotoDB   `gorm:"foreignKey:UserID"`
}

//StorageDB  Таблица хранилищпользователя в БД
type StorageDB struct {
	gorm.Model
	NameStorage string     `gorm:"size:60"`
	Address     string     `gorm:"text"`
	Objects     []ObjectDB `gorm:"foreignKey:StorageID"`
}

// ObjectID Таблица автомобилей в хранилище
type ObjectDB struct {
	gorm.Model
	NameObject string    `gorm:"size:50`
	Photos     []PhotoDB `gorm:"foreignKey:ObjectID"` //
	StorageID  int
}

// PhotoDB Фотографии автомобилей
type PhotoDB struct {
	gorm.Model
	Path      string `gorm:"size:128`
	Longitude string
	Latitude  string
	UserID    int
	ObjectID  int
}

// User Таблица пользователей в БД
type User struct {
	ID       uint      `json:"id"`
	FIO      string    `json:"fio"`
	Storages []Storage `json:"storages,omitempty"`
}

//Storage  /
type Storage struct {
	ID          uint     `json:"id"`
	NameStorage string   `json:"name_storage"`
	Address     string   `json:"address"`
	Objects     []Object `json:"objects"`
}

// Object /
type Object struct {
	ID         uint    `json:"id"`
	NameObject string  `json:"name_object"`
	Photos     []Photo `json:"photos"`
}

// Photo /
type Photo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Longitude string    `json:"longitude"`
	Latitude  string    `json:"latitude"`
	Path      string    `json:"path"`
	User      User      `json:"user"`
}
