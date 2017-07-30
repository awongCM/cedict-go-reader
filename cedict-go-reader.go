// Credits - https://blog.gopheracademy.com/advent-2014/parsers-lexers/

package main

//Go libraries to import
import (
	"bufio"
	"fmt"
	"strings"
	"io"
	"os"
)

//Data structures defined

//Our starting AST construct to parse
type ChineseCEDictReader struct {
	*bufio.Scanner
	TokenType int
	entry *Entry
	lineInput string 
}

//Define line tokens
const (
	DICT_ENTRY = iota
	COMMENT_ENTRY    // #
	ERR_ENTRY  //NIL
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
			e.TokenType = COMMENT_ENTRY
			advance, token, err = processCommentEntry(data, atEOF)
			e.lineInput = string(token)
			
		} else {
			e.TokenType = DICT_ENTRY
			advance, token, err = processDictEntry(data, atEOF)
			e.lineInput = string(token)
		}
		return
	}

	s.Split(splitFunc)
	return e

}

func processCommentEntry(data []byte, atEOF bool) (int, []byte, error){
	var tokens []byte

	for i, b := range data {
		if b =='\n' || (atEOF && i == len(data) - 1){
			return i + 1, tokens, nil
		} else {
			tokens = append(tokens, b)
		}
	}

	if atEOF {
		return len(data), tokens, nil
	}

	return 0, nil, nil
}

func processDictEntry(data []byte, atEOF bool) (int, []byte, error) {
	var tokens []byte

	for i, b:= range data {
		if b == '\n' {
			return i + 1, tokens, nil
		} else {
			tokens = append(tokens, b)
		}
	}

	if atEOF {
		return len(data), tokens, nil
	}

	return 0, nil, nil
}



func main() {
	input := "# Comment \n 世界 世界 [shi4 jie4] /world/CL:個|个[ge4]/ hello world \n"

	r := io.Reader(strings.NewReader(input))
	
	startingEntry := NewEntry(r)


	for startingEntry.Scan() {

		if startingEntry.TokenType == DICT_ENTRY  {
			fmt.Println("dict entry found", startingEntry.lineInput)
		} else if startingEntry.TokenType == COMMENT_ENTRY{
			fmt.Println("comment entry found", startingEntry.lineInput)
		}

		if err := startingEntry.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		
	}	

}
