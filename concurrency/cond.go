package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	c := sync.NewCond(&sync.Mutex{})
	var ready int

	// 启动 10 个运动员 goroutine
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			c.L.Lock()
			ready++
			c.L.Unlock()
			log.Printf("运动员 #%d 已准备", i)
			c.Broadcast()
		}(i)
	}
	// 裁判 A
	go func() {
		c.L.Lock()
		for ready != 10 {
			c.Wait()
			log.Println("裁判 A 被唤醒一次")
		}
		log.Println("裁判 A：所有运动员准备好了，开始！")
		c.L.Unlock()
	}()

	// 裁判 B
	go func() {
		c.L.Lock()
		for ready != 10 {
			c.Wait()
			log.Println("裁判 B 被唤醒一次")
		}
		log.Println("裁判 B：所有运动员准备好了，开始！")
		c.L.Unlock()
	}()

	time.Sleep(5 * time.Second) // 等待 goroutine 完成
}
