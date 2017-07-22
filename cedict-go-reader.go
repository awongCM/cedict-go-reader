// Credits - https://blog.gopheracademy.com/advent-2014/parsers-lexers/

package main

//Go libraries to import
import (
	"bufio"
	"fmt"
	"strings"
	"io"
)

//Data structures defined

//Our starting AST construct to parse
type ChineseCEDictReader struct {
	*bufio.Scanner
	TokenType int
	entry *Entry
	lineInput string 
}

//Define Tokens
const (
	DICTENTRY = iota
	COMMENTENTRY    // #
	ERRORENTRY  //NIL
)

type Entry struct {
	Simplified string	
	Traditional string	
	Pinyin string	
	PinyinWithTones string	
	PinyinNoTones string	
	Definitions [] string
}

//Scanning inputs

func NewEntry(r io.Reader) *ChineseCEDictReader {
	s := bufio.NewScanner(r)
	e := &ChineseCEDictReader{
		Scanner: s,
	}

	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) == 0 {		
			return			
		}

		if data[0] == '#' {
			e.TokenType = COMMENTENTRY
			advance, token, err = bufio.ScanWords(data, atEOF)
			fmt.Println("debugger comment ", string(data))		
			e.lineInput = string(data)
		} else {
			e.TokenType = DICTENTRY
			advance, token, err = bufio.ScanWords(data, atEOF)
			fmt.Println("debugger dict", string(data))
			e.lineInput = string(data)
		}
		return

	}

	s.Split(splitFunc)
	return e

}



func main() {
	input := "世界 世界 [shi4 jie4] /world/CL:個|个[ge4]/ hello world"

	r := io.Reader(strings.NewReader(input))
	
	startingEntry := NewEntry(r)

	
	startingEntry.Scan();
	fmt.Println("lineinput", startingEntry.lineInput);

}
