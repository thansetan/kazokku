package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type DB struct {
	UserName string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
}

type App struct {
	Host    string `mapstructure:"APP_HOST"`
	Port    int    `mapstructure:"APP_PORT"`
	SaveDir string `mapstructure:"SAVE_DIR"`
}

type Config struct {
	Database DB
	App      App
}

func LoadConfig(configFilePath string) (Config, error) {
	var conf Config
	var dbConf DB
	var appConf App

	_, err := os.Stat(configFilePath)
	if err != nil {
		return conf, err
	}

	v := viper.New()

	v.SetConfigFile(configFilePath)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return conf, err
	}

	if err := v.Unmarshal(&dbConf); err != nil {
		return conf, err
	}

	if err := v.Unmarshal(&appConf); err != nil {
		return conf, err
	}

	conf.Database = dbConf
	conf.App = appConf
	os.Setenv("SAVE_DIR", appConf.SaveDir)

	return conf, nil
}

func (db DB) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.UserName, db.Password, db.Host, db.Port, db.Name)
}
