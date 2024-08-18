package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Join(filepath.Dir(b), "..")
	Env        envSchema
)

type dbConnection struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	SslMode  string `mapstructure:"DB_SSL_MODE"`
}

func (d dbConnection) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", d.Host, d.Port, d.Name, d.Username, d.Password, d.SslMode)
}

type envSchema struct {
	AppPort   string `mapstructure:"APP_PORT"`
	Db        dbConnection
	SecretKey string `mapstructure:"SECRET_KEY"`
}

func Init() {
	Env = envSchema{}
	dbConn := dbConnection{}
	viper.SetConfigType("env")
	viper.AddConfigPath(basepath) // path to look for the config file in
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&Env)

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&dbConn)

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	Env.Db = dbConn
}
