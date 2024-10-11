package main

import (
	"fmt"
	"os"
	"strings"

	tui "github.com/mikolajgs/terminal-ui"
)

type gameInterface struct {
	t           *tui.TUI
	words       *tui.TUIPane
	leftLetter  *tui.TUIPane
	rightLetter *tui.TUIPane
	score       *tui.TUIPane
	leftTop     *tui.TUIPane
	g           *game
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
		p.Write(4, 1, gi.g.getLeftLetter(), false)
		return 0
	})
	gi.rightLetter.SetOnDraw(func(p *tui.TUIPane) int {
		p.Write(3, 0, "->", false)
		p.Write(3, 1, gi.g.getRightLetter(), false)
		return 0
	})
	gi.score.SetOnIterate(func(p *tui.TUIPane) int {
		p.Write(0, 0, "Correct:", false)
		p.Write(1, 1, fmt.Sprintf("%d/%d", g.getNumberOfCorrectAnswers(), g.getNumberOfUsedWords()), false)
		p.Write(0, 3, "Total:", false)
		p.Write(1, 4, fmt.Sprintf("%d", g.getNumberOfAllWords()), false)
		return 0
	})

	gi.initKeyboard()

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

func (gi *gameInterface) initIteration() {
	f := func(p *tui.TUIPane) int {
		if !gi.g.isStarted() {
			p.Write(2, 1, "Press the S key to start the game", false)
			return NOT_STARTED
		}

		gi.g.setAvailableLines(p.GetHeight())

		// Run game loop iteration and update UI depending on the result
		r := gi.g.iterate()
		if r == GAME_OVER {
			p.Write(2, 0, "** Game over! **", false)
			return r
		}

		l := gi.g.getCurrentLine()
		if r == CONTINUE_GAME || r == INCORRECT_GUESS || r == CORRECT_GUESS {
			// Write word
			if l > 0 {
				gi.clearPaneLine(gi.words, l)
			}
			gi.writeWord(gi.g.getCurrentWord(), l+1)
		}

		if r == CORRECT_GUESS {
			gi.clearPaneLine(gi.words, l+1)
		}

		return r
	}

	gi.words.SetOnDraw(f)
	gi.words.SetOnIterate(f)
}

func (gi *gameInterface) setSpeed(i int) {
	gi.t.SetLoopSleep(i)
}

func (gi *gameInterface) initKeyboard() {
	gi.t.SetOnKeyPress(func(t *tui.TUI, b []byte) {
		if string(b) == "x" {
			t.Exit(0)
		}
		if string(b) == "s" {
			if !gi.g.isStarted() {
				gi.clearPane(gi.words)

				gi.g.startGame()
			}
			return
		}
		// TODO: Keys should be handled differently, maybe in raw mode
		// left arrow pressed
		if string(b) == "D" {
			gi.g.setCurrentWordWithLeftLetter()
			gi.writeCurrentWord()
			return
		}
		// right arrow pressed
		if string(b) == "C" {
			gi.g.setCurrentWordWithRightLetter()
			gi.writeCurrentWord()
			return
		}
		// down arrow pressed
		if string(b) == "B" {
			l := gi.g.getCurrentLine()
			if l == gi.g.getLastLine() {
				return
			}
			if l > 0 {
				gi.clearPaneLine(gi.words, l+1)
			}
			gi.g.setNextLineToLast()
			gi.writeCurrentWord()
		}
	})
}

func (gi *gameInterface) writeCurrentWord() {
	l := gi.g.getCurrentLine()
	if l > 0 {
		gi.clearPaneLine(gi.words, l)
	}
	gi.writeWord(gi.g.getCurrentWord(), l+1)
}

func (gi *gameInterface) clearPane(p *tui.TUIPane) {
	for y := 0; y < p.GetHeight()-2; y++ {
		gi.clearPaneLine(p, y)
	}
}

func (gi *gameInterface) clearPaneLine(p *tui.TUIPane, y int) {
	p.Write(0, y, strings.Repeat(" ", p.GetWidth()-2), false)
}

func (gi *gameInterface) writeWord(w string, l int) {
	gi.words.Write((gi.words.GetWidth()-2-len(w))/2, l, w, false)
}

func (gi *gameInterface) run() {
	gi.t.Run(os.Stdout, os.Stderr)
}
