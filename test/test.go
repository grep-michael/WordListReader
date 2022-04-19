package test

import (
	"fmt"
	"sync"

	"github.com/michaelknudsen/WordListReader/wordlistreader"
)

/*

	Methods will sometimes fail and i blieve thats because the
	thread that closes a channel (i.e one or two) closes the channel as
	the 3rd thread tries to read from it.

*/

func Itertest() {
	var wg sync.WaitGroup
	wlr := wordlistreader.MakeNewWordListReader("rockyou.txt")
	defer wlr.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for word := range wlr.Iter() {
			fmt.Println("one", word)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for word := range wlr.Iter() {
			fmt.Println("two", word)
		}
	}()
	wg.Wait()

}

/*UNSUPPORTED
func ReadLineTest() {
	var done = make(chan int, 2)


		individual threads access readline

	one := make(chan string, 1)
	two := make(chan string, 1)

	var wg sync.WaitGroup
	wlr := wordlistreader.MakeNewWordListReader("./rockyou.txt")
	defer wlr.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cont := true
		str := ""
		for cont {
			str, cont = wlr.ReadLine()
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
			str, cont = wlr.ReadLine()
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

}*/
