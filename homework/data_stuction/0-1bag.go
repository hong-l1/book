package data

import (
	"math"
)

func bag1(values []int, v []int, w int) (ans int) {
	dp := make([][]int, len(values)+1)
	for k := range dp {
		dp[k] = make([]int, len(v)+1)
		for v := range dp[k] {
			dp[k][v] = math.MinInt
		}
	}
	dp[0][0] = 0
	for i := 1; i <= len(values); i++ {
		for j := 1; j <= w; j++ {
			if j >= v[i-1] {
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-v[i-1]]+values[i-1])
			} else {
				dp[i][j] = dp[i-1][j]
			}
			ans = max(ans, dp[i][j])
		}
	}
	return ans
}
