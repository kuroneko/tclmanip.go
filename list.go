package tclmanip

import (
	"unicode"
	"errors"
	"strings"
)

// TclList is a list encoded using Tcl semantics. 
type TclList string

var (
	// ErrIndexOutOfRange indicates that the index specified was out of range for the given list.
	ErrIndexOutOfRange = errors.New("List Index out of Range")
)

// needsQuoting() returns true if the TclList its called against requires quoting when included into
// a list
func (listin *TclList) needsQuoting() bool {
	if (len(*listin) == 0) {
		return true
	}
	for _, curRune := range *listin {
		switch (curRune) {
		case ' ':
			fallthrough
		case '\r':
			fallthrough
		case '\n':
			fallthrough
		case '\t':
			fallthrough
		case '\\':
			fallthrough
		case '\'':
			fallthrough
		case '"':
			return true
		}
	}
	return false
}

// Set assigns a new TclList using the given string
func Set(liststring string) *TclList {
	newList := TclList(liststring)
	return &newList
}

// String() returns the encoded Tcl List
func (listin *TclList) String() string {
	return string(*listin)
}

// Join joins a split TclList back together again.
func Join(listin []TclList) TclList {
	listOut := make([]string, len(listin))
	for i, s := range listin {
		listOut[i] = string(s)
	}
	return TclList(strings.Join(listOut, " "))
}


// Set implements the Tcl lset method
func (listin *TclList) Set(index int, sublist TclList) error {
	parts := listin.Split()
	if index >= len(parts) {
		return ErrIndexOutOfRange
	}
	if sublist.needsQuoting() {
		parts[index] = "{" + sublist + "}"
	} else {
		parts[index] = sublist
	}
	*listin = Join(parts)
	return nil
}

// Split breaks down a tcl list into a go slice.  You can use Join to put the
// pieces back together again.
//
// This only splits a single level.  Call against results to sub-split.
func (listin *TclList) Split() []TclList {
	var rs []TclList
	curlyDepth := 0

	curTok := []rune{}

	for _, curRune := range *listin {
		if (curlyDepth == 0) {
			if unicode.IsSpace(curRune) {
				// end the existing token
				if len(curTok) > 0 {
					rs = append(rs, TclList(curTok))
					curTok = []rune{}
				}
			} else {
				switch (curRune) {
				case '{':
					curlyDepth++
				default:
					curTok = append(curTok, curRune)
				}
			}
		} else {
			// depth > 0.  track curlies only.
			if (curRune == '{') {
				curlyDepth++;
			}
			if (curRune == '}') {
				curlyDepth--;
			}
			if (curlyDepth > 0) {
				curTok = append(curTok, curRune)
			}
		}
	}
	if len(curTok) > 0 {
		rs = append(rs, TclList(string(curTok)))
		curTok = []rune{}
	}
	return rs
}

// Index implements the Tcl command lindex.  It returns nil if the index is out of range.
func (listin *TclList) Index(index int) *TclList {
	parts := listin.Split()
	if (index >= len(parts)) {
		return nil
	}
	return &(parts[index])
}