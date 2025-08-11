package main

import (
	"fmt"
	"sort"
)

func sack01(n int, c int, v []int, w []int) int {
	ans := make([]int, c+1)
	ans[0] = 0
	for i := 0; i < n; i++ {
		for j := c; j >= v[i]; j-- {
			ans[j] = max(ans[j], ans[j-v[i]]+w[i])
		}
	}
	sort.Ints(ans)
	return ans[c]
}
func main() {
	n := 4
	C := 5
	v := []int{1, 3, 4, 2}
	w := []int{15, 20, 30, 10}
	fmt.Println("最大价值为:", sack01(n, C, v, w))
}
