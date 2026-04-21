package render

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Render_SimpleText_ShoulbBeRendredAsIs(t *testing.T) {
	out, err := Render("hello world!")
	assert.NoError(t, err)
	assert.Equal(t, "hello world!", out)
}

func Test_Render_HtmlTags_ShouldBeEscaped(t *testing.T) {
	out, err := Render("<div></div>")
	assert.NoError(t, err)
	assert.Equal(t, "&lt;div&gt;&lt;/div&gt;", out)
}

func Test_Render_BoldMarkdown_ShouldBeWrappedInBTag(t *testing.T) {
	out, err := Render("Hello *world* `world`!")
	assert.NoError(t, err)
	assert.Equal(t, `Hello <b>world</b> <span class="inlinecode">world</span>!`, out)
}

func Test_OneBacktick(t *testing.T) {
	out, err := Render("a `b`")
	assert.NoError(t, err)
	assert.Equal(t, "a <span class=\"inlinecode\">b</span>", out)
}

func Test_Render_BigCodeBlock_ShouldKeepChromaStylesInsideBlock(t *testing.T) {
	out, err := Render("```fmt.Println(\"<hello>\")```")
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(out, `<div class="multicode"><pre`))
	assert.Contains(t, out, "&lt;hello&gt;")
	assert.NotContains(t, out, "<style")
	assert.NotContains(t, out, "<body")
	assert.NotContains(t, out, "<html")
}
