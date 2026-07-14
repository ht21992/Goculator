package main

import (
	"fmt"
	"math/big"
	"strings"
)

func evalArbitraryPrecision(expr string) (*big.Float, error) {
	expr = strings.ReplaceAll(expr, "x", "*")
	expr = strings.TrimSpace(expr)

	var tokens []string
	var currentToken strings.Builder

	expectSign := true

	for _, ch := range expr {
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

	// multiplication, division, modulo first
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

	// addition and subtraction
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