package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {

	if pattern == "\\w" {
		validAlpha := isAlpha(line)
		validDigit := isDigit(line)

		return validAlpha || validDigit, nil
	}

	if pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
		if pattern[1] == '^' {
			return matchNegativeGroupCharacter(line), nil
		}
		return matchPositiveCharacter(line, pattern), nil
	}

	if pattern == "\\d" {

		ok := isDigit(line)
		return ok, nil
	}

	if utf8.RuneCountInString(pattern) != 1 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool

	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	ok = bytes.ContainsAny(line, pattern)

	return ok, nil
}

func isDigit(line []byte) bool {

	//since line[i] returns the decimal of the byte value, we can compare directly
	// to '0' and '9'
	for i := range line {
		if line[i] >= '0' && line[i] <= '9' {
			return true
		}
	}
	return false

}

func isAlpha(line []byte) bool {
	for i := range line {
		if (line[i] >= 'a' && line[i] <= 'z') || (line[i] >= 'A' && line[i] <= 'Z') || (line[i] == '_') {
			return true
		}
	}
	return false
}

func matchPositiveCharacter(line []byte, pattern string) bool {

	pattern = pattern[1 : len(pattern)-1]

	chars := strings.Split(pattern, "")

	set := make(map[byte]struct{})
	for _, v := range chars {
		set[v[0]] = struct{}{}
	}

	for i := range line {

		if _, ok := set[line[i]]; ok {
			return true
		}
	}

	return false
}

func matchNegativeGroupCharacter(line []byte, pattern string) bool {

	//Trick here is checking if there is one char different from the ones in the pattern.
	pattern = pattern[2 : len(pattern)-1]

	chars := strings.Split(pattern, "")

	fmt.Println("chars", chars)
	set := make(map[byte]struct{})

	for _, v := range chars {
		set[v[0]] = struct{}{}
	}

	for i := range line {

		if _, ok := set[line[i]]; !ok {
			return true
		}
	}

	return false
}
