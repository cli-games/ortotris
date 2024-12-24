package main

import (
	"fmt"
	"os"
	"strings"

	tui "github.com/go-phings/terminal-ui"
)

const Black = 2
const Green = 3
const Yellow = 4
const DarkBlue = 5
const DarkGreen = 6
const DarkYellow = 7
const DarkMagenta = 8
const DarkCyan = 9
const DarkGray = 10
const Red = 11
const Magenta = 12
const Cyan = 17

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
		p.Write(0, 0, 
			gi.wrapInColors("   <-   ", Black, DarkGreen),
			false,
		)
		p.Write(0, 1,
			gi.wrapInColors("    "+gi.g.getLeftLetter()+"   ", Black, DarkGreen),
			false,
		)
		return 0
	})
	gi.rightLetter.SetOnDraw(func(p *tui.TUIPane) int {
		p.Write(0, 0, 
			gi.wrapInColors("   ->   ", Black, DarkYellow),
			false,
		)
		p.Write(0, 1,
			gi.wrapInColors("   "+gi.g.getRightLetter()+"    ", Black, DarkYellow),
			false,
		)
		return 0
	})
	gi.score.SetOnIterate(func(p *tui.TUIPane) int {
		p.Write(0, 0, gi.wrapInColors("Correct:", Magenta, 0), false)
		p.Write(1, 1, fmt.Sprintf("%d/%d", g.getNumberOfCorrectAnswers(), g.getNumberOfUsedWords()), false)
		p.Write(0, 3, gi.wrapInColors("Total:", Magenta, 0), false)
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

	tl := &tui.TUIPaneStyle{
		NE: gi.wrapInColors("╗", Cyan, DarkCyan),
		N: gi.wrapInColors("═", Cyan, DarkCyan),
		NW: gi.wrapInColors("╔", Cyan, DarkCyan),
		W: gi.wrapInColors("║", Cyan, DarkCyan),
		SW: gi.wrapInColors("╚", Cyan, DarkCyan),
		S: gi.wrapInColors("═", Cyan, DarkCyan),
		SE: gi.wrapInColors("╝", Cyan, DarkCyan),
		E: gi.wrapInColors("║", Cyan, DarkCyan),
	}
	gi.leftTop.SetStyle(tl)

	cl := &tui.TUIPaneStyle{
		NE: gi.wrapInColors("╗", Green, DarkGreen),
		N: gi.wrapInColors("═", Green, DarkGreen),
		NW: gi.wrapInColors("╔", Green, DarkGreen),
		W: gi.wrapInColors("║", Green, DarkGreen),
		SW: gi.wrapInColors("╚", Green, DarkGreen),
		S: gi.wrapInColors("═", Green, DarkGreen),
		SE: gi.wrapInColors("╝", Green, DarkGreen),
		E: gi.wrapInColors("║", Green, DarkGreen),
	}

	cr := &tui.TUIPaneStyle{
		NE: gi.wrapInColors("╗", Yellow, DarkYellow),
		N: gi.wrapInColors("═", Yellow, DarkYellow),
		NW: gi.wrapInColors("╔", Yellow, DarkYellow),
		W: gi.wrapInColors("║", Yellow, DarkYellow),
		SW: gi.wrapInColors("╚", Yellow, DarkYellow),
		S: gi.wrapInColors("═", Yellow, DarkYellow),
		SE: gi.wrapInColors("╝", Yellow, DarkYellow),
		E: gi.wrapInColors("║", Yellow, DarkYellow),
	}
	gi.leftLetter.SetStyle(cl)
	gi.rightLetter.SetStyle(cr)

	ss := &tui.TUIPaneStyle{
		NE: gi.wrapInColors("╗", Magenta, DarkMagenta),
		N: gi.wrapInColors("═", Magenta, DarkMagenta),
		NW: gi.wrapInColors("╔", Magenta, DarkMagenta),
		W: gi.wrapInColors("║", Magenta, DarkMagenta),
		SW: gi.wrapInColors("╚", Magenta, DarkMagenta),
		S: gi.wrapInColors("═", Magenta, DarkMagenta),
		SE: gi.wrapInColors("╝", Magenta, DarkMagenta),
		E: gi.wrapInColors("║", Magenta, DarkMagenta),
	}
	gi.score.SetStyle(ss)

	sw := &tui.TUIPaneStyle{
		NE: gi.wrapInColors("╗", Black, Red),
		N: gi.wrapInColors("═", Black, Red),
		NW: gi.wrapInColors("╔", Black, Red),
		W: gi.wrapInColors("║", Black, Red),
		SW: gi.wrapInColors("╚", Black, Red),
		S: gi.wrapInColors("═", Black, Red),
		SE: gi.wrapInColors("╝", Black, Red),
		E: gi.wrapInColors("║", Black, Red),
	}
	gi.words.SetStyle(sw)
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

func (gi gameInterface) wrapInColors(s string, fg int, bg int) string {
	f := ""
	b := ""
	r := "\033[0m"
	switch (fg) {
	case DarkCyan:
		f = "\033[0;36m"
	case Cyan:
		f = "\033[1;96m"
	case DarkMagenta:
		f = "\033[0;35m"
	case Black:
		f = "\033[0;30m"
	case Yellow:
		f = "\033[1;93m"
	case Green:
		f = "\033[1;92m"
	case Magenta:
		f = "\033[1;95m"
	default:
		f = ""
	}

	switch (bg) {
	case DarkBlue:
		b = "\033[44m"
	case DarkMagenta:
		b = "\033[45m"
	case DarkCyan:
		b = "\033[46m"
	case DarkGray:
		b = "\033[100m"
	case Red:
		b = "\033[101m"
	case DarkGreen:
		b = "\033[42m"
	case DarkYellow:
		b = "\033[43m"
	case Black:
		b = "\033[40m"
	case Yellow:
		b = "\033[103m"
	default:
		b = ""
	}

	return fmt.Sprintf("%s%s%s%s", f, b, s, r)
}
