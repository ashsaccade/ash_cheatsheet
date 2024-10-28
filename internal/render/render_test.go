package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Render_SimpleText_ShoulbBeRendredAsIs(t *testing.T) {
	out := Render("hello world!")
	assert.Equal(t, "hello world!", out)
}

func Test_Render_HtmlTags_ShouldBeEscaped(t *testing.T) {
	out := Render("<div></div>")
	assert.Equal(t, "&lt;div&gt;&lt;/div&gt;", out)
}

func Test_Render_BoldMarkdown_ShouldBeWrappedInBTag(t *testing.T) {
	out := Render("Hello **world**!")
	// out := Render("Hello **world**! `code` ```go``` ")
	assert.Equal(t, "Hello <b>world</b>!", out)
}

func Test_Render_Incorrect(t *testing.T) {
	out := Render("Hello **wo*rld**!")
	// out := Render("Hello **world**! `code` ```go``` ")
	assert.Equal(t, "Hello <b>wo*rld</b>!", out)
}
