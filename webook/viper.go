package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigFile("webook/config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
