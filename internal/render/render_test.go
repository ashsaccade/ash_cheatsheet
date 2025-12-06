package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
