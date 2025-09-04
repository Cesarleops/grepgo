package main

import (
	"fmt"
	"io"
	"os"
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

	for i := 0; i < len(line); i++ {
		isMatch := match(pattern, line[i:])

		if isMatch {
			return true, nil
		}
	}

	return false, nil

}

func match(pattern string, text []byte) bool {

	println("pattern len", len(pattern))
	println("text len", len(text))

	if len(pattern) == 0 {
		fmt.Println("nice")
		return true
	}

	if len(pattern) > 0 && len(text) == 0 {
		fmt.Println("fail")
		return false
	}

	println("current pattern", pattern)
	println("current text", string(text))
	println("current char", string(text[0]))

	if pattern[0] == '\\' {

		if pattern[1] == 'd' {
			isAMatch := isDigit(text[0])
			if !isAMatch {
				return false
			} else {
				return match(pattern[2:], text[1:])
			}
		}

		if pattern[1] == 'w' {
			isAMatch := isAlpha(text[0]) || isDigit(text[0])
			if !isAMatch {
				return false
			} else {
				return match(pattern[2:], text[1:])
			}
		}

	}

	if pattern[0] == '[' {
		if pattern[1] == '^' {
			isAMatch, newPattern := matchNegativeGroupCharacter(text[0], pattern[2:])
			if !isAMatch {
				return false
			} else {
				return match(newPattern, text[1:])
			}
		} else {
			isAMatch, newPattern := matchPositiveGroupCharacter(text[0], pattern[1:])
			if !isAMatch {
				return false
			} else {
				return match(newPattern, text[1:])
			}
		}

	}
	fmt.Println("normal check")
	if pattern[0] == text[0] {
		fmt.Println("normal match")
		return match(pattern[1:], text[1:])

	} else {
		return false
	}
}

func isDigit(c byte) bool {

	//since line[i] returns the decimal of the byte value, we can compare directly
	// to '0' and '9'
	if c >= '0' && c <= '9' {
		return true
	}
	return false

}

func isAlpha(c byte) bool {
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_') {
		return true
	}
	return false
}

func matchPositiveGroupCharacter(c byte, pattern string) (bool, string) {

	fmt.Println("pattern for positive group", pattern)
	//Trick here is checking if there is one char equal to any from the ones in the pattern.
	// The pattern already has trimmed the first [ need to iterate until the next ]
	i := 0

	//Here, we need to advance the i until the group ends
	// event if the match happens early.
	matchAny := false
	for i < len(pattern) {
		if pattern[i] == ']' {
			break
		}
		if c == pattern[i] && !matchAny {
			matchAny = true
		}
		i++
	}

	if matchAny {
		return true, pattern[i+1:]
	} else {
		fmt.Println("didn't match positive group")
		return false, ""
	}

}

func matchNegativeGroupCharacter(c byte, pattern string) (bool, string) {
	//abc]m
	fmt.Println("pattern for negative group", pattern)
	//Trick here is checking if there is one char different from the ones in the pattern.
	// The pattern already has trimmed the first [ need to iterate until the next ]
	i := 0
	for i < len(pattern) {
		if pattern[i] == ']' {
			break
		}
		if c == pattern[i] {
			return false, ""
		}
		i++
	}

	return true, pattern[i+1:]

}
