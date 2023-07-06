package linkcxo

import (
	"os"

	"log"

	"github.com/spf13/viper"
)

func InitializeEnvironment() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
	keys := viper.AllKeys()

	for _, k := range keys {
		envValue, found := os.LookupEnv(k)
		if !found || envValue == "" {
			os.Setenv(k, viper.GetString(k))
		}
	}
}
