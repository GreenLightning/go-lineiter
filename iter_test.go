package lineiter

import (
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

func (tester *LineIteratorTester) line(expected string) {
	actual := tester.it.Text()
	if actual != expected {
		tester.t.Helper()
		tester.t.Errorf("expected %q, but found %q", expected, actual)
	}
}

func (tester *LineIteratorTester) expect(expected, actual int, name string) {
	if actual != expected {
		tester.t.Helper()
		tester.t.Errorf("expected %s=%d, but found %d", name, expected, actual)
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
