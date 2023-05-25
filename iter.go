// Package lineiter provides an allocation-free, zero-copy line iterator.
//
// Iterators have two special states: before-the-beginning and past-the-end.
// A new iterator starts in the before-the-beginning state and calls to Next()
// will move the iterator to the next line. Next() will return false when the
// iterator reaches the past-the-end state. This allows for the usual iteration
// pattern:
//
//	it := MakeLineIterator(...)
//	for it.Next() {
//	    var line []byte = it.Bytes()
//	    fmt.Printf("%s\n", line)
//
//	    // More convenient:
//	    // fmt.Printf("%s\n", it.Text())
//	}
//
// Semantically, this is equivalent to the following code, except that a
// carriage return before the newline is trimmed by the line iterator:
//
//	for _, line := range strings.Split(..., "\n") {
//	    fmt.Printf("%s\n", line)
//	}
//
// In particular, the empty string contains one line and "a\n" contains two
// lines ("a" and "").
//
// In the special states, the line data accessors return default values as makes
// sense (e.g. Bytes() and Text() will return an empty slice or string).
//
// String functions are provided for convenience, but will obviously allocate.
package lineiter

import "bytes"

// LineIterator contains the state of the iterator.
// Copying a LineIterator will produce a new iterator with the same state.
type LineIterator struct {
	data []byte

	start   int
	end     int
	newline int

	lineNumber int
	lineCount  int
}

func MakeLineIteratorString(data string) LineIterator {
	return MakeLineIterator([]byte(data))
}

func MakeLineIterator(data []byte) LineIterator {
	return LineIterator{
		data:       data,
		start:      -1,
		end:        -1,
		newline:    -1,
		lineNumber: 0,
	}
}

func MakeLineIteratorEndString(data string) LineIterator {
	return MakeLineIteratorEnd([]byte(data))
}

func MakeLineIteratorEnd(data []byte) LineIterator {
	return LineIterator{
		data:       data,
		start:      len(data) + 1,
		end:        len(data) + 1,
		newline:    len(data) + 1,
		lineNumber: -1,
	}
}

// Moves the iterator into the before-the-beginning state.
func (it *LineIterator) SeekStart() {
	it.start = -1
	it.end = -1
	it.newline = -1
	it.lineNumber = 0
}

// Moves the iterator into the past-the-end state.
func (it *LineIterator) SeekEnd() {
	it.start = len(it.data) + 1
	it.end = len(it.data) + 1
	it.newline = len(it.data) + 1
	it.lineNumber = -1
}

func (it *LineIterator) Next() bool {
	if it.newline > len(it.data) {
		return false
	}

	it.start = it.newline + 1
	if it.start > len(it.data) {
		it.end = it.start
		it.newline = it.start
		if it.lineNumber > 0 {
			it.lineCount = it.lineNumber
		}
		it.lineNumber = -1
		return false
	}

	index := bytes.IndexByte(it.data[it.start:], '\n')
	if index == -1 {
		it.newline = len(it.data)
	} else {
		it.newline = it.start + index
	}

	it.end = it.newline
	if it.end != 0 && it.data[it.end-1] == '\r' {
		it.end--
	}

	it.lineNumber++
	return true
}

func (it *LineIterator) Previous() bool {
	if it.start < 0 {
		return false
	}

	it.newline = it.start - 1
	if it.newline < 0 {
		it.start = -1
		it.end = -1
		if it.lineNumber < 0 {
			it.lineCount = -(it.lineNumber + 1)
		}
		it.lineNumber = 0
		return false
	}

	index := bytes.LastIndexByte(it.data[:it.newline], '\n')
	if index == -1 {
		it.start = 0
	} else {
		it.start = index + 1
	}

	it.end = it.newline
	if it.end != 0 && it.data[it.end-1] == '\r' {
		it.end--
	}

	it.lineNumber--
	return true
}

// Offset of current line in underlying data slice.
// Returns -1 in before-the-beginning state.
// Returns FullLength()+1 in past-the-end state.
func (it *LineIterator) Offset() int {
	return it.start
}

// Returns current line as byte slice.
// Returns empty slice if in special state.
func (it *LineIterator) Bytes() []byte {
	if it.start < 0 {
		return []byte{}
	}
	return it.data[it.start:it.end]
}

// Returns current line as string.
// Returns empty string if in special state.
func (it *LineIterator) Text() string {
	return string(it.Bytes())
}

// Returns length of current line.
// Returns 0 if in special state.
func (it *LineIterator) Length() int {
	return it.end - it.start
}

// Returns true if iterator points to valid line and false if iterator is in special state.
func (it *LineIterator) Valid() bool {
	return it.start >= 0 && it.start <= len(it.data)
}

// Returns underlying data as byte slice.
func (it *LineIterator) FullBytes() []byte {
	return it.data
}

// Returns underlying data as string.
func (it *LineIterator) FullText() string {
	return string(it.data)
}

// Returns length of underlying data.
func (it *LineIterator) FullLength() int {
	return len(it.data)
}

// Returns line number relative to the last special state (beginning or end).
// If relative to the beginning, the result is positive and the first line has line number 1.
// If relative to the end, the result is negative and the last line has line number -1.
// Returns 0 if in special state.
func (it *LineIterator) RelativeLineNumber() int {
	if it.lineNumber < 0 {
		return it.lineNumber + 1
	} else {
		return it.lineNumber
	}
}

// Returns absolute line number.
// The first line has line number 1.
// Returns 0 or LineCount()+1 if in special state.
// May have to read the entire data slice to count the number of lines,
// which can be expensive for large data sets.
func (it *LineIterator) LineNumber() int {
	if it.lineNumber < 0 {
		return it.LineCount() + it.lineNumber + 2
	}
	return it.lineNumber
}

// Returns the total number of lines, if known, or 0, otherwise.
func (it *LineIterator) OptionalLineCount() int {
	return it.lineCount
}

// Returns the total number of lines.
// May have to read the entire data slice to count the number of lines,
// which can be expensive for large data sets.
// The value is cached in the line iterator.
func (it *LineIterator) LineCount() int {
	if it.lineCount == 0 {
		it.lineCount = bytes.Count(it.data, []byte{'\n'}) + 1
	}
	return it.lineCount
}
