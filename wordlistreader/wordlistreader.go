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
	iterStarted bool
	itermu      sync.Mutex
}

func (wlr *WordListReader) ReadLine() (string, bool) {
	wlr.readlinemu.Lock()
	defer wlr.readlinemu.Unlock()
	end := wlr.scanner.Scan()
	return wlr.scanner.Text(), end

}
func (wlr *WordListReader) Close() {
	wlr.file.Close()
}

func (wlr *WordListReader) startIter() {
	wlr.iterStarted = true
	wlr.iterChannel = make(chan string)
	go func() {
		cont := true
		str := ""
		for cont {
			str, cont = wlr.ReadLine()
			wlr.iterChannel <- str
		}
		close(wlr.iterChannel)
	}()
}

func (wlr *WordListReader) Iter() chan string {
	wlr.itermu.Lock()
	if !wlr.iterStarted {
		wlr.startIter()
	}
	wlr.itermu.Unlock()
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
