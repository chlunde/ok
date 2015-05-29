package main

import (
	"regexp"
	"unicode/utf8"
)

type Editor struct {
	Text   string
	Cursor int
	Word   *regexp.Regexp
}

func NewEditor() Editor {
	return Editor{Word: regexp.MustCompile(`\s+|-|\.|,`)}
}

func (editor *Editor) State() (string, int) {
	return editor.Text, editor.Cursor
}

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

func (editor *Editor) Insert(r rune) {
	s, c := editor.State()
	runes := []rune(s)
	editor.Update(string(runes[:c])+string(r)+string(runes[c:]), c+1)
}

func (editor *Editor) Remove() bool {
	s, c := editor.State()
	runes := []rune(s)
	if len(s) == 0 || c >= len(runes) {
		return false
	}
	editor.Update(string(runes[:c])+string(runes[c+1:]), c)
	return true
}

func (editor *Editor) RemoveBackwards() bool {
	s, c := editor.State()
	if c == 0 || len(s) == 0 {
		return false
	}
	runes := []rune(s)
	editor.Update(string(runes[:c-1])+string(runes[c:]), c-1)
	return true
}

func (editor *Editor) RemoveToEnd() bool {
	s, c := editor.State()
	if len(s) == 0 || c == len(s) {
		return false
	}
	runes := []rune(s)
	editor.Update(string(runes[:c]), c)
	return true
}

func (editor *Editor) RemoveWord() bool {
	s, c := editor.State()
	if len(s) == 0 || c == 0 {
		return false
	}

	runes := []rune(s)
	text := runes[:c]
	word := editor.Word.FindAllStringIndex(string(text), -1)
	if word == nil {
		return editor.Clear()
	}
	end := word[len(word)-1]

	if len(word) > 1 && c <= end[1] {
		text = text[:word[len(word)-2][1]]
	} else if c <= end[1] {
		text = nil
	} else {
		text = text[:end[1]]
	}

	editor.Update(string(text)+string(runes[c:]), len(text))

	return true
}

func (editor *Editor) MoveStart() {
	editor.Cursor = 0
}

func (editor *Editor) MoveEnd() {
	editor.Cursor = len([]rune(editor.Text))
}

func (editor *Editor) MoveForward() {
	if editor.Cursor < len([]rune(editor.Text)) {
		editor.Cursor += 1
	}
}

func (editor *Editor) MoveBackward() {
	if editor.Cursor > 0 {
		editor.Cursor -= 1
	}
}
