package main

import "fmt"

func longestCommonSubsequence(text1 string, text2 string) int {
	m := len(text1)
	n := len(text2)
	ans := make([][]int, m+1)
	path := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		ans[i] = make([]int, n+1)
		path[i] = make([]int, n+1)
	}
	for k := 1; k < m+1; k++ {
		for v := 1; v < n+1; v++ {
			if text1[k-1] != text2[v-1] {
				if ans[k][v-1] > ans[k-1][v] {
					ans[k][v] = ans[k][v-1]
					path[k][v] = 1
				} else {
					ans[k][v] = ans[k-1][v]
					path[k][v] = 2
				}
			} else {
				path[k][v] = 3
				ans[k][v] = ans[k-1][v-1] + 1
			}
		}
	}
	list := make([]byte, 0)
	row, col := m, n
	for row > 0 && col > 0 {
		if path[row][col] == 3 {
			list = append([]byte{text1[row-1]}, list...) // 注意这里是 row-1
			row--
			col--
		} else if path[row][col] == 2 {
			row--
		} else {
			col--
		}
	}
	fmt.Println(string(list))
	return ans[m][n]
}
func main() {
	text1 := "abcde"
	text2 := "ace"
	fmt.Println(longestCommonSubsequence(text1, text2))
}
