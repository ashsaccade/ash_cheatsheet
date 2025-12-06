package render

import (
	"fmt"
	"html"

	"strings"

	"ash_cheatsheet/internal/grammar"

	"github.com/alecthomas/chroma/v2/quick"
)

func Render(in string) (string, error) {
	md := &grammar.AshMd{Buffer: in}
	if err := md.Init(); err != nil {
		return "", fmt.Errorf("failed to init md: %v", err)
	}

	if err := md.Parse(); err != nil {
		return "", fmt.Errorf("failed to parse md: %v", err)
	}
	// md.PrettyPrintSyntaxTree(in)

	out := &strings.Builder{}
	for _, block := range md.ParseAST(in) {
		switch block.Type {
		case grammar.BlockTypeText:
			out.WriteString(processText(block.Str))
		case grammar.BlockTypeBold:
			out.WriteString("<b>")
			out.WriteString(processText(block.Str))
			out.WriteString("</b>")
		case grammar.BlockTypeCode:
			out.WriteString(`<span class="inlinecode">`)
			out.WriteString(html.EscapeString(block.Str))
			out.WriteString("</span>")
		case grammar.BlockTypeBigCode:
			out.WriteString(`<div class="multicode">`)
			quick.Highlight(out, block.Str, "go", "html", "monokai") // dark theme
			//quick.Highlight(out, block.Str, "go", "html", "xcode") // light theme
			out.WriteString(`</div>`)
		}
	}
	return out.String(), nil
}

func processText(in string) string {
	in = html.EscapeString(in)
	return strings.ReplaceAll(in, "\n", "<br>")
}
