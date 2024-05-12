package main

import (
	"Diploma/internal/handler"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env %s", err.Error())
	}

	if err := handler.Start(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
