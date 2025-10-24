package defer_and_recover

import (
	"fmt"
	"sync"
)

// defer is a stack, so the last defer will be executed first
// i can use recover to catch the panic in foo, and continue the execution of bar
// like try catch in java
func foo() {
	panic("foo")
}

func bar(i int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in bar:", r, "i:", i)
		}
	}()
	foo()
}

func deferAndRecover() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("goroutine start", i)
			bar(i)
			i++

		}()
	}
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in main:", r)
			}
		}()
		foo()
	}()
	wg.Wait()
}
