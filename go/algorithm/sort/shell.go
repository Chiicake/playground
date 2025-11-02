package sort

func shellSort(nums []int) {
	k := 1
	for k < len(nums) {
		k = 3*k + 1
	}
	for k >= 1 {
		for i := k; i < len(nums); i++ {
			for j := i; j >= k && nums[j] < nums[j-k]; j -= k {
				nums[j], nums[j-k] = nums[j-k], nums[j]
			}
		}
		k /= 3
	}
}
