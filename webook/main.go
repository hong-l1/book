package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	//InitViperRemote()
	//keys := viper.AllKeys()
	//fmt.Println(keys)
	//setting := viper.AllSettings()
	//fmt.Println(setting)
	//Initviper11()
	InitPrometheus()
	app := InitWebServer()
	for _, c := range app.Consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	app.Server.Run(":8080")
}
func initViper() {
	viper.SetDefault("db.mysql.dsn", "root:123456@tcp(localhost:6380)/webook?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
func initViperv1() {
	viper.SetConfigFile("./config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
func Initviper11() {
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(e.Name, e.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
func InitViperRemote() {
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.SetConfigType("yaml")
	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.ReadRemoteConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(e.Name, e.Op)
	})
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
func InitPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}
