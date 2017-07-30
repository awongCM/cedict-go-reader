// Credits goes - https://blog.gopheracademy.com/advent-2014/parsers-lexers/

package main

//Go libraries to import
import (
	"bufio"
	"fmt"
	"strings"
	"io"
	"os"
	"regexp"
	"bytes"
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

func extractPinyinWithTones(p string) string {
	pv := strings.Replace(p, "u:", "ü", -1)
	py := strings.Split(pv, " ")

	var output bytes.Buffer
	for _, pySyllable := range py {
		pyNoTone, tone := extractTone(pySyllable)
		pyWithTone, err := replaceWithToneMark(pyNoTone, tone)
		if err != nil {
			output.WriteString(pySyllable)
		} else {
			output.WriteString(pyWithTone)
		}
	}
	return output.String()

}

func extractPinyinWithoutTones(p string) string {
	pv := strings.Replace(p, "u:", "v", -1)
	py := strings.Split(pv, " ")

	var output bytes.Buffer
	for _, pySyllable := range py {
		pyNoTone, _ := extractTone(pySyllable)
		output.WriteString(pyNoTone)
	}
	return output.String()
}

func extractTone(p string) (string ,int){
	tone := int(p[len(p)-1]) - 48

	if tone > 5 || tone < 0 {
		return p, 0
	}

	return p[0: len(p)-1], tone
}

func replaceWithToneMark(s string, tone int) (string, error) {
	lookup, err :=toneLookUpTable(tone)

	if err != nil {
		return "", err
	}

	if strings.Contains(s, "a") {
		return strings.Replace(s, "a", lookup["a"], -1), nil
	}

	if strings.Contains(s, "e") {
		return strings.Replace(s, "e", lookup["e"], -1), nil
	}

	if strings.Contains(s, "ou") {
		return strings.Replace(s, "o", lookup["o"], -1), nil
	}
	index := strings.LastIndexAny(s, "iüou")

	if index != -1 {
		var output bytes.Buffer
		for i, runeVal := range s {
			if i == index {
				output.WriteString(lookup[string(runeVal)])
			} else {
				output.WriteString(string(runeVal))
			}
		}
		return output.String(), nil
	}

	return "", fmt.Errorf("No tone match")

}

func toneLookUpTable(tone int) (map[string]string, error) {

	if tone < 0 || tone > 5 {
		return nil, fmt.Errorf("Tried to create tone lookup table with tone %i", tone)
	}

	lookupTable := map[string][]string{

	}

	toneLookup := make(map[string]string)

	for vowel, toneRunes:= range lookupTable {
		toneLookup[vowel] = toneRunes[tone]
	}

	return toneLookup, nil
}

var regExEntry = regexp.MustCompile(`(?P<trad>\S*?) (?P<simp>\S*?) \[(?P<pinyin>.+)\] \/(?P<defs>.+)\/`)

func parseEntry(s string) (*Entry, error) {
	match := regExEntry.FindStringSubmatch(s)

	if match == nil {
		return nil, fmt.Errorf("Format Error for entry: %v", s)
	}

	e := Entry{}

	for i, repattern := range regExEntry.SubexpNames()  {
		if i == 0 || repattern == "" {
			continue
		}

		switch repattern {
		case "trad":
			e.Simplified = match[i]
		case "simp":
			e.Traditional = match[i]
		case "pinyin":
			e.Pinyin = match[i]
		case "defs":
			e.Definitions = strings.Split(match[i], "/")
		}
	}

	e.PinyinWithTones = extractPinyinWithTones(e.Pinyin)
	e.PinyinNoTones = extractPinyinWithoutTones(e.Pinyin)

	return &e, nil
}


func main() {
	input := "# Comment \n 世界 世界 [shi4 jie4] /world/CL:個|个[ge4]/ \n"

	r := io.Reader(strings.NewReader(input))
	
	startingEntry := NewEntry(r)

	for startingEntry.Scan() {

		if startingEntry.TokenType == DICT_ENTRY  {
			
			e, err := parseEntry(startingEntry.lineInput)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot parse entry", err)
			}
			startingEntry.entry = e
			fmt.Println("Dict Entry: ", startingEntry.entry.Simplified, startingEntry.entry.Definitions[0] )

		} else if startingEntry.TokenType == COMMENT_ENTRY{
			fmt.Println("comment entry found", startingEntry.lineInput)
		}

		if err := startingEntry.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		
	}	

}
