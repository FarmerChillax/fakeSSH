package utils

import (
	"fmt"

	"github.com/FarmerChillax/fakeSSH/config"
	"github.com/FarmerChillax/fakeSSH/model"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() (db *gorm.DB) {
	mysqlConfig := config.GetMysql()
	if mysqlConfig != nil {
		db = NewMysql(mysqlConfig)
	} else {
		db = NewSqlite(config.GetSqlite())
	}

	// 迁移 schema
	AutoMigrationSchema(db)
	return db
}

func NewSqlite(sqliteConfig *config.SqliteConfig) *gorm.DB {
	sqlitePath := "./data.db"
	if sqliteConfig != nil {
		sqlitePath = sqliteConfig.Path
	}
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		panic("[NewSqlite] failed to connect database")
	}
	logrus.Infof("connect to sqlite, path: %s", sqlitePath)

	return db
}

func NewMysql(mysqlConf *config.MysqlConfig) *gorm.DB {
	dsnTmplate := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"

	dsn := fmt.Sprintf(dsnTmplate, mysqlConf.Username, mysqlConf.Password,
		mysqlConf.Host, mysqlConf.Port, mysqlConf.DBName)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("NewMysql err: %v", err))
	}

	sqldb, err := db.DB()
	if err != nil {
		logrus.Errorf("NewMysql.db.DB err: %v", err)
		panic(err)
	}

	maxIdle := 10
	maxOpen := 30
	sqldb.SetMaxIdleConns(maxIdle)
	sqldb.SetMaxOpenConns(maxOpen)

	return db
}

func AutoMigrationSchema(db *gorm.DB) (err error) {
	logrus.Infof("Start AutoMigrationSchema.")
	// 迁移 schema
	err = db.AutoMigrate(&model.Data{})
	if err != nil {
		logrus.Errorf("db.AutoMigrate(&model.Data{}) err: %v", err)
		return err
	}
	logrus.Infof("AutoMigrationSchema done.")
	return nil
}
