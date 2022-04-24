package wordlistreader

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type WordListReader struct {
	file        *os.File
	scanner     *bufio.Scanner
	iterChannel chan string
	itermu      sync.Once
}

func (wlr *WordListReader) readLine() (string, bool) {

	return wlr.scanner.Text(), wlr.scanner.Scan()

}
func (wlr *WordListReader) Close() {
	wlr.file.Close()
}

func (wlr *WordListReader) startIter() {
	go func() {
		cont := true
		str := ""
		for cont {
			str, cont = wlr.readLine()
			wlr.iterChannel <- str
		}
		close(wlr.iterChannel)
	}()
	return
}

func (wlr *WordListReader) Iter() <-chan string {
	wlr.itermu.Do(wlr.startIter)
	return wlr.iterChannel
}

func MakeBufferedWordListReader(filename string, buffSize int) *WordListReader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)

	return &WordListReader{
		file:        f,
		scanner:     scanner,
		iterChannel: make(chan string, buffSize),
	}
}

func MakeUnbufferedWordListReader(filename string) *WordListReader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)

	return &WordListReader{
		file:        f,
		scanner:     scanner,
		iterChannel: make(chan string),
	}
}
