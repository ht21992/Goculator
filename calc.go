package main

import (
	"strings"
)

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
				c.cursorX = x
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
			return
		}

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
				return
			}
		}

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

	if strings.ContainsAny(string(trimmed[len(trimmed)-1]), "/x-+%") {
		return false
	}

	res, err := evalArbitraryPrecision(c.display)
	if err != nil {
		c.display = "ERROR"
		return true
	}

	c.display = res.Text('g', 15)
	return true
}