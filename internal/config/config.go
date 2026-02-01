package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Load() error {

	//get from env varible
	viper.AutomaticEnv()

	// change DB.HOST -> DB_HOST
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// if already decalre in env
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		
		viper.SetConfigType("env")

		if err:= viper.ReadInConfig(); err != nil{
			return  err
		}
	}

	return nil


}