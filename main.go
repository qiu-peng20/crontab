package main

import (
	"fmt"
	"time"
)

func main()  {
	num := 0
	go func() {
		for  {
			fmt.Print(22222)
			num++
		}
	}()
	time.Sleep(time.Second)
	fmt.Print(num)
}

