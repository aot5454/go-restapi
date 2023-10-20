package config

import (
	"fmt"
	"go-restapi/app"
	"strings"

	"github.com/spf13/viper"
)

func GetConfig() (config app.Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	fmt.Printf("%#v\n", config)
	return
}
