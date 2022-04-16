package wordlistreader

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type WordListReader struct {
	mu       sync.Mutex
	filename string
	file     *os.File
	scanner  *bufio.Scanner
}

func (wlr *WordListReader) Readline() (string, bool) {
	wlr.mu.Lock()
	defer wlr.mu.Unlock()
	end := wlr.scanner.Scan()
	return wlr.scanner.Text(), end

}
func (wlr *WordListReader) Close() {
	wlr.file.Close()
}

func MakeNewWordListReader(filename string) WordListReader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
	}
	scanner := bufio.NewScanner(f)

	return WordListReader{
		filename: filename,
		file:     f,
		scanner:  scanner,
	}

}
