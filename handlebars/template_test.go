package handlebars

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Render(t *testing.T) {
	r := require.New(t)

	ctx := NewContext()
	ctx.Set("name", "Tim")
	s, err := Render("{{name}}", ctx)
	r.NoError(err)
	r.Equal("Tim", s)
}

func Test_Render_with_Content(t *testing.T) {
	r := require.New(t)

	ctx := NewContext()
	ctx.Set("name", "Tim")
	s, err := Render("<p>{{name}}</p>", ctx)
	r.NoError(err)
	r.Equal("<p>Tim</p>", s)
}

func Test_Render_Unknown_Value(t *testing.T) {
	r := require.New(t)

	ctx := NewContext()
	_, err := Render("<p>{{name}}</p>", ctx)
	r.Error(err)
	r.Equal("could not find value for name [line 1:3]", err.Error())
}

func Test_Render_with_String(t *testing.T) {
	r := require.New(t)

	ctx := NewContext()
	s, err := Render(`<p>{{"Tim"}}</p>`, ctx)
	r.NoError(err)
	r.Equal("<p>Tim</p>", s)
}

func Test_Render_with_Math(t *testing.T) {
	r := require.New(t)

	ctx := NewContext()
	_, err := Render(`<p>{{2 + 1}}</p>`, ctx)
	r.Error(err)
}

func Test_Render_with_Comments(t *testing.T) {
	r := require.New(t)
	ctx := NewContext()
	s, err := Render(`<p><!-- comment --></p>`, ctx)
	r.NoError(err)
	r.Equal("<p><!-- comment --></p>", s)
}

func Test_Render_with_Func(t *testing.T) {
	r := require.New(t)
	ctx := NewContext()
	ctx.Set("user", user{First: "Mark", Last: "Bates"})
	s, err := Render("{{user.FullName}}", ctx)
	r.NoError(err)
	r.Equal("Mark Bates", s)
}

func Test_Render_Array(t *testing.T) {
	r := require.New(t)

	ctx := NewContext()
	ctx.Set("names", []string{"mark", "bates"})
	s, err := Render("{{names}}", ctx)
	r.NoError(err)
	r.Equal("mark bates", s)
}

type user struct {
	First string
	Last  string
}

func (u user) FullName() string {
	return u.First + " " + u.Last
}
