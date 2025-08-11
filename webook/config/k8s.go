//go:build k8s

package config

var Config = WeBookConfig{
	DbConfig{
		Dns: "root:123456@tcp(webook-mysql:6380)/webook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	RedisConfig{
		Addr: "webook-redis:3308",
	},
}
