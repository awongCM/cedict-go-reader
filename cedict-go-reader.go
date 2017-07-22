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
	s *bufio.Scanner
	TokenType int
	entry *Entry
}

//Define Tokens
const (
	ENTRY = iota
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

func NewEntry(r io.Reader) []string {
	bufio_s := bufio.NewScanner(r)
	e := &ChineseCEDictReader{
		s: bufio_s,
	}

	var line_input [] string

	for e.s.Scan() {
		line_input = append(line_input, e.s.Text())
	}

	return line_input
}

func main() {
	input := "世界 世界 [shi4 jie4] /world/CL:個|个[ge4]/ \n 你好 你好 [ni3 hao3] /Hello!/Hi!/How are you?/\n"

	r := io.Reader(strings.NewReader(input))
	
	// startingEntry := NewEntry(r)

	// fmt.Println("startingEntry: %s", startingEntry)

	lineInputs := NewEntry(r)

	for index, element := range lineInputs {
		fmt.Println("Each line input: [%d] index is [%s]", index, element)
	}
}
