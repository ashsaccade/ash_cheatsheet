package main

import (
	"ash_cheatsheet/internal/grammar"
	"fmt"
	"github.com/alecthomas/chroma/v2/quick"
	"html"
	"strings"
)

func main() {
	input := "text *bold* some other text *dfsdfsdf* my pretty `code` ```big code \n block```"
	ashMd := &grammar.AshMd{Buffer: input}
	err := ashMd.Init()
	if err != nil {
		panic(err)
	}

	err = ashMd.Parse()
	if err != nil {
		panic(err)
	}
	ashMd.PrettyPrintSyntaxTree(input)

	blocks := ashMd.ParseAST(input)

	out := &strings.Builder{}
	for _, block := range blocks {
		switch block.Type {
		case grammar.BlockTypeText:
			out.WriteString(html.EscapeString(block.Str))
		case grammar.BlockTypeBold:
			out.WriteString("<b>")
			out.WriteString(html.EscapeString(block.Str))
			out.WriteString("</b>")
		case grammar.BlockTypeCode:
			out.WriteString(`<span class="inline-code">`)
			out.WriteString(html.EscapeString(block.Str))
			out.WriteString("</span>")
		case grammar.BlockTypeBigCode:
			out.WriteString(`<pre><code>`)
			quick.Highlight(out, block.Str, "go", "html", "xcode")

			out.WriteString(html.EscapeString(block.Str))
			out.WriteString(`</code></pre>`)
		}
	}
	fmt.Println(out.String())
}
