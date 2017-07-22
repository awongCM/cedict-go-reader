// Credits - https://blog.gopheracademy.com/advent-2014/parsers-lexers/

package cedict-go-reader

//Go libraries to import
import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strings"
	"io"
)

//Data structures defined

//Our starting AST construct to parse
type ChineseCEDictReader struct {
	*bufio.Scanner
	TokenType int
	entry *Entry
}

//Define Tokens
const (
	ENTRY = iota
	COMMENTENTRY    // #

	SQUAREBRACKETS // [] 
	FORWARDSLASH // \
)

type Entry struct {
	Simplified string	
	Traditional string	
	Pinyin string	
	PinyinWithTones string	
	PinyinNoTones string	
	Definitions [] string
}

func main() {
	//Todo

	// const input = "1,2,3,4,"
	// scanner := bufio.NewScanner(strings.NewReader(input))

	// onComma := func(data []byte, atEOF bool) (advance int, token []byte, err error){
	// 	for i := 0; i < len(data); i++ {
	// 		if data[i] == ',' {
	// 			return i + 1, data[:1], nil
	// 		}
	// 	}

	// 	return 0, data, bufio.ErrFinalToken
	// }
	// scanner.Split(onComma)
	// for scanner.Scan() {
	// 	fmt.Printf("%q ", scanner.Text())
	// }
	// if err := scanner.Err(); err != nil {
	// 	fmt.Fprintln(os.Stderr, "reading input:", err)
	// 	log.Fatal(err)
	// }
}
