package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/nsf/termbox-go"
)

var buttons = [][]string{
	{"C", "⌫", "%", "/"},
	{"7", "8", "9", "x"},
	{"4", "5", "6", "-"},
	{"1", "2", "3", "+"},
	{"0", ".", "+/-", "="},
}

type Calc struct {
	cursorX  int
	cursorY  int
	display  string
	isResult bool
}

func (c *Calc) handleKeyboardChar(ch rune) {
	str := string(ch)

	if str == "*" {
		str = "x"
	}

	for y, row := range buttons {
		for x, btn := range row {
			if btn == str {
				c.cursorX = x // move yellow cursor highlight dynamically
				c.cursorY = y
				c.executeAction(btn)
				return
			}
		}
	}
}

func (c *Calc) executeAction(val string) {
	switch val {
	case "C":
		c.display = ""
		c.isResult = false

	case "⌫":
		if c.isResult {
			c.display = ""
			c.isResult = false
			return
		}
		trimmed := strings.TrimSpace(c.display)
		if len(trimmed) > 0 {
			c.display = trimmed[:len(trimmed)-1]
		}

	case "+/-":
		if c.display == "" || c.display == "0" || c.display == "ERROR" {
			return
		}
		if c.isResult {
			if strings.HasPrefix(c.display, "-") {
				c.display = c.display[1:]
			} else {
				c.display = "-" + c.display
			}
			return
		}

		operators := "/x-+%"
		idx := strings.LastIndexAny(c.display, operators)

		if idx == -1 {
			if strings.HasPrefix(c.display, "-") {
				c.display = c.display[1:]
			} else {
				c.display = "-" + c.display
			}
		} else if idx == len(c.display)-1 {
			return // can not toggle sign on a raw floating operator
		}

		// parantehesis removed
		// else {
		// 	prefix := c.display[:idx+1]
		// 	suffix := c.display[idx+1:]
		// 	if strings.HasPrefix(suffix, "(-") && strings.HasSuffix(suffix, ")") {
		// 		c.display = prefix + suffix[2:len(suffix)-1]
		// 	} else {
		// 		c.display = prefix + "(-" + suffix + ")"
		// 	}
		// }

	case "=":
		if c.evaluate() {
			c.isResult = true
		}

	default:
		if c.isResult {
			if strings.ContainsAny(val, "/x-+%") {
				c.isResult = false
			} else {
				c.display = ""
				c.isResult = false
			}
		}

		// double decimal blocking logic , prevents 12.1.2
		if val == "." {
			operators := "/x-+%"
			idx := strings.LastIndexAny(c.display, operators)

			currentNumber := ""
			if idx == -1 {
				currentNumber = c.display
			} else {
				currentNumber = c.display[idx+1:]
			}

			if strings.Contains(currentNumber, ".") {
				return // ignore input entirely
			}
		}

		// operator swapping logic
		// if user hits '2 +', then hits 'x', it safely swaps them to read '2 x'
		if strings.ContainsAny(val, "/x-+%") {
			trimmed := strings.TrimSpace(c.display)
			if len(trimmed) > 0 {
				lastChar := string(trimmed[len(trimmed)-1])
				if strings.ContainsAny(lastChar, "/x-+%") {
					c.display = trimmed[:len(trimmed)-1] + val
					return
				}
			} else {
				if val == "x" || val == "/" || val == "%" {
					return
				}
			}
		}

		c.display += val
	}
}

func (c *Calc) evaluate() bool {
	if c.display == "" {
		return false
	}

	trimmed := strings.TrimSpace(c.display)

	// blocking incomplete evaluation
	if strings.ContainsAny(string(trimmed[len(trimmed)-1]), "/x-+%") {
		return false
	}

	res, err := evalArbitraryPrecision(c.display)
	if err != nil {
		c.display = "ERROR"
		return true
	} else {
		c.display = res.Text('g', 15)
		return true
	}

}
func evalArbitraryPrecision(expr string) (*big.Float, error) {
	expr = strings.ReplaceAll(expr, "x", "*")
	expr = strings.TrimSpace(expr)

	var tokens []string
	var currentToken strings.Builder

	expectSign := true

	for _, ch := range expr {
		//  pranthesis removed
		// if ch == '(' || ch == ')' || ch == ' ' {
		// 	continue
		// }
		if strings.ContainsRune("+-*/%", ch) {
			if expectSign && (ch == '-' || ch == '+') {
				currentToken.WriteRune(ch)
				expectSign = false
				continue
			}
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(ch))
			expectSign = true
		} else {
			currentToken.WriteRune(ch)
			expectSign = false
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	if len(tokens) == 0 {
		return big.NewFloat(0), nil
	}

	// handle multiplication, division, and modulo first
	for i := 1; i < len(tokens); i += 2 {
		if i+1 >= len(tokens) {
			break
		}
		operator := tokens[i]
		if operator == "*" || operator == "/" || operator == "%" {
			result := new(big.Float).SetPrec(512)
			_, _, err := result.Parse(tokens[i-1], 10)
			if err != nil {
				return nil, err
			}

			nextNum := new(big.Float).SetPrec(512)
			_, _, err = nextNum.Parse(tokens[i+1], 10)
			if err != nil {
				return nil, err
			}

			switch operator {
			case "*":
				result.Mul(result, nextNum)
			case "/":
				if nextNum.Cmp(big.NewFloat(0)) == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				result.Quo(result, nextNum)
			case "%":
				q := new(big.Float).Quo(result, nextNum)
				z, _ := q.Int(nil)
				zFloat := new(big.Float).SetInt(z)
				rem := new(big.Float).Mul(zFloat, nextNum)
				result.Sub(result, rem)
			}

			tokens[i-1] = result.Text('f', -1)
			tokens = append(tokens[:i], tokens[i+2:]...)
			i -= 2
		}
	}

	// handle addition and subtraction then
	total := new(big.Float).SetPrec(512)
	_, _, err := total.Parse(tokens[0], 10)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(tokens); i += 2 {
		if i+1 >= len(tokens) {
			break
		}
		operator := tokens[i]
		nextNumStr := tokens[i+1]

		nextNum := new(big.Float).SetPrec(512)
		_, _, err := nextNum.Parse(nextNumStr, 10)
		if err != nil {
			return nil, err
		}

		switch operator {
		case "+":
			total.Add(total, nextNum)
		case "-":
			total.Sub(total, nextNum)
		}
	}
	return total, nil
}

func printString(x, y int, fg, bg termbox.Attribute, msg string) {
	for i, ch := range msg {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func drawBox(x, y, width, height int, fg, bg termbox.Attribute) {

	// for i := 0; i < width; i++ {
	for i := range width {
		termbox.SetCell(x+i, y, '=', fg, bg)
		termbox.SetCell(x+i, y+height-1, '=', fg, bg)
	}
	// for i := 0; i < height; i++ {
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

	// Label
	drawBox(2, 1, 40, 26, accentColor, termbox.ColorDefault)
	printString(12, 1, termbox.ColorBlack|termbox.AttrBold, accentColor, " Goculator v1.0 ")

	// LCD
	drawBox(4, 3, 36, 5, termbox.ColorGreen, termbox.ColorDefault)

	displayStr := calc.display
	if displayStr == "" {
		displayStr = "0"
	}

	// slice offset truncation
	if len(displayStr) > 30 {
		displayStr = displayStr[len(displayStr)-30:]
	}

	paddedDisplay := fmt.Sprintf("%32s", displayStr)
	printString(6, 5, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, paddedDisplay)

	// buttons grid
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

			// draw blocks
			// for h := 0; h < bHeight; h++ {
			for h := range bHeight {
				for w := range bWidth {
					// for w := 0; w < bWidth; w++ {
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

			startX += bWidth + 2 // offset to draw next btn block
		}
		startY += 3 + 1 // offset down to draw next row

	}

	printString(4, 28, termbox.ColorWhite, termbox.ColorDefault, "▲▼◀▶ Arrows or Type Directly | Esc Exit")
	termbox.Flush()
}

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

		// event listener
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			// ev .Ch catches regular keyboard characters (e.g. '7', '+')
			if ev.Ch != 0 {
				calc.handleKeyboardChar(ev.Ch)
				continue
			}
			//  ev.Key handles btns like arrows, backspace, and escape
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
