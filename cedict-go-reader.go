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
	bufio_s := bufio.NewScanner(r)
	e := &ChineseCEDictReader{
		Scanner: bufio_s,
	}

	// var line_input [] string

	// for e.Scanner.Scan() {
	// 	line_input = append(line_input, e.Scanner.Text())
	// }

	// return line_input

	splitFunc := func (data [] byte, atEOF bool) (advance int, token []byte, err error) {
		
		if len(data) {
			return			
		}

		if data[0] == "#" {
			
			e.TokenType = COMMENTENTRY
			e.lineInput = data
		}
		else {

			e.TokenType = DICTENTRY
			e.lineInput = data
		}
		return
	}

	bufio_s.Split(splitFunc)

	return e

}

func main() {
	input := "世界 世界 [shi4 jie4] /world/CL:個|个[ge4]/"

	r := io.Reader(strings.NewReader(input))
	
	startingEntry := NewEntry(r)

	fmt.Println("startingEntry: %s", startingEntry.lineInput)

	// lineInputs := NewEntry(r)

	// for index, element := range lineInputs {
	// 	fmt.Println("Each line input: [%d] index is [%s]", index, element)
	// }
}
