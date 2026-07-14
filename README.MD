# Goculator v1.0

A terminal-based calculator built in Go using `termbox-go` for TUI rendering
and `math/big` for arbitrary precision arithmetic.

## Features

- Basic arithmetic: addition, subtraction, multiplication, division, modulo
- Keyboard input: type numbers and operators directly, or navigate with arrow keys
- Sign toggle (`+/-`): negate the current number
- Cursor navigation with visual highlight (yellow)

![Project Image](https://private-user-images.githubusercontent.com/47816410/621541109-0f539916-adf7-47ee-ac6f-e19fcfcb582d.png?jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3ODQwNDA1MDgsIm5iZiI6MTc4NDA0MDIwOCwicGF0aCI6Ii80NzgxNjQxMC82MjE1NDExMDktMGY1Mzk5MTYtYWRmNy00N2VlLWFjNmYtZTE5ZmNmY2I1ODJkLnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFWQ09EWUxTQTUzUFFLNFpBJTJGMjAyNjA3MTQlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjYwNzE0VDE0NDMyOFomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPTMzNTJkOTBhZGFkOTc5ZTQ4NTU4MDM2ODY1ZTA1OWRiZjkxZmM4ZDFjZmQ2ZjI2ZGU3MmRjMTVmYWQ3NzExMDMmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JnJlc3BvbnNlLWNvbnRlbnQtdHlwZT1pbWFnZSUyRnBuZyJ9.q0PAtwiZl4yPvDtuhggbJ-Qn7ewXxBcQxnPLHqCl3Yc)



## Install & Run

```bash
# Install dependencies
go mod tidy

# Run
go run .

# Or build a binary
go build .
```


# Controls

| Key   |      Action      |
|----------|:-------------:|
| 0-9 , . |  Input digits and decimal |
| +, -, x, /, % |    Operators   |
| = or Enter | Evaluate |
| C | 	Clear all |
| Backspace | Delete last character |
| +/- | Toggle sign of current number|
| Arrow keys | 	Navigate button grid |
| Esc / Ctrl+C | Exit |



