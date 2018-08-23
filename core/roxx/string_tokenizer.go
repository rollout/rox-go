package roxx

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type StringTokenizer struct {
	text            string
	currentPosition int
	newPosition     int
	maxPosition     int
	delimiters      string
	retDelims       bool
}

func NewStringTokenizer(text, delim string, returnDelims bool) *StringTokenizer {
	return &StringTokenizer{
		text:            text,
		delimiters:      delim,
		retDelims:       returnDelims,
		currentPosition: 0,
		newPosition:     -1,
		maxPosition:     len(text),
	}
}

func (st *StringTokenizer) skipDelimiters(startPos int) int {
	position := startPos
	for !st.retDelims && position < st.maxPosition {
		r, size := utf8.DecodeRuneInString(st.text[position:])
		if !st.isDelimiter(r) {
			break
		}
		position += size
	}

	return position

}

func (st *StringTokenizer) scanToken(startPos int) int {
	position := startPos
	for position < st.maxPosition {
		r, size := utf8.DecodeRuneInString(st.text[position:])
		if st.isDelimiter(r) {
			break
		}
		position += size
	}

	if st.retDelims && startPos == position {
		r, size := utf8.DecodeRuneInString(st.text[position:])
		if st.isDelimiter(r) {
			position += size
		}
	}

	return position
}

func (st *StringTokenizer) isDelimiter(r rune) bool {
	return strings.IndexRune(st.delimiters, r) >= 0
}

func (st *StringTokenizer) hasMoreTokens() bool {
	st.newPosition = st.skipDelimiters(st.currentPosition)
	return st.newPosition < st.maxPosition
}

func (st *StringTokenizer) nextToken(delim string) string {
	delimsChanged := false
	if delim != "" {
		st.delimiters = delim
		delimsChanged = true
	}

	if st.newPosition >= 0 && !delimsChanged {
		st.currentPosition = st.newPosition
	} else {
		st.currentPosition = st.skipDelimiters(st.currentPosition)
	}

	st.newPosition = -1

	if st.currentPosition >= st.maxPosition {
		panic(fmt.Sprintf("StringTokenizer: currentPosition %d should be less than maxPosition %d", st.currentPosition, st.maxPosition))
	}

	start := st.currentPosition
	st.currentPosition = st.scanToken(st.currentPosition)
	return st.text[start:st.currentPosition]
}
