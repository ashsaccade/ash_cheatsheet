package render

import (
	"fmt"
	"html"
	"strings"

	"ash_cheatsheet/internal/grammar"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
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
			if err := highlightGo(out, block.Str); err != nil {
				return "", fmt.Errorf("failed to highlight code: %v", err)
			}
			out.WriteString(`</div>`)
		}
	}
	return out.String(), nil
}

func highlightGo(out *strings.Builder, source string) error {
	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		return err
	}

	formatter := chromahtml.New(chromahtml.WithClasses(false))
	return formatter.Format(out, style, iterator)
}

func processText(in string) string {
	in = html.EscapeString(in)
	return strings.ReplaceAll(in, "\n", "<br>")
}
