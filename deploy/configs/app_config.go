package configs

import (
	"github.com/spf13/viper"
	"log"
	"path"
	"runtime"
)

func LoadConfig() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)

	viper.AddConfigPath(dir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
}
