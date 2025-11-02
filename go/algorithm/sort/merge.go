package sort

func mergeDown2Up(nums []int) {
	// merge [start...mid] and [mid+1...end]
	aux := make([]int, len(nums))
	merge := func(nums []int, start, mid, end int) {
		i, j := start, mid+1
		for k := start; k <= end; k++ {
			aux[k] = nums[k]
		}

		for k := start; k <= end; k++ {
			if i > mid {
				nums[k] = aux[j]
				j++
			} else if j > end {
				nums[k] = aux[i]
				i++
			} else if aux[i] < aux[j] {
				nums[k] = aux[i]
				i++
			} else {
				nums[k] = aux[j]
				j++
			}
		}
	}
	for sz := 1; sz < len(nums); sz *= 2 {
		for i := 0; i < len(nums)-sz; i += sz * 2 {
			merge(nums, i, i+sz-1, min(len(nums)-1, i+2*sz-1))
		}
	}
}
