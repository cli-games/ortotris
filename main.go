package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/mikolajgs/broccli"

	tui "github.com/mikolajgs/terminal-ui"
)

var words []string
var letters [2]string
var gameStarted bool
var nextWordIndex int
var currentWordTemplate string
var currentWord string
var currentWordCorrect string
var nextWordLine int // 0 - top, 20 - bottom
var wordsNotGuessed []string
var centralPane *tui.TUIPane
var leftBottomPane *tui.TUIPane
var rightBottomPane *tui.TUIPane
var rightTopPane *tui.TUIPane
var lastAvailableLine int
var wordsGiven int

func versionHandler(c *broccli.CLI) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func main() {
	cli := broccli.NewCLI("ortotris", "Klon ortotrisa", "")
	cmd := cli.AddCmd("zagraj", "Załącza grę", startHandler)
	cmd.AddFlag("plik", "f", "", "Plik z wyrazami", broccli.TypePathFile, broccli.IsExistent|broccli.IsRequired)
	_ = cli.AddCmd("wersja", "Pokazuje wersję programu", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--wersja") {
		os.Args = []string{"App", "version"}
	}
	os.Exit(cli.Run())
}

func startHandler(c *broccli.CLI) int {
	// Create UI
	t := tui.NewTUI()
	p := t.GetPane()

	pLeft, pMiddleAndRight := p.SplitVertically(-10, tui.UNIT_CHAR)
	pMiddle, pRight := pMiddleAndRight.SplitVertically(10, tui.UNIT_CHAR)
	pLeftTop, pLeftBottom := pLeft.SplitHorizontally(4, tui.UNIT_CHAR)
	pRightTop, pRightBottom := pRight.SplitHorizontally(4, tui.UNIT_CHAR)

	centralPane = pMiddle
	leftBottomPane = pLeftBottom
	rightBottomPane = pRightBottom
	rightTopPane = pRightTop

	s := tui.NewTUIPaneStyleFrame()

	for _, pn := range []*tui.TUIPane{pLeftTop, leftBottomPane, pRightTop, rightBottomPane, centralPane} {
		pn.SetStyle(s)
	}

	centralPane.SetOnDraw(drawMiddle())
	centralPane.SetOnIterate(drawMiddle())

	leftBottomPane.SetOnDraw(func(p *tui.TUIPane) int {
		p.Write(3, 0, "<-", false)
		p.Write(4, 1, letters[0], false)
		return 0
	})
	rightBottomPane.SetOnDraw(func(p *tui.TUIPane) int {
		p.Write(3, 0, "->", false)
		p.Write(3, 1, letters[1], false)
		return 0
	})
	rightTopPane.SetOnIterate(func(p *tui.TUIPane) int {
		p.Write(0, 0, "Dobrze:", false)
		p.Write(1, 1, fmt.Sprintf("%d/%d", wordsGiven - len(wordsNotGuessed), wordsGiven), false)
		return 0
	})

	t.SetLoopSleep(300)

	t.SetOnKeyPress(func(t *tui.TUI, b []byte) {
		if string(b) == "x" {
			t.Exit(0)
		}
		if string(b) == "s" {
			if !gameStarted {
				startGame()
			}
		}
		// TODO: Keys should be handled differently, maybe in raw mode
		// left arrow pressed
		if string(b) == "D" {
			currentWord = strings.Replace(currentWordTemplate, "_", letters[0], 1)
			clearLineBeforeWord()
			writeCurrentWord()
		}
		// right arrow pressed 
		if string(b) == "C" {
			currentWord = strings.Replace(currentWordTemplate, "_", letters[1], 1)
			clearLineBeforeWord()
			writeCurrentWord()
		}
		// down arrow pressed
		if string(b) == "B" {
			clearLineBeforeWord()
			nextWordLine = lastAvailableLine-1
			writeCurrentWord()
		}
	})

	readWordsFromFile(c.Flag("plik"))
	randomizeWords()

	// Run UI
	t.Run(os.Stdout, os.Stderr)

	return 0
}

func drawMiddle() func(*tui.TUIPane) int {
	fn := func(p *tui.TUIPane) int {
		if !gameStarted {
			p.Write(2, 1, "Naciśnij 's' aby zacząć grę", false)
			return 0
		}

		// If there is no word then take the next one
		if currentWord == "" {
			currentWordArr := strings.Split(words[nextWordIndex], ":")
			currentWordTemplate = currentWordArr[0]
			currentWord = currentWordTemplate
			currentWordCorrect = strings.Replace(currentWordTemplate, "_", currentWordArr[1], 1)
			nextWordIndex++
			nextWordLine = 0
		}

		// We need a position that is at the very bottom
		lastAvailableLine = p.GetHeight()-2-len(wordsNotGuessed)

		if lastAvailableLine == 0 || nextWordIndex == len(words) {
			p.Write(2, 0, "** Koniec gry! **", false)
			gameStarted = false
			return 2
		}

		// Draw word
		clearLineBeforeWord()
		writeCurrentWord()

		// If the word is already in the last line
		if nextWordLine == lastAvailableLine-1 {
			wordsGiven++
			if currentWord != currentWordCorrect {
				wordsNotGuessed = append(wordsNotGuessed, currentWord)
			} else {
				clearPaneLine(centralPane, nextWordLine)
			}
			currentWord = ""
			return 1
		}

		// Increment the line for the next iteration
		nextWordLine++

		return 0
	}
	return fn
}

func clearPane(p *tui.TUIPane) {
	for y := 0; y < p.GetHeight()-2; y++ {
		clearPaneLine(p, y)
	}
}

func clearPaneLine(p *tui.TUIPane, y int) {
	p.Write(0, y, strings.Repeat(" ", p.GetWidth()-2), false)
}

func randomizeWords() {
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
}

func startGame() {
	clearPane(centralPane)
	gameStarted = true
	nextWordIndex = 0
	currentWord = ""
	nextWordLine = 0
	wordsNotGuessed = []string{}
	wordsGiven = 0
}

func readWordsFromFile(fp string) {
	fn := fp
	f, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Błąd przy otwieraniu pliku z wyrazami %s: %s", fn, err.Error())
		os.Exit(1)
	}
	defer f.Close()

	// TODO: Validation - for now, code assumes that the file contains correct data
	i := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		i++
		if i == 1 {
			lineArr := strings.Split(line, ":")
			letters = [2]string{lineArr[0], lineArr[1]}
			continue
		}
		words = append(words, line)
	}
}

func writeCurrentWord() {
	wordLen := len(currentWord)
	leftMargin := (centralPane.GetWidth()-2-wordLen)/2
	centralPane.Write(leftMargin, nextWordLine, currentWord, false)
}

func clearLineBeforeWord() {
	lineToDrawOn := nextWordLine
	if nextWordLine > 0 {
		clearPaneLine(centralPane, lineToDrawOn-1)
	}
	clearPaneLine(centralPane, lineToDrawOn)
}
