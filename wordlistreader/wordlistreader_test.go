package wordlistreader

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*

	Methods will sometimes fail and i blieve thats because the
	thread that closes a channel (i.e one or two) closes the channel as
	the 3rd thread tries to read from it.

*/

func TestSpeedOneThread(t *testing.T) {
	var wg sync.WaitGroup
	wlr := MakeNewWordListReader("../rockyou.txt")
	wg.Add(1)
	start := time.Now()
	go func() {
		defer wg.Done()
		for range wlr.Iter() { //simply read
		}
	}()
	duration := time.Since(start)
	fmt.Printf("Duation of unbuffered One thread: %v\n", duration)
	wlr.Close()

	wlr = MakeNewWordListReader("../rockyou.txt")
	wlr.testBool = true
	wg.Add(1)
	start = time.Now()
	go func() {
		defer wg.Done()
		for range wlr.Iter() { //simply read
		}
	}()
	duration = time.Since(start)
	fmt.Printf("Duation of unbuffered One thread: %v\n", duration)
	wlr.Close()

}

func TestSpeedTwoThread(t *testing.T) {
	var wg sync.WaitGroup
	wlr := MakeNewWordListReader("../rockyou.txt")
	wg.Add(1)
	start := time.Now()
	go func() {
		defer wg.Done()
		for range wlr.Iter() { //simply read
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range wlr.Iter() { //simply read
		}
	}()
	duration := time.Since(start)
	fmt.Printf("Duation of unbuffered Two thread: %v\n", duration)
	wlr.Close()

	wlr = MakeNewWordListReader("../rockyou.txt")
	wlr.testBool = true
	wg.Add(1)
	start = time.Now()
	go func() {
		defer wg.Done()
		for range wlr.Iter() { //simply read
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range wlr.Iter() { //simply read
		}
	}()
	duration = time.Since(start)
	fmt.Printf("Duation of unbuffered Two thread: %v\n", duration)
	wlr.Close()

}

func TestIter(t *testing.T) {

	var wg sync.WaitGroup
	wlr := MakeNewWordListReader("../rockyou.txt")
	defer wlr.Close()

	threadOneWords := []string{}
	threadTwoWords := []string{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for word := range wlr.Iter() {
			threadOneWords = append(threadOneWords, word)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for word := range wlr.Iter() {
			threadTwoWords = append(threadTwoWords, word)
		}
	}()
	wg.Wait()
	for _, v1 := range threadOneWords {
		for _, v2 := range threadTwoWords {
			if v1 == v2 {
				t.Error("Duplicate words in lists")
			}
		}
	}
}

func TestIterWithChannels(t *testing.T) {
	var done = make(chan int, 2)

	one := make(chan string, 1)
	two := make(chan string, 1)

	var wg sync.WaitGroup
	wlr := MakeNewWordListReader("../rockyou.txt")
	defer wlr.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for word := range wlr.Iter() {
			one <- word
		}
		close(one)
		done <- 1
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for word := range wlr.Iter() {
			two <- word
		}
		close(two)
		done <- 2
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var oneS string
		var twoS string
		doneCount := 0
		for doneCount < 1 {
			select {
			case x := <-one:
				oneS = x
			case y := <-two:
				twoS = y
			case <-done:
				doneCount += 1
			}
			if twoS == oneS && oneS != "" && twoS != "" {
				t.Errorf("Duplicate Words, \"%s\" : \"%s\"", oneS, twoS)
				return
			}

		}

	}()
	wg.Wait()

}
