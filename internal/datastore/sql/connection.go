package datastore

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

type UserCredential struct {
	Id          int    `gorm:"primaryKey"`
	PhoneNumber string `gorm:"not null"`
	Pin         string `gorm:"not null"`
	SharedKey   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func InitDb() {
	fmt.Println("start initiating the db")
	var err error
	Db, err = gorm.Open(sqlite.Open("./biometric.dat"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Db.AutoMigrate(&UserCredential{})
	fmt.Println("init db done.")
}
