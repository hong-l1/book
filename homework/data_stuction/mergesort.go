package data

func MergeSort(nums []int) []int {
	if len(nums) <= 1 {
		return nums
	}
	left := MergeSort(nums[:len(nums)/2])
	right := MergeSort(nums[len(nums)/2:])
	return Merge(left, right)
}
func Merge(nums1 []int, nums2 []int) []int {
	left := 0
	right := 0
	ans := make([]int, 0, len(nums1)+len(nums2))
	for left < len(nums1) && right < len(nums2) {
		if (left < len(nums1) && nums1[left] <= nums2[right]) || right == len(nums2) {
			ans = append(ans, nums1[left])
			left++
		} else {
			ans = append(ans, nums2[right])
			right++
		}
	}
	return ans
}
