package wordlistreader

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type wordListReader struct {
	file        *os.File
	scanner     *bufio.Scanner
	iterChannel chan string
	itermu      sync.Once
}

func (wlr *wordListReader) readLine() (string, bool) {

	return wlr.scanner.Text(), wlr.scanner.Scan()

}
func (wlr *wordListReader) Close() {
	wlr.file.Close()
}

func (wlr *wordListReader) startIter() {
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

func (wlr *wordListReader) Iter() <-chan string {
	wlr.itermu.Do(wlr.startIter)
	return wlr.iterChannel
}

func MakeBufferedWordListReader(filename string, buffSize int) *wordListReader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)

	return &wordListReader{
		file:        f,
		scanner:     scanner,
		iterChannel: make(chan string, buffSize),
	}
}

func MakeUnbufferedWordListReader(filename string) *wordListReader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)

	return &wordListReader{
		file:        f,
		scanner:     scanner,
		iterChannel: make(chan string),
	}
}
