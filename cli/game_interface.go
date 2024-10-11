package main

import (
	"fmt"
	"os"
	"strings"

	tui "github.com/mikolajgs/terminal-ui"
)

type gameInterface struct {
	t *tui.TUI
	words *tui.TUIPane
	leftLetter *tui.TUIPane
	rightLetter *tui.TUIPane
	score *tui.TUIPane
	leftTop *tui.TUIPane
	g *game
}

func newGameInterface(g *game) *gameInterface {
	gi := &gameInterface{}

	gi.g = g
	gi.t = tui.NewTUI()
	p := gi.t.GetPane()

	pLeft, pMiddleAndRight := p.SplitVertically(-10, tui.UNIT_CHAR)
	pMiddle, pRight := pMiddleAndRight.SplitVertically(10, tui.UNIT_CHAR)
	pLeftTop, pLeftBottom := pLeft.SplitHorizontally(4, tui.UNIT_CHAR)
	pRightTop, pRightBottom := pRight.SplitHorizontally(4, tui.UNIT_CHAR)

	gi.words = pMiddle
	gi.leftLetter = pLeftBottom
	gi.rightLetter = pRightBottom
	gi.score = pRightTop
	gi.leftTop = pLeftTop

	gi.initStyle()

	gi.initIteration()

	gi.leftLetter.SetOnDraw(func(p *tui.TUIPane) int {
		p.Write(3, 0, "<-", false)
		p.Write(4, 1, gi.g.letters[0], false)
		return 0
	})
	gi.rightLetter.SetOnDraw(func(p *tui.TUIPane) int {
		p.Write(3, 0, "->", false)
		p.Write(3, 1, gi.g.letters[1], false)
		return 0
	})
	gi.score.SetOnIterate(func(p *tui.TUIPane) int {
		p.Write(0, 0, "Correct:", false)
		p.Write(1, 1, fmt.Sprintf("%d/%d", gi.g.wordsGiven - len(gi.g.wordsNotGuessed), gi.g.wordsGiven), false)
		return 0
	})

	gi.initKeyboard(g)

	return gi
}

func (gi *gameInterface) initStyle() {
	s := tui.NewTUIPaneStyleFrame()
	gi.words.SetStyle(s)
	gi.leftLetter.SetStyle(s)
	gi.rightLetter.SetStyle(s)
	gi.score.SetStyle(s)
	gi.leftTop.SetStyle(s)
}

func (gi *gameInterface) setSpeed(i int) {
	gi.t.SetLoopSleep(i)
}

func (gi *gameInterface) initKeyboard(g *game) {
	gi.t.SetOnKeyPress(func(t *tui.TUI, b []byte) {
		if string(b) == "x" {
			t.Exit(0)
		}
		if string(b) == "s" {
			if !g.started {
				gi.startGame()
			}
		}
		// TODO: Keys should be handled differently, maybe in raw mode
		// left arrow pressed
		if string(b) == "D" {
			g.currentWord = strings.Replace(g.currentWordTemplate, "_", g.letters[0], 1)
			gi.clearLineBeforeWord()
			gi.writeCurrentWord()
		}
		// right arrow pressed 
		if string(b) == "C" {
			g.currentWord = strings.Replace(g.currentWordTemplate, "_", g.letters[1], 1)
			gi.clearLineBeforeWord()
			gi.writeCurrentWord()
		}
		// down arrow pressed
		if string(b) == "B" {
			gi.clearLineBeforeWord()
			g.nextWordLine = g.lastAvailableLine-1
			gi.writeCurrentWord()
		}
	})
}

func (gi *gameInterface) initIteration() {
	f := func(p *tui.TUIPane) int {
		if !gi.g.started {
			p.Write(2, 1, "Naciśnij 's' aby zacząć grę", false)
			return 0
		}

		// If there is no word then take the next one
		if gi.g.currentWord == "" {
			currentWordArr := strings.Split(gi.g.words[gi.g.nextWordIndex], ":")
			gi.g.currentWordTemplate = currentWordArr[0]
			gi.g.currentWord = gi.g.currentWordTemplate
			gi.g.currentWordCorrect = strings.Replace(gi.g.currentWordTemplate, "_", currentWordArr[1], 1)
			gi.g.nextWordIndex++
			gi.g.nextWordLine = 0
		}

		// We need a position that is at the very bottom
		gi.g.lastAvailableLine = p.GetHeight()-2-len(gi.g.wordsNotGuessed)

		if gi.g.lastAvailableLine == 0 || gi.g.nextWordIndex == len(gi.g.words) {
			p.Write(2, 0, "** Koniec gry! **", false)
			gi.g.started = false
			return 2
		}

		// Draw word
		gi.clearLineBeforeWord()
		gi.writeCurrentWord()

		// If the word is already in the last line
		if gi.g.nextWordLine == gi.g.lastAvailableLine-1 {
			gi.g.wordsGiven++
			if gi.g.currentWord != gi.g.currentWordCorrect {
				gi.g.wordsNotGuessed = append(gi.g.wordsNotGuessed, gi.g.currentWord)
			} else {
				gi.clearPaneLine(gi.words, gi.g.nextWordLine)
			}
			gi.g.currentWord = ""
			return 1
		}

		// Increment the line for the next iteration
		gi.g.nextWordLine++

		return 0
	}
	gi.words.SetOnDraw(f)
	gi.words.SetOnIterate(f)

}

func (gi *gameInterface) clearPane(p *tui.TUIPane) {
	for y := 0; y < p.GetHeight()-2; y++ {
		gi.clearPaneLine(p, y)
	}
}

func (gi *gameInterface) clearPaneLine(p *tui.TUIPane, y int) {
	p.Write(0, y, strings.Repeat(" ", p.GetWidth()-2), false)
}

func (gi *gameInterface) writeCurrentWord() {
	wordLen := len(gi.g.currentWord)
	leftMargin := (gi.words.GetWidth()-2-wordLen)/2
	gi.words.Write(leftMargin, gi.g.nextWordLine, gi.g.currentWord, false)
}

func (gi *gameInterface) clearLineBeforeWord() {
	lineToDrawOn := gi.g.nextWordLine
	if gi.g.nextWordLine > 0 {
		gi.clearPaneLine(gi.words, lineToDrawOn-1)
	}
	gi.clearPaneLine(gi.words, lineToDrawOn)
}

func (gi *gameInterface) run() {
	gi.t.Run(os.Stdout, os.Stderr)
}

func (gi *gameInterface) startGame() {
	gi.clearPane(gi.words)
	gi.g.started = true
	gi.g.nextWordIndex = 0
	gi.g.currentWord = ""
	gi.g.nextWordLine = 0
	gi.g.wordsNotGuessed = []string{}
	gi.g.wordsGiven = 0
}
