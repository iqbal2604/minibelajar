package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var App = struct {
	Mode        string
	Port        string
	PrefixApi   string
	Environment string
}{}

var DB = struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Name     string
	SSL      string
	Timezone string
}{}

var Redis = struct {
	Host     string
	Password string
}{}

var RabbitMQ = struct {
	Host     string
	Username string
	Password string
}{}

var Token = struct {
	AccessSecretKey string
}{}

func Load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found, using environment variables")
	}

	App.Mode = os.Getenv("MODE")
	App.Port = os.Getenv("PORT")
	App.PrefixApi = os.Getenv("PREFIX_API")
	App.Environment = os.Getenv("ENVIRONMENT")

	if App.Port == "" {
		log.Panic("PORT cannot be empty")
	}

	if App.PrefixApi == "" {
		log.Panic("PREFIX_API cannot be empty")
	}

	if App.Mode != "debug" && App.Mode != "release" && App.Mode != "test" {
		log.Panicf("invalid MODE: %s", App.Mode)
	}

	DB.Host = os.Getenv("DB_HOST")
	DB.Port = os.Getenv("DB_PORT")
	DB.User = os.Getenv("DB_USER")
	DB.Pass = os.Getenv("DB_PASS")
	DB.Name = os.Getenv("DB_NAME")
	DB.SSL = os.Getenv("DB_SSL")
	DB.Timezone = os.Getenv("DB_TIMEZONE")

	if DB.Host == "" || DB.Port == "" || DB.User == "" || DB.Pass == "" || DB.Name == "" || DB.SSL == "" || DB.Timezone == "" {
		log.Panic("DB config incomplete")
	}

	Redis.Host = os.Getenv("REDIS_HOST")
	Redis.Password = os.Getenv("REDIS_PASSWORD")

	RabbitMQ.Host = os.Getenv("RABBITMQ_HOST")
	RabbitMQ.Username = os.Getenv("RABBITMQ_USERNAME")
	RabbitMQ.Password = os.Getenv("RABBITMQ_PASSWORD")

	Token.AccessSecretKey = os.Getenv("TOKEN_ACCESS_SECRET_KEY")
}
