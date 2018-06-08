package main

import (
	"fmt"
	"sync"
	"time"
)

type x struct {
	sync.Mutex
	val int
}

func main() {
	// waitGroup()
	mutexLock()
}

func waitGroup() {
	x := &x{val: 1}

	wait := &sync.WaitGroup{}

	wait.Add(1)

	go func() {
		x.val = 2
		wait.Done()
	}()

	wait.Wait()

	fmt.Println(x.val)
}

func mutexLock() {
	x := &x{val: 1}

	x.Lock()

	go func() {
		x.Lock()
		fmt.Println(x.val)
		x.Unlock()
	}()

	go func() {
		time.Sleep(1 * time.Second)
		x.val = 2
		x.Unlock()
	}()

	time.Sleep(2 * time.Second)
}
