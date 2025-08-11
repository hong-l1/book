//package main
//
//import "fmt"
//
//func lengthOfLIS(nums []int) int {
//	l := len(nums)
//	ans := make([]int, l)
//	list := make([]int, l)
//	for temp := range ans {
//		ans[temp] = 1
//		list[temp] = -1
//	}
//	for k := 0; k < l; k++ {
//		for v := 0; v < k; v++ {
//			if nums[k] > nums[v] && ans[k] < ans[v]+1 {
//				ans[k] = ans[v] + 1
//				list[k] = v
//			}
//		}
//	}
//	maxLen := ans[0]
//	lastIdx := 0
//	for i := 1; i < l; i++ {
//		if ans[i] > maxLen {
//			maxLen = ans[i]
//			lastIdx = i
//		}
//	}
//	path := make([]int, maxLen)
//	for lastIdx != -1 {
//		path = append([]int{nums[lastIdx]}, path...)
//		lastIdx = list[lastIdx]
//	}
//	fmt.Println(path)
//	return maxLen
//}
//func main() {
//	nums := []int{10, 9, 2, 5, 3, 7, 101, 18}
//	fmt.Println(lengthOfLIS(nums))
//}
