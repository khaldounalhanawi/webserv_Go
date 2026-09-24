package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	t := time.Now()
	start := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			println("order received")
			<-start
			time.Sleep(time.Second * 3)
			fmt.Printf("%d: Printing at %d,%d,%d\n",
					i,
					time.Now().Hour(),
					time.Now().Minute(),
					time.Now().Second(),
				)
			wg.Done()
		} ()
		time.Sleep(time.Second * 1)
	}
	println("All orders has been created!")

	close(start)
	wg.Wait()
	println("Program is done")
	fmt.Println(time.Since(t))
}