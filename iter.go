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
// String functions are provided for convenience, but will obviously allocate.
package lineiter

import "bytes"

// LineIterator contains the state of the iterator.
// Copying a LineIterator will produce a new iterator with the same state.
type LineIterator struct {
	data    []byte
	start   int
	end     int
	newline int
}

func MakeLineIteratorString(data string) LineIterator {
	return MakeLineIterator([]byte(data))
}

func MakeLineIterator(data []byte) LineIterator {
	return LineIterator{
		data:    data,
		start:   -1,
		end:     -1,
		newline: -1,
	}
}

func MakeLineIteratorEndString(data string) LineIterator {
	return MakeLineIteratorEnd([]byte(data))
}

func MakeLineIteratorEnd(data []byte) LineIterator {
	return LineIterator{
		data:    data,
		start:   len(data) + 1,
		end:     len(data) + 1,
		newline: len(data) + 1,
	}
}

func (it *LineIterator) Next() bool {
	if it.newline > len(it.data) {
		return false
	}

	it.start = it.newline + 1
	if it.start > len(it.data) {
		it.end = it.start
		it.newline = it.start
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
	return true
}

func (it *LineIterator) Offset() int {
	return it.start
}

func (it *LineIterator) NextOffset() int {
	return it.newline + 1
}

// Returns current line as byte slice.
func (it *LineIterator) Bytes() []byte {
	return it.data[it.start:it.end]
}

// Returns current line as string.
func (it *LineIterator) Text() string {
	return string(it.Bytes())
}
