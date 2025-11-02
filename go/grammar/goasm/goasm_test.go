package goasm

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	result := add(5, 3)
	fmt.Println(result)
	if result != 8 {
		t.Errorf("add(5, 3) = %d; want 8", result)
	}
}
