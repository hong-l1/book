package data

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {
	nums := make([]int, 10)
	for k := range nums {
		nums[k] = k
	}
	stack := make([]int, 0) //先进后出
	queue := make([]int, 0) //先进先出
	for k := range nums {
		stack = append(stack, nums[k])
		queue = append(queue, k)
	}
	for len(stack) > 0 {
		fmt.Println(stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	for len(queue) > 0 {
		fmt.Println(queue[0])
		queue = queue[1:]
	}
	//flag := 1
	//for len(queue) > 0 {
	//	if flag%2 == 0 {
	//		fmt.Println(queue[0])
	//		queue = queue[1:]
	//
	//	} else {
	//		fmt.Println(queue[:len(queue)-1])
	//		queue = queue[:len(queue)-1]
	//	}
	//	flag++
	//}
	type priorityqueue struct {
		val  int
		rank int
	}
	priorityQueue := make([]priorityqueue, 0)
	for k := range nums {
		suiji := rand.Intn(100)
		priorityQueue = append(priorityQueue, priorityqueue{
			val:  k,
			rank: suiji,
		})
	}
	sort.Slice(priorityQueue, func(i, j int) bool {
		return priorityQueue[i].rank > priorityQueue[j].rank
	})
	for len(priorityQueue) > 0 {
		fmt.Println(priorityQueue[0].val)
		priorityQueue = priorityQueue[1:]
	}
}

//队列，栈，优先队列，双端队列
