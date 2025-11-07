package editor

import "testing"

func TestInsertTextWithinLine(t *testing.T) {
	buf := NewBuffer("abc")
	buf.cursor.Line = 0
	buf.cursor.Col = 1

	buf.InsertText("XY")

	if got, want := buf.Line(0), "aXYbc"; got != want {
		t.Fatalf("line mismatch: got %q want %q", got, want)
	}
	if got, want := buf.cursor.Line, 0; got != want {
		t.Fatalf("cursor line got %d want %d", got, want)
	}
	if got, want := buf.cursor.Col, 3; got != want {
		t.Fatalf("cursor col got %d want %d", got, want)
	}
}

func TestInsertTextWithNewline(t *testing.T) {
	buf := NewBuffer("abcd")
	buf.cursor.Col = 2

	buf.InsertText("X\nY")

	if got, want := buf.LineCount(), 2; got != want {
		t.Fatalf("line count got %d want %d", got, want)
	}
	if got, want := buf.Line(0), "abX"; got != want {
		t.Fatalf("line 0 got %q want %q", got, want)
	}
	if got, want := buf.Line(1), "Ycd"; got != want {
		t.Fatalf("line 1 got %q want %q", got, want)
	}
	if buf.cursor.Line != 1 {
		t.Fatalf("cursor line got %d want 1", buf.cursor.Line)
	}
	if buf.cursor.Col != 1 {
		t.Fatalf("cursor col got %d want 1", buf.cursor.Col)
	}
}

func TestDeleteBackwardWithinLine(t *testing.T) {
	buf := NewBuffer("abcd")
	buf.cursor.Col = 2

	if !buf.DeleteBackward() {
		t.Fatalf("expected delete backward success")
	}
	if got, want := buf.Line(0), "acd"; got != want {
		t.Fatalf("line got %q want %q", got, want)
	}
	if buf.cursor.Col != 1 {
		t.Fatalf("cursor col got %d want 1", buf.cursor.Col)
	}
}

func TestDeleteBackwardAtLineStartMerges(t *testing.T) {
	buf := NewBuffer("ab\ncd")
	buf.cursor.Line = 1
	buf.cursor.Col = 0

	if !buf.DeleteBackward() {
		t.Fatalf("expected merge delete success")
	}
	if got, want := buf.LineCount(), 1; got != want {
		t.Fatalf("line count got %d want %d", got, want)
	}
	if got, want := buf.Line(0), "abcd"; got != want {
		t.Fatalf("line got %q want %q", got, want)
	}
	if buf.cursor.Line != 0 || buf.cursor.Col != 2 {
		t.Fatalf("cursor got %d:%d want 0:2", buf.cursor.Line, buf.cursor.Col)
	}
}

func TestDeleteForward(t *testing.T) {
	buf := NewBuffer("ab\ncd")
	buf.cursor.Line = 0
	buf.cursor.Col = 2

	if !buf.DeleteForward() {
		t.Fatalf("expected delete forward success")
	}
	if got, want := buf.LineCount(), 1; got != want {
		t.Fatalf("line count got %d want %d", got, want)
	}
	if got, want := buf.Line(0), "abcd"; got != want {
		t.Fatalf("line got %q want %q", got, want)
	}
	if buf.cursor.Line != 0 || buf.cursor.Col != 2 {
		t.Fatalf("cursor got %d:%d want 0:2", buf.cursor.Line, buf.cursor.Col)
	}
}

func TestDeleteLinesRange(t *testing.T) {
	buf := NewBuffer("l1\nl2\nl3\nl4")
	buf.DeleteLines(1, 2)

	if got, want := buf.LineCount(), 2; got != want {
		t.Fatalf("line count got %d want %d", got, want)
	}
	if got, want := buf.Line(0), "l1"; got != want {
		t.Fatalf("line0 got %q want %q", got, want)
	}
	if got, want := buf.Line(1), "l4"; got != want {
		t.Fatalf("line1 got %q want %q", got, want)
	}
	if buf.cursor.Line != 1 {
		t.Fatalf("cursor line got %d want 1", buf.cursor.Line)
	}
}

func TestDeleteLinesAll(t *testing.T) {
	buf := NewBuffer("l1")
	buf.DeleteLines(0, 0)

	if got := buf.LineCount(); got != 1 {
		t.Fatalf("expected 1 line after delete, got %d", got)
	}
	if got := buf.Line(0); got != "" {
		t.Fatalf("expected empty line after delete, got %q", got)
	}
	if buf.cursor.Line != 0 {
		t.Fatalf("cursor line got %d want 0", buf.cursor.Line)
	}
}

func TestLinePrefix(t *testing.T) {
	buf := NewBuffer("héllo\nworld")
	if got := buf.LinePrefix(0, 2); got != "hé" {
		t.Fatalf("prefix got %q want %q", got, "hé")
	}
	if got := buf.LinePrefix(1, 10); got != "world" {
		t.Fatalf("prefix beyond length got %q want %q", got, "world")
	}
	if got := buf.LinePrefix(0, 0); got != "" {
		t.Fatalf("zero prefix got %q want empty", got)
	}
}

func TestLinesRange(t *testing.T) {
	buf := NewBuffer("l1\nl2\nl3")
	lines := buf.LinesRange(0, 1)
	if len(lines) != 2 || lines[0] != "l1" || lines[1] != "l2" {
		t.Fatalf("lines range got %v", lines)
	}
	lines[0] = "xx"
	if buf.Line(0) != "l1" {
		t.Fatalf("lines range did not copy, line 0 mutated to %q", buf.Line(0))
	}
	if got := buf.LinesRange(5, 6); len(got) != 0 {
		t.Fatalf("out of range should return empty slice, got %v", got)
	}
}

func TestInsertLines(t *testing.T) {
	buf := NewBuffer("l1\nl2")
	buf.InsertLines(1, []string{"X", "Y"})

	if got, want := buf.LineCount(), 4; got != want {
		t.Fatalf("line count got %d want %d", got, want)
	}
	if got, want := buf.Line(0), "l1"; got != want {
		t.Fatalf("line0 got %q want %q", got, want)
	}
	if got, want := buf.Line(1), "X"; got != want {
		t.Fatalf("line1 got %q want %q", got, want)
	}
	if got, want := buf.Line(2), "Y"; got != want {
		t.Fatalf("line2 got %q want %q", got, want)
	}
	if got, want := buf.Line(3), "l2"; got != want {
		t.Fatalf("line3 got %q want %q", got, want)
	}
	if buf.cursor.Line != 2 {
		t.Fatalf("cursor line expected 2 got %d", buf.cursor.Line)
	}
}
