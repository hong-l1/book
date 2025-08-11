package config

type WeBookConfig struct {
	DbConfig
	RedisConfig
}
type DbConfig struct {
	Dns string
}
type RedisConfig struct {
	Addr string
}
