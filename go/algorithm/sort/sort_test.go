package sort

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShellSort(t *testing.T) {
	nums := []int{5, 4, 3, 2, 1}
	shellSort(nums)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, nums)
}

func TestMergeDown2Up(t *testing.T) {
	nums := []int{5, 4, 3, 2, 1}
	mergeDown2Up(nums)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, nums)
}
