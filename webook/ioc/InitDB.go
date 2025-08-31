package ioc

import (
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
	"log"
	"time"
)

func InitDb(l logger.Loggerv1) *gorm.DB {
	type config struct {
		Dsn string `mapstructure:"dsn"`
	}
	cfg := config{
		Dsn: "root:123456@tcp(localhost:6380)/webook?charset=utf8mb4&parseTime=True&loc=Local",
	}
	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
		QueryFields:            false,
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			//慢查询阈值，执行时间超过阈值，才会使用
			//50，100ms
			SlowThreshold:             time.Millisecond * 10,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  glogger.Info,
			ParameterizedQueries:      true,
		}),
	})
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, field ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}
