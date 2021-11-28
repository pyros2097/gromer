package gromer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func writeIndent(w io.Writer, indent int) {
	for i := 0; i < indent*2; i++ {
		w.Write([]byte(" "))
	}
}

func mergeAttributes(parent *Element, uis ...interface{}) *Element {
	elems := []*Element{}
	for _, v := range uis {
		switch c := v.(type) {
		case Attribute:
			parent.setAttr(c.Key, c.Value)
		case M:
			for k, v := range c {
				if a, ok := v.(string); ok {
					parent.setAttr(k, a)
				} else {
					// store some server state on the component ctx
					// parent.setAttr(k, a)
				}

			}
		case *Element:
			elems = append(elems, c)
		case nil:
			// dont need to add nil items
		default:
			panic(fmt.Sprintf("unknown type in render %+v", v))
		}
	}
	if !parent.selfClosing {
		parent.children = elems
	}
	return parent
}

type HtmlPage struct {
	ctx         context.Context
	classLookup map[string]bool
	css         *bytes.Buffer
	js          *bytes.Buffer
	Head        *Element
	Body        *Element
}

func (p *HtmlPage) computeCss(elems []*Element) {
	for _, el := range elems {
		if v, ok := el.attrs["class"]; ok {
			classes := strings.Split(v, " ")
			for _, c := range classes {
				if s, ok := twClassLookup[c]; ok {
					if _, ok2 := p.classLookup[c]; !ok2 {
						p.classLookup[c] = true
						p.css.WriteString("      ." + c + " { " + s + " }\n")
					}
				}
			}
		}
		if len(el.children) > 0 {
			p.computeCss(el.children)
		}
	}
}

func (p *HtmlPage) WriteHtml(w io.Writer) {
	w.Write([]byte("<!DOCTYPE html>\n"))
	w.Write([]byte("<html lang=\"en\">\n"))
	p.computeCss(p.Body.children)
	p.Head.children = append(p.Head.children, StyleTag(Text(normalizeStyles+p.css.String())))
	p.Head.writeHtmlIndent(w, 1)
	// p.Body.children = append(p.Body.children, Script(Text(fmt.Sprintf(`
	// 		document.addEventListener('alpine:init', () => {%s
	// 		});
	// `, p.js.String()))))
	p.Body.writeHtmlIndent(w, 1)
	w.Write([]byte("\n</html>"))
}

func Html(h *Element, b *Element) HtmlPage {
	return HtmlPage{
		// ctx:         context.WithValue(ctx, "state", &StateContext{}),
		classLookup: map[string]bool{},
		js:          bytes.NewBuffer(nil),
		css:         bytes.NewBuffer(nil),
		Head:        h,
		Body:        b,
	}
}

func Head(elems ...*Element) *Element {
	basic := []*Element{
		{tag: "meta", selfClosing: true, attrs: map[string]string{"charset": "UTF-8"}},
		{tag: "meta", selfClosing: true, attrs: map[string]string{"http-equiv": "Content-Type", "content": "text/html;charset=utf-8"}},
		{tag: "meta", selfClosing: true, attrs: map[string]string{"http-equiv": "encoding", "content": "utf-8"}},
	}
	return &Element{tag: "head", children: append(basic, elems...)}
}

func Body(elems ...*Element) *Element {
	return &Element{tag: "body", children: elems}
}

type Element struct {
	tag         string
	attrs       map[string]string
	children    []*Element
	selfClosing bool
	text        string
}

func NewElement(tag string, selfClosing bool, uis ...interface{}) *Element {
	return mergeAttributes(&Element{tag: tag, selfClosing: selfClosing}, uis...)
}

func (e *Element) setAttr(k string, v string) {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}

	switch k {
	case "style", "allow":
		s := e.attrs[k] + v + ";"
		e.attrs[k] = s
		return

	case "class":
		s := e.attrs[k]
		if s != "" {
			s += " "
		}
		s += v
		e.attrs[k] = s
		return
	}
	if v == "false" {
		delete(e.attrs, k)
		return
	} else if v == "true" {
		e.attrs[k] = ""
	} else {
		e.attrs[k] = v
	}
}

func (e *Element) writeHtmlIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	if e.tag == "text" {
		writeIndent(w, indent)
		w.Write([]byte(e.text))
		return
	}
	w.Write([]byte("<"))
	w.Write([]byte(e.tag))

	for k, v := range e.attrs {
		w.Write([]byte(" "))
		w.Write([]byte(k))

		if v != "" {
			w.Write([]byte(`="`))
			w.Write([]byte(v))
			w.Write([]byte(`"`))
		}
	}

	w.Write([]byte(">"))

	if e.selfClosing {
		return
	}

	for _, c := range e.children {
		if len(e.children) > 1 {
			w.Write([]byte("\n"))
		}
		if c != nil {
			c.writeHtmlIndent(w, indent+1)
		}
	}

	if len(e.children) != 0 {
		// w.Write([]byte("\n"))
		writeIndent(w, indent)
	}

	w.Write([]byte("</"))
	w.Write([]byte(e.tag))
	w.Write([]byte(">\n"))
}

func Title(v string) *Element {
	return &Element{
		tag:      "title",
		children: []*Element{Text(v)},
	}
}
func Text(v string) *Element {
	return &Element{
		tag:  "text",
		text: v,
	}
}

func Meta(name, content string) *Element {
	return &Element{
		tag:         "meta",
		selfClosing: true,
		attrs: map[string]string{
			"name":    name,
			"content": content,
		},
	}
}

func Link(rel, href string) *Element {
	return &Element{
		tag:         "link",
		selfClosing: true,
		attrs: map[string]string{
			"rel":  rel,
			"href": href,
		},
	}
}

func Script(uis ...interface{}) *Element {
	return NewElement("script", false, uis...)
}

func StyleTag(uis ...interface{}) *Element {
	return NewElement("style", false, uis...)
}

func Div(uis ...interface{}) *Element {
	return NewElement("div", false, uis...)
}

func A(uis ...interface{}) *Element {
	return NewElement("a", false, uis...)
}

func P(uis ...interface{}) *Element {
	return NewElement("p", false, uis...)
}

func H1(uis ...interface{}) *Element {
	return NewElement("h1", false, uis...)
}
func H2(uis ...interface{}) *Element {
	return NewElement("h2", false, uis...)
}
func H3(uis ...interface{}) *Element {
	return NewElement("h3", false, uis...)
}
func H4(uis ...interface{}) *Element {
	return NewElement("h4", false, uis...)
}
func H5(uis ...interface{}) *Element {
	return NewElement("h5", false, uis...)
}
func H6(uis ...interface{}) *Element {
	return NewElement("h6", false, uis...)
}

func Span(uis ...interface{}) *Element {
	return NewElement("span", false, uis...)
}

func Input(uis ...interface{}) *Element {
	return NewElement("input", true, uis...)
}

func Image(uis ...interface{}) *Element {
	return NewElement("image", true, uis...)
}

func Button(uis ...interface{}) *Element {
	return NewElement("button", false, uis...)
}

func Svg(uis ...interface{}) *Element {
	return NewElement("svg", false, uis...)
}

func SvgText(uis ...interface{}) *Element {
	return NewElement("text", false, uis...)
}

func Ul(uis ...interface{}) *Element {
	return NewElement("ul", false, uis...)
}

func Li(uis ...interface{}) *Element {
	return NewElement("li", false, uis...)
}

type Attribute struct {
	Key   string
	Value string
}

func Attr(k, v string) Attribute {
	return Attribute{k, v}
}

func OnClick(v string) Attribute {
	return Attribute{"@click", v}
}

func ID(v string) Attribute {
	return Attribute{"id", v}
}

func Style(v string) Attribute {
	return Attribute{"style", v}
}

func Accept(v string) Attribute {
	return Attribute{"accept", v}
}

func AutoComplete(v bool) Attribute {
	return Attribute{"autocomplete", strconv.FormatBool(v)}
}

func Checked(v bool) Attribute {
	return Attribute{"checked", strconv.FormatBool(v)}
}

func Disabled(v bool) Attribute {
	return Attribute{"disabled", strconv.FormatBool(v)}
}

func Name(v string) Attribute {
	return Attribute{"name", v}
}

func Type(v string) Attribute {
	return Attribute{"type", v}
}

func Value(v string) Attribute {
	return Attribute{"value", v}
}

func Placeholder(v string) Attribute {
	return Attribute{"placeholder", v}
}

func Src(v string) Attribute {
	return Attribute{"src", v}
}

func Defer() Attribute {
	return Attribute{"defer", "true"}
}

func ViewBox(v string) Attribute {
	return Attribute{"viewBox", v}
}

func X(v string) Attribute {
	return Attribute{"x", v}
}

func Y(v string) Attribute {
	return Attribute{"y", v}
}

func Href(v string) Attribute {
	return Attribute{"href", v}
}

func Target(v string) Attribute {
	return Attribute{"target", v}
}

func Rel(v string) Attribute {
	return Attribute{"rel", v}
}

func Css(v string) Attribute {
	return Attribute{"class", v}
}

func XData(v string) Attribute {
	return Attribute{"x-data", v}
}

func XText(v string) Attribute {
	return Attribute{"x-text", v}
}
