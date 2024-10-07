# ortotris

This project is a small, terminal-based orthography game designed for my kids. It is inspired by the classic DOS game Ortotris, released in 1992, which was similar to Tetris but focused on improving spelling skills. The game runs in the terminal and follows a similar concept to help players with orthography.

See screenshot below:

![Ortotris](screenshot.png)

### Running

To play the game just run:

    go run *.go zagraj -f words-u-o.txt

### Instructions
Words descend from the top of the screen, similar to Tetris, but with one or two missing letters, indicated by an underscore (_). Use the left and right arrow keys to select one of the available letters before the word reaches the bottom. If an incorrect letter is chosen, the word will remain at the bottom. You can also press the down arrow to drop the word immediately.

