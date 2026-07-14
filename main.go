package main

import (
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)
	calc := Calc{cursorX: 0, cursorY: 0, display: ""}

	for {
		draw(calc)

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch != 0 {
				calc.handleKeyboardChar(ev.Ch)
				continue
			}
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return
			case termbox.KeyArrowUp:
				if calc.cursorY > 0 {
					calc.cursorY--
				}
			case termbox.KeyArrowDown:
				if calc.cursorY < len(buttons)-1 {
					calc.cursorY++
				}
			case termbox.KeyArrowLeft:
				if calc.cursorX > 0 {
					calc.cursorX--
				}
			case termbox.KeyArrowRight:
				if calc.cursorX < len(buttons[calc.cursorY])-1 {
					calc.cursorX++
				}
			case termbox.KeyEnter:
				calc.executeAction(buttons[calc.cursorY][calc.cursorX])
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				calc.executeAction("⌫")
			}
		}
	}
}