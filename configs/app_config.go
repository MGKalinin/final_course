package configs

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
