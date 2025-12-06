package archive

import "strings"

type State int

const (
	StateWaitingForFirstAsterisk State = iota
	StateWaitingForSecondAsterisk
	StateBold
	StateWaitingSecondClosing
)

// renderBold2 на основе стейт машины
func renderBold2(in string) string {
	res := strings.Builder{}
	state := StateWaitingForFirstAsterisk
	for _, letter := range in {
		switch state {
		case StateWaitingForFirstAsterisk:
			if letter == '*' {
				state = StateWaitingForSecondAsterisk
			} else {
				res.WriteRune(letter)
			}
		case StateWaitingForSecondAsterisk:
			if letter == '*' {
				state = StateBold
				res.WriteString("<b>")
			} else {
				state = StateWaitingForFirstAsterisk
				res.WriteByte('*')
				res.WriteRune(letter)
			}
		case StateBold:
			if letter == '*' {
				state = StateWaitingSecondClosing
			} else {
				res.WriteRune(letter)
			}
		case StateWaitingSecondClosing:
			if letter == '*' {
				state = StateWaitingForFirstAsterisk
				res.WriteString("</b>")
			} else {
				state = StateBold
				res.WriteByte('*')
				res.WriteRune(letter)
			}
		}
	}
	return res.String()
}
