package test

import (
	"fmt"
	"sync"

	"github.com/michaelknudsen/WordListReader/wordlistreader"
)

var done = make(chan int, 2)

func checkdone() bool {
	return len(done) >= 1
}

func main() {
	/*
		testing
	*/
	one := make(chan string)
	two := make(chan string)

	var wg sync.WaitGroup
	wlr := wordlistreader.MakeNewWordListReader("./rockyou.txt")
	defer wlr.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cont := true
		str := ""
		for cont {
			str, cont = wlr.Readline()
			one <- str
		}
		fmt.Println("thread1 : return one")
		close(one)
		done <- 1
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		cont := true
		str := ""
		for cont {
			str, cont = wlr.Readline()
			two <- str
		}
		fmt.Println("thread2 : return two")
		close(two)
		done <- 2
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var oneS string
		var twoS string
		for {
			oneS = "default"
			twoS = "default"
			select {
			case <-one:
				oneS = <-one
			case <-two:
				twoS = <-two
			case <-done:
				if len(done) >= 2 {
					return
				}
			}
			if twoS == oneS {
				fmt.Println("Same")
				return
			} else {
				fmt.Println(oneS, twoS)
			}
		}

	}()
	wg.Wait()
}
