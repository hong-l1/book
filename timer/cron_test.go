package timer

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	expr.AddFunc("@every 1s", func() {
		fmt.Println("运行！")
		time.Sleep(5 * time.Second)
		fmt.Println("结束！")
	})
	expr.Start()
	time.Sleep(3 * time.Second)
	c := expr.Stop()
	fmt.Println("发出停止信号")
	<-c.Done()
	fmt.Println("全部结束")
}
