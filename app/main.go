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
	fmt.Println("match", ok)
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
	if pattern[0] == '^' {
		v := match(pattern[1:], line)
		if !v {
			return false, nil
		}
	}
	for i := 0; i < len(line); i++ {

		isMatch := match(pattern, line[i:])
		fmt.Println("Pattern found ? ", isMatch)
		if isMatch {
			return true, nil
		}
	}

	return false, nil

}

func match(pattern string, text []byte) bool {

	//This is our base case, if we reach the point were we consumed all the pattern
	// without failure, it means we've found a match
	if len(pattern) == 0 {
		return true
	}

	println("current pattern", pattern)
	println("current text", string(text))

	if len(pattern) > 1 && pattern[1] == '?' {

		return matchOptionalOperator(text, pattern)
	}

	if len(pattern) > 1 && pattern[1] == '+' {
		return matchPlusQuantifier(text, pattern)

	}

	if len(text) > 0 && pattern[0] == '.' {
		return matchesAnyCharacter(text, pattern)
	}

	//If the current char in the pattern it's a slash, it means is a character class
	// and we need to handle it in a special way.
	if pattern[0] == '\\' && len(text) > 0 {
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

	if pattern[0] == '^' {

		isAMatch, newText, newPattern := matchStartAnchor(text, pattern[1:])
		if !isAMatch {
			return false
		} else {
			return match(newPattern, newText)
		}
	}

	if len(pattern) == 1 && pattern[0] == '$' {
		return len(text) == 0
	}

	//If we reach this point, it means the pattern hasn't been fully matched yet
	// if it doesn't match or there is no more text, its a no match and we return fail
	if len(text) > 0 && pattern[0] == text[0] {
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

func matchStartAnchor(text []byte, pattern string) (bool, []byte, string) {

	i := 0

	for i < len(pattern) {
		if pattern[i] == '$' {
			return true, text[i:], pattern[i:]
		}
		if text[i] != pattern[i] {
			return false, nil, ""
		}
		i++
	}

	return true, text[i:], pattern[i:]
}

func matchPlusQuantifier(text []byte, pattern string) bool {
	//Here we need to answer the question, are you the letter x, or any of the letters
	// after you are equal to x
	letterToMatch := pattern[0]
	fmt.Println("letter to match", string(letterToMatch))

	if letterToMatch == '.' {
		i := 1
		for i < len(text) {
			if text[i] == pattern[2] {
				break
			}
			i++
		}
		return match(pattern[2:], text[i:])
	}

	//ca+r  caaaaaarla
	if text[0] != letterToMatch {
		fmt.Println("ups")
		return false
	}

	//Start from 1 because we already read the current char
	i := 1

	//Need to be careful of not consuming the character in the pattern right after the +
	// even if its the same

	for i < len(text) {

		if text[i] != letterToMatch {
			break
		}
		i++
	}

	if pattern[2] == letterToMatch {
		return match(pattern[2:], text[i-1:])
	}
	return match(pattern[2:], text[i:])
}

func matchOptionalOperator(text []byte, pattern string) bool {
	charToIgnore := pattern[0]

	// the current character is different from the optional one
	if len(text) >= len(pattern)-1 {
		fmt.Println("Text has the same characters that pattern, check if the character matches the optional")
		if text[0] != charToIgnore {
			return false
		} else {
			return match(pattern[2:], text[1:])
		}
	}

	//Reaching here means the optional character does not exist, so we move in the pattern
	// but did not read the text
	return match(pattern[2:], text)
}

func matchesAnyCharacter(text []byte, pattern string) bool {
	fmt.Println("breaks here")
	return match(pattern[1:], text[1:])
}
