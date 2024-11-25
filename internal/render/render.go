package render

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

type State int

const (
	StateWaitingForFirstAsterisk State = iota
	StateWaitingForSecondAsterisk
	StateBold
	StateWaitingSecondClosing
)

var singleCodeRe = regexp.MustCompile("(?U)`([^`]+)`")
var multiCodeRe = regexp.MustCompile("(?Ums)```(.+)```")

func renderCode(in string) string {
	res := in

	res = multiCodeRe.ReplaceAllStringFunc(res, func(match string) string {
		match = strings.Trim(match, "```")
		match = strings.TrimSpace(match)

		buf := bytes.NewBuffer(nil)
		quick.Highlight(buf, match, "go", "html", "xcode")
		match = buf.String()

		return `<div class="multicode">` + match + "</div>"
	})

	// res =singleCodeRe.ReplaceAllString(in, `<span class="mycode">$1</span>`)

	return res
}

var boldRe = regexp.MustCompile(`(?Um)\*\*(.+)\*\*`)

func renderBold(in string) string {
	return boldRe.ReplaceAllString(in, `<b>$1</b>`)
}

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

func Render(in string) string {
	// in = html.EscapeString(in)
	res := in

	res = renderBold(res)
	res = renderCode(res)

	// res = strings.ReplaceAll(res, "\n", "<br>")

	return res
}
