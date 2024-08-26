package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"jwt-auth/models"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "secret" //Enter your password for the DB
	dbname   = "go-jwt-auth"
)

var dsn string = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
	host, port, user, password, dbname)

var DB *gorm.DB

func DBconn() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	db.AutoMigrate(&models.Users{})
	db.AutoMigrate(&models.Sessions{})
}
