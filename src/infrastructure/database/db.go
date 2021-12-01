package database

import (
	"auth-service/config"
	"auth-service/src/domain/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	host := config.C.Database.Host
	user := config.C.Database.User
	pass := config.C.Database.Password
	dbname := config.C.Database.DBName
	port := fmt.Sprint(config.C.Database.Port)
	sslm := config.C.Database.SSLMode
	tz := config.C.Database.TimeZone

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, pass, dbname, port, sslm, tz)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if !db.Migrator().HasTable(&models.User{}) {
		db.Migrator().CreateTable(&models.User{})
	}

	return db, nil
}
