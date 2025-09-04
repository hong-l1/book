package ioc

import (
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	promsdk "github.com/prometheus/client_golang/prometheus"
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
	pcb := NewCallbaks()
	pcb.registerAll(db)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	//err = dao.InitTables(db)
	//if err != nil {
	//	panic(err)
	//}
	return db
}

type gormLoggerFunc func(msg string, field ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

type Callbacks struct {
	vector *promsdk.SummaryVec
}

func NewCallbaks() *Callbacks {
	vector := promsdk.NewSummaryVec(promsdk.SummaryOpts{
		Namespace: "book",
		Subsystem: "webook",
		Name:      "gorm_query_time",
		Help:      "统计 GORM 的执行时间",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})
	promsdk.MustRegister(vector)
	return &Callbacks{
		vector: vector,
	}
}
func (c Callbacks) registerAll(db *gorm.DB) {
	err := db.Callback().Create().Before("*").Register("prometheus_create_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").Register("prometheus_create_after", c.after("create"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Before("*").Register("prometheus_Update_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").Register("prometheus_Update_after", c.after("Update"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Delete().Before("*").Register("prometheus_delete_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").Register("prometheus_delete_after", c.after("delete"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Raw().Before("*").Register("prometheus_raw_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").Register("prometheus_raw_after", c.after("raw"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Row().Before("*").Register("prometheus_row_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").Register("prometheus_row_after", c.after("row"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Query().Before("*").Register("prometheus_Query_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Query().After("*").Register("prometheus_Query_after", c.after("Query"))
	if err != nil {
		panic(err)
	}
}
func (c Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("startTime", startTime)
	}
}
func (c Callbacks) after(t string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		data, _ := db.Get("startTime")
		d, ok := data.(time.Time)
		if !ok {
			return
		}
		table := db.Statement.Table
		if len(table) == 0 {
			table = "unknown"
		}
		c.vector.WithLabelValues(t, table).Observe(float64(time.Since(d).Milliseconds()))
	}
}
