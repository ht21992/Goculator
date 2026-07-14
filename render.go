package main

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

func printString(x, y int, fg, bg termbox.Attribute, msg string) {
	for i, ch := range msg {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func drawBox(x, y, width, height int, fg, bg termbox.Attribute) {
	for i := range width {
		termbox.SetCell(x+i, y, '=', fg, bg)
		termbox.SetCell(x+i, y+height-1, '=', fg, bg)
	}
	for i := range height {
		termbox.SetCell(x, y+i, '║', fg, bg)
		termbox.SetCell(x+width-1, y+i, '║', fg, bg)
	}

	termbox.SetCell(x, y, '╔', fg, bg)
	termbox.SetCell(x+width-1, y, '╗', fg, bg)
	termbox.SetCell(x, y+height-1, '╚', fg, bg)
	termbox.SetCell(x+width-1, y+height-1, '╝', fg, bg)
}

func draw(calc Calc) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	accentColor := termbox.ColorBlue
	textColor := termbox.ColorWhite

	drawBox(2, 1, 40, 26, accentColor, termbox.ColorDefault)
	printString(12, 1, termbox.ColorBlack|termbox.AttrBold, accentColor, " Goculator v1.0 ")

	drawBox(4, 3, 36, 5, termbox.ColorGreen, termbox.ColorDefault)

	displayStr := calc.display
	if displayStr == "" {
		displayStr = "0"
	}

	if len(displayStr) > 30 {
		displayStr = displayStr[len(displayStr)-30:]
	}

	paddedDisplay := fmt.Sprintf("%32s", displayStr)
	printString(6, 5, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, paddedDisplay)

	startY := 9
	for y, row := range buttons {
		startX := 4
		for x, btn := range row {
			bWidth := 7
			bHeight := 3

			fg := textColor
			bg := termbox.ColorDefault
			isSelected := (calc.cursorX == x && calc.cursorY == y)

			if isSelected {
				fg = termbox.ColorBlack | termbox.AttrBold
				bg = termbox.ColorYellow
			} else if strings.ContainsAny(btn, "/x-+%=") {
				fg = termbox.ColorMagenta | termbox.AttrBold
			} else if btn == "C" || btn == "⌫" {
				fg = termbox.ColorRed | termbox.AttrBold
			}

			for h := range bHeight {
				for w := range bWidth {
					ch := ' '
					if h == 0 && w == 0 {
						ch = '┌'
					}
					if h == 0 && w == bWidth-1 {
						ch = '┐'
					}
					if h == bHeight-1 && w == 0 {
						ch = '└'
					}
					if h == bHeight-1 && w == bWidth-1 {
						ch = '┘'
					}

					if isSelected {
						ch = ' '
					}

					termbox.SetCell(startX+w, startY+h, ch, fg, bg)
				}
			}

			labelX := startX + ((bWidth - len(btn)) / 2)
			labelY := startY + 1
			printString(labelX, labelY, fg, bg, btn)

			startX += bWidth + 2
		}
		startY += 3 + 1
	}

	printString(4, 28, termbox.ColorWhite, termbox.ColorDefault, "▲▼◀▶ Arrows or Type Directly | Esc Exit")
	termbox.Flush()
}