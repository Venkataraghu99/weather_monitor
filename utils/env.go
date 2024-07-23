package utils

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	ApiKey         string `mapstructure:"API_KEY"`
	OpenWeatherUrl string `mapstructure:"OPEN_WEATHER_URL"`
	Port           int    `mapstructure:"PORT"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}
