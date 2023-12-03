package database

import (
	"errors"
	"fmt"
	"log"

	conf "github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func AutoMigrate() {
	DB.AutoMigrate(models.Account{})
	/*
		DB.AutoMigrate(models.ForumThreadsContent{})
		DB.AutoMigrate(models.ForumThreads{})
		DB.AutoMigrate(models.Forum{})
	*/
	DB.AutoMigrate(models.Player{})
	DB.AutoMigrate(models.Todays{}, models.Comments{}, models.Voted{})
	DB.AutoMigrate(models.NewsTicket{})
	DB.AutoMigrate(models.News{})
	DB.AutoMigrate(models.Files{})
	DB.AutoMigrate(models.Forum{}, models.ForumThreads{}, models.ForumThreadsContent{})

}

func connectionSqlite(database conf.DataBase, debugMode bool) error {
	var err error

	logDebug := logger.Silent
	if debugMode {
		logDebug = logger.Warn
	}

	DB, err = gorm.Open(sqlite.Open(database.SqlitePath), &gorm.Config{
		Logger: logger.Default.LogMode(logDebug),
	})
	if err != nil {
		return err
	}
	AutoMigrate()
	log.Println("[SQLITE] Connection Sqlite ok")
	return nil

}

func connectionMysql(database conf.DataBase, debugMode bool) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		database.User, database.Password, database.Host, database.Port, database.NameDB)

	logDebug := logger.Silent
	if debugMode {
		logDebug = logger.Warn
	}
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logDebug),
	})
	if err != nil {
		return err
	}
	AutoMigrate()
	log.Println("[MYSQL] Connection Mysql Ok")
	return nil
}

func ConnectionDataBase() error {
	switch conf.Config.DataBase.Engine {
	case "sqlite":
		return connectionSqlite(conf.Config.DataBase, conf.Config.Server.Debug)

	case "mysql":
		return connectionMysql(conf.Config.DataBase, conf.Config.Server.Debug)

	default:
		return errors.New(conf.ERROR_DATABASE_GET)
	}

}
