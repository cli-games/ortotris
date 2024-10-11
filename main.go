package main

import (
	"fmt"
	"os"

	"github.com/mikolajgs/broccli"
)

func versionHandler(c *broccli.CLI) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func main() {
	cli := broccli.NewCLI("ortotris", "Clone of a classic Ortotris game", "")
	cmd := cli.AddCmd("start", "Starts the game", startHandler)
	cmd.AddFlag("words", "f", "", "Text file with wordlist", broccli.TypePathFile, broccli.IsExistent|broccli.IsRequired)
	_ = cli.AddCmd("version", "Shows version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}
	os.Exit(cli.Run())
}

func startHandler(c *broccli.CLI) int {
	g := newGame()
	g.readWordsFromFile(c.Flag("words"))
	g.randomizeWords()

	gi := newGameInterface(g)
	gi.setSpeed(500)

	gi.run()
	return 0
}
