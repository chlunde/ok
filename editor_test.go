package main

import "testing"

/*
type Editor struct {
	Text   string
	Cursor int
	Word   *regexp.Regexp
}
*/

/*
func TestNewEditor(t *testing.T) {
	return Editor{Word: regexp.MustCompile(`\s+|-|\.|,`)}
}

func TestEditorState(t *testing.T) (string, int) {
	return editor.Text, editor.Cursor
}
*/

/*
func (editor *Editor) Clear() bool {
	if len(editor.Text) == 0 {
		return false
	}
	editor.Update("", 0)
	return true
}

func (editor *Editor) Update(text string, cursor int) {
	editor.Text, editor.Cursor = text, cursor
	if editor.Cursor < 0 {
		editor.Cursor = utf8.RuneCountInString(editor.Text)
	}
}
*/

func TestEditorInsert(t *testing.T) {
	e := NewEditor()
	e.Insert('a')
	if e.Text != "a" {
		t.Error("initial insert failed: a", e.Text)
	}

	e.Insert('b')
	if e.Text != "ab" {
		t.Fatal("append failed: ab", e.Text)
	}

	e.MoveStart()
	e.Insert('c')
	if e.Text != "cab" {
		t.Fatal("prepend failed: cab", e.Text)
	}

	e.Insert('d')
	if e.Text != "cdab" {
		t.Fatal("insert in the middle failed: cdab", e.Text)
	}
}

func TestEditorInsertUnicode(t *testing.T) {
	answer := []string{
		"æabc",
		"aæbc",
		"abæc",
		"abcæ",
	}

	for pos := 0; pos <= 3; pos++ {
		e := NewEditor()
		e.Insert('a')
		e.Insert('b')
		e.Insert('c')
		if e.Text != "abc" {
			t.Error("test setup failed: abc", e.Text)
		}

		r := 'æ'
		e.MoveStart()
		for i := 0; i < pos; i++ {
			e.MoveForward()
		}
		e.Insert(r)

		if e.Text != answer[pos] {
			t.Errorf("inserting %v at pos %v failed: %v vs %v", r, pos, e.Text, answer[pos])
		}
	}
}

func TestEditorInsertUnicodeInUnicode(t *testing.T) {
	answer := []string{
		"æaåc",
		"aæåc",
		"aåæc",
		"aåcæ",
	}

	for pos := 0; pos <= 3; pos++ {
		e := NewEditor()
		e.Insert('a')
		e.Insert('å')
		e.Insert('c')
		if e.Text != "aåc" {
			t.Error("test setup failed: abc", e.Text)
		}

		r := 'æ'
		e.MoveStart()
		for i := 0; i < pos; i++ {
			e.MoveForward()
		}
		e.Insert(r)

		if e.Text != answer[pos] {
			t.Errorf("inserting %v at pos %v failed: %v vs %v", r, pos, e.Text, answer[pos])
		}
	}
}

func TestEditorRemove(t *testing.T) {
	answer := []string{
		"bc",
		"ac",
		"ab",
		"abc",
	}

	for pos := 0; pos <= 3; pos++ {
		e := NewEditor()
		e.Insert('a')
		e.Insert('b')
		e.Insert('c')
		if e.Text != "abc" {
			t.Error("test setup failed: abc", e.Text)
		}

		e.MoveStart()
		for i := 0; i < pos; i++ {
			e.MoveForward()
		}
		e.Remove()

		if e.Text != answer[pos] {
			t.Errorf("remove pos %v failed: %v vs %v", pos, e.Text, answer[pos])
		}
	}
}

func TestEditorRemoveUnicode(t *testing.T) {
	answer := []string{
		"åc",
		"ac",
		"aå",
		"aåc",
	}

	for pos := 0; pos <= 3; pos++ {
		e := NewEditor()
		e.Insert('a')
		e.Insert('å')
		e.Insert('c')
		if e.Text != "aåc" {
			t.Error("test setup failed: abc", e.Text)
		}

		e.MoveStart()
		for i := 0; i < pos; i++ {
			e.MoveForward()
		}
		e.Remove()

		if e.Text != answer[pos] {
			t.Errorf("remove pos %v failed: %v vs %v", pos, e.Text, answer[pos])
		}
	}
}

func TestEditorRemoveBackward(t *testing.T) {
	answer := []string{
		"abc",
		"bc",
		"ac",
		"ab",
	}

	for pos := 0; pos <= 3; pos++ {
		e := NewEditor()
		e.Insert('a')
		e.Insert('b')
		e.Insert('c')
		if e.Text != "abc" {
			t.Error("test setup failed: abc", e.Text)
		}

		e.MoveStart()
		for i := 0; i < pos; i++ {
			e.MoveForward()
		}
		e.RemoveBackwards()

		if e.Text != answer[pos] {
			t.Errorf("remove pos %v failed: %v vs %v", pos, e.Text, answer[pos])
		}
	}
}

func TestEditorRemoveBackwardsUnicode(t *testing.T) {
	answer := []string{
		"aåc",
		"åc",
		"ac",
		"aå",
	}

	for pos := 0; pos <= 3; pos++ {
		e := NewEditor()
		e.Insert('a')
		e.Insert('å')
		e.Insert('c')
		if e.Text != "aåc" {
			t.Error("test setup failed: abc", e.Text)
		}

		e.MoveStart()
		for i := 0; i < pos; i++ {
			e.MoveForward()
		}
		e.RemoveBackwards()

		if e.Text != answer[pos] {
			t.Errorf("remove pos %v failed: %v vs %v", pos, e.Text, answer[pos])
		}
	}
}

func TestEditorRemoveToEndUnicode(t *testing.T) {
	answer := []string{
		"",
		"a",
		"aå",
		"aåc",
	}

	for pos := 0; pos <= 3; pos++ {
		e := NewEditor()
		e.Insert('a')
		e.Insert('å')
		e.Insert('c')
		if e.Text != "aåc" {
			t.Error("test setup failed: abc", e.Text)
		}

		e.MoveStart()
		for i := 0; i < pos; i++ {
			e.MoveForward()
		}
		e.RemoveToEnd()

		if e.Text != answer[pos] {
			t.Errorf("remove pos %v failed: %v vs %v", pos, e.Text, answer[pos])
		}
	}
}

func TestEditorRemoveWordUnicode(t *testing.T) {
	answer := []string{
		"aå.åa-åå",
		"å.åa-åå",
		".åa-åå",
		"åa-åå",
		"aå.a-åå",
		"aå.-åå",
		"aå.åå",
		"aå.åa-åå",
		"aå.åa-å",
		"aå.åa-å",
		"aå.åa-",
	}

	for pos := 0; pos <= 9; pos++ {
		e := NewEditor()
		e.Insert('a')
		e.Insert('å')
		e.Insert('.')
		e.Insert('å')
		e.Insert('a')
		e.Insert('-')
		e.Insert('å')
		e.Insert('å')

		if e.Text != "aå.åa-åå" {
			t.Error("test setup failed: abc", e.Text)
		}

		e.MoveStart()
		for i := 0; i < pos; i++ {
			e.MoveForward()
		}
		e.RemoveWord()

		if e.Text != answer[pos] {
			t.Errorf("removeword pos %v failed: %v vs %v", pos, e.Text, answer[pos])
		}
	}
}

/*
func (editor *Editor) MoveStart() // implicit test
func (editor *Editor) MoveForward() // implicit test

func (editor *Editor) MoveEnd()
func (editor *Editor) MoveBackward()
*/
