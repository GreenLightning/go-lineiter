package lineiter

import (
	"bytes"
	"testing"
)

type LineIteratorTester struct {
	it LineIterator
	t  *testing.T
}

func MakeLineIteratorTester(t *testing.T, data string) LineIteratorTester {
	return LineIteratorTester{
		it: MakeLineIteratorString(data),
		t:  t,
	}
}

func MakeLineIteratorTesterEnd(t *testing.T, data string) LineIteratorTester {
	return LineIteratorTester{
		it: MakeLineIteratorEndString(data),
		t:  t,
	}
}

func (tester *LineIteratorTester) next(expected bool) {
	actual := tester.it.Next()
	if actual != expected {
		if expected {
			tester.t.Helper()
			tester.t.Fatal("expected line, but Next() returned false")
		} else {
			tester.t.Helper()
			tester.t.Fatal("expected end, but Next() returned true")
		}
	}
}

func (tester *LineIteratorTester) previous(expected bool) {
	actual := tester.it.Previous()
	if actual != expected {
		if expected {
			tester.t.Helper()
			tester.t.Fatal("expected line, but Previous() returned false")
		} else {
			tester.t.Helper()
			tester.t.Fatal("expected end, but Previous() returned true")
		}
	}
}

func (tester *LineIteratorTester) line(expected string) {
	actual := tester.it.Text()
	if actual != expected {
		tester.t.Helper()
		tester.t.Errorf("expected %q, but found %q", expected, actual)
	}
}

func (tester *LineIteratorTester) valid(expected bool) {
	actual := tester.it.Valid()
	if actual != expected {
		if expected {
			tester.t.Helper()
			tester.t.Fatal("expected valid, but Valid() returned false")
		} else {
			tester.t.Helper()
			tester.t.Fatal("expected invalid, but Valid() returned true")
		}
	}
}

func TestLineIteratorNextEmpty(t *testing.T) {
	it := MakeLineIteratorTester(t, "")

	it.next(true)
	it.line("")
	it.next(false)
}

func TestLineIteratorNextSingleNewline(t *testing.T) {
	it := MakeLineIteratorTester(t, "\n")

	it.next(true)
	it.line("")
	it.next(true)
	it.line("")
	it.next(false)
}

func TestLineIteratorNextNonEmpty(t *testing.T) {
	it := MakeLineIteratorTester(t, "a")

	it.next(true)
	it.line("a")
	it.next(false)
}

func TestLineIteratorNextTwoLines(t *testing.T) {
	it := MakeLineIteratorTester(t, "a\nb")

	it.next(true)
	it.line("a")
	it.next(true)
	it.line("b")
	it.next(false)
}

func TestLineIteratorNextTrailingNewline(t *testing.T) {
	it := MakeLineIteratorTester(t, "a\n")

	it.next(true)
	it.line("a")
	it.next(true)
	it.line("")
	it.next(false)
}

func TestLineIteratorNextLong(t *testing.T) {
	it := MakeLineIteratorTester(t, "abc\ndef\nxyz")

	it.next(true)
	it.line("abc")
	it.next(true)
	it.line("def")
	it.next(true)
	it.line("xyz")
	it.next(false)
}

func TestLineIteratorNextLongTrailingNewline(t *testing.T) {
	it := MakeLineIteratorTester(t, "abc\ndef\nxyz\n")

	it.next(true)
	it.line("abc")
	it.next(true)
	it.line("def")
	it.next(true)
	it.line("xyz")
	it.next(true)
	it.line("")
	it.next(false)
}

func TestLineIteratorNextTrimCarriageReturn(t *testing.T) {
	it := MakeLineIteratorTester(t, "abc\r\nxyz\n")

	it.next(true)
	it.line("abc")
	it.next(true)
	it.line("xyz")
	it.next(true)
	it.line("")
	it.next(false)
}

func TestLineIteratorNextTrimCarriageReturnLastLine(t *testing.T) {
	it := MakeLineIteratorTester(t, "a\nb\r")

	it.next(true)
	it.line("a")
	it.next(true)
	it.line("b")
	it.next(false)
}

func TestLineIteratorNextTrimCarriageReturnOnly(t *testing.T) {
	it := MakeLineIteratorTester(t, "\r")

	it.next(true)
	it.line("")
	it.next(false)
}

func TestLineIteratorPreviousEmpty(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "")

	it.previous(true)
	it.line("")
	it.previous(false)
}

func TestLineIteratorPreviousSingleNewline(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "\n")

	it.previous(true)
	it.line("")
	it.previous(true)
	it.line("")
	it.previous(false)
}

func TestLineIteratorPreviousNonEmpty(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "a")

	it.previous(true)
	it.line("a")
	it.previous(false)
}

func TestLineIteratorPreviousTwoLines(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "a\nb")

	it.previous(true)
	it.line("b")
	it.previous(true)
	it.line("a")
	it.previous(false)
}

func TestLineIteratorPreviousTrailingNewline(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "a\n")

	it.previous(true)
	it.line("")
	it.previous(true)
	it.line("a")
	it.previous(false)
}

func TestLineIteratorPreviousLong(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "abc\ndef\nxyz")

	it.previous(true)
	it.line("xyz")
	it.previous(true)
	it.line("def")
	it.previous(true)
	it.line("abc")
	it.previous(false)
}

func TestLineIteratorPreviousLongTrailingNewline(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "abc\ndef\nxyz\n")

	it.previous(true)
	it.line("")
	it.previous(true)
	it.line("xyz")
	it.previous(true)
	it.line("def")
	it.previous(true)
	it.line("abc")
	it.previous(false)
}

func TestLineIteratorPreviousTrimCarriageReturn(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "abc\r\nxyz\n")

	it.previous(true)
	it.line("")
	it.previous(true)
	it.line("xyz")
	it.previous(true)
	it.line("abc")
	it.previous(false)
}

func TestLineIteratorPreviousTrimCarriageReturnLastLine(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "a\nb\r")

	it.previous(true)
	it.line("b")
	it.previous(true)
	it.line("a")
	it.previous(false)
}

func TestLineIteratorPreviousTrimCarriageReturnOnly(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "\r")

	it.previous(true)
	it.line("")
	it.previous(false)
}

func TestLineIteratorNextValid(t *testing.T) {
	it := MakeLineIteratorTester(t, "a")

	it.valid(false)
	it.next(true)
	it.valid(true)
	it.next(false)
	it.valid(false)
}

func TestLineIteratorPreviousValid(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "a")

	it.valid(false)
	it.previous(true)
	it.valid(true)
	it.previous(false)
	it.valid(false)
}

func TestLineIteratorLineSpecial(t *testing.T) {
	it := MakeLineIteratorTester(t, "a")

	it.line("")
	it.next(true)
	it.next(false)
	it.line("")
}

func TestLineIteratorMixed(t *testing.T) {
	it := MakeLineIteratorTester(t, "a\nb\nc")

	it.next(true)
	it.line("a")
	it.next(true)
	it.line("b")
	it.next(true)
	it.line("c")
	it.previous(true)
	it.line("b")
	it.previous(true)
	it.line("a")
	it.next(true)
	it.line("b")
	it.next(true)
	it.line("c")
	it.next(false)
}

func TestLineIteratorSeekStart(t *testing.T) {
	it := MakeLineIteratorTesterEnd(t, "a\nb\nc")

	it.it.SeekStart()
	it.next(true)
	it.line("a")
	it.next(true)
	it.line("b")
	it.next(true)
	it.line("c")
	it.next(false)
}

func TestLineIteratorSeekEnd(t *testing.T) {
	it := MakeLineIteratorTester(t, "a\nb\nc")

	it.it.SeekEnd()
	it.previous(true)
	it.line("c")
	it.previous(true)
	it.line("b")
	it.previous(true)
	it.line("a")
	it.previous(false)
}

func BenchmarkBytes(b *testing.B) {
	var buffer bytes.Buffer
	for i := 0; i < 1000; i++ {
		buffer.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ\n")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it := MakeLineIterator(buffer.Bytes())
		for it.Next() {
			it.Bytes()
		}
	}
}

func BenchmarkText(b *testing.B) {
	var buffer bytes.Buffer
	for i := 0; i < 1000; i++ {
		buffer.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ\n")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it := MakeLineIterator(buffer.Bytes())
		for it.Next() {
			it.Text()
		}
	}
}
