package wordlistreader

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type WordListReader struct {
	readlinemu  sync.Mutex
	filename    string
	file        *os.File
	scanner     *bufio.Scanner
	iterChannel chan string
	itermu      sync.Once
	testBool    bool
}

func (wlr *WordListReader) readLine() (string, bool) {
	//wlr.readlinemu.Lock()
	//defer wlr.readlinemu.Unlock()
	return wlr.scanner.Text(), wlr.scanner.Scan()

}
func (wlr *WordListReader) Close() {
	wlr.file.Close()
}

func (wlr *WordListReader) startIter() {
	if wlr.testBool {
		wlr.iterChannel = make(chan string, 3)
	} else {
		wlr.iterChannel = make(chan string)
	}

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

func (wlr *WordListReader) Iter() chan string {
	wlr.itermu.Do(wlr.startIter)
	return wlr.iterChannel

}

func MakeNewWordListReader(filename string) WordListReader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)

	return WordListReader{
		filename: filename,
		file:     f,
		scanner:  scanner,
	}

}
