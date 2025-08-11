package main

import (
	"fmt"
	"github.com/petermattis/goid"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type RecursiveMutex struct {
	sync.Mutex       // 内嵌标准互斥锁
	owner      int64 // 当前持有锁的goroutine ID
	recursion  int32 // 重入次数计数器
}
type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

// GoID 拿到当前goroutine的id
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
func (m *RecursiveMutex) Lock() {
	gid := goid.Get() // 获取当前goroutine ID

	// 检查是否是当前持有锁的goroutine再次尝试获取锁
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++ // 增加重入计数
		return        // 直接返回，不实际获取锁
	}

	// 新goroutine尝试获取锁
	m.Mutex.Lock() // 阻塞直到获取底层锁

	// 记录新持有者信息
	atomic.StoreInt64(&m.owner, gid) // 原子存储goroutine ID
	m.recursion = 1                  // 初始化重入计数为1
}
func (m *RecursiveMutex) Unlock() {
	gid := goid.Get() // 获取当前goroutine ID
	// 检查是否是锁的持有者尝试释放
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	// 减少重入计数
	m.recursion--
	// 检查是否完全释放
	if m.recursion != 0 { // 还有嵌套锁未释放
		return
	}
	// 完全释放锁
	atomic.StoreInt64(&m.owner, -1) // 清除持有者信息
	m.Mutex.Unlock()                // 释放底层互斥锁
}
func (m *TokenRecursiveMutex) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token { //如果传入的token和持有锁的token一致，说明是递归调用
		m.recursion++
		return
	}
	m.Mutex.Lock() // 传入的token不一致，说明不是递归调用
	// 抢到锁之后记录这个token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

// 释放锁
func (m *TokenRecursiveMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token { // 释放其它token持有的锁
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.token, token))
	}
	m.recursion--         // 当前持有这个锁的token释放锁
	if m.recursion != 0 { // 还没有回退到最初的递归调用
		return
	}
	atomic.StoreInt64(&m.token, 0) // 没有递归调用了，释放锁
	m.Mutex.Unlock()
}
