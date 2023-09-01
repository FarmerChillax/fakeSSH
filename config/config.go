package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Mysql  *MysqlConfig  `json:"mysql,omitempty"`
	Sqlite *SqliteConfig `json:"sqlite,omitempty"`
}

var (
	config *Config
	ENV    string
)

func Load() error {
	ENV = os.Getenv("ENV")
	log.Infof("loading env: %v", ENV)
	switch ENV {
	case "dev":
		viper.SetConfigName("dev.config")
	case "test":
		viper.SetConfigName("test.config")
	default:
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml")
	viper.SetConfigType("json")
	viper.AddConfigPath("./local")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

type MysqlConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int32  `json:"port,omitempty"`
	DBName   string `json:"db_name,omitempty" mapstructure:"db_name"`
}

func GetMysql() *MysqlConfig {
	return config.Mysql
}

type SqliteConfig struct {
	Path string `json:"path,omitempty"`
}

func GetSqlite() *SqliteConfig {
	return config.Sqlite
}
