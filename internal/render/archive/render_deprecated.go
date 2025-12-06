package archive

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

var (
	singleCodeRe = regexp.MustCompile("(?U)`([^`]+)`")
	multiCodeRe  = regexp.MustCompile("(?Ums)```(.+)```")
)

var boldRe = regexp.MustCompile(`(?Um)\*\*(.+)\*\*`)

func renderBold(in string) string {
	return boldRe.ReplaceAllString(in, `<b>$1</b>`)
}

func Render(in string) string {
	// in = html.EscapeString(in)
	res := in

	res = renderBold(res)
	res = renderCode(res)

	// res = strings.ReplaceAll(res, "\n", "<br>")

	return res
}

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
