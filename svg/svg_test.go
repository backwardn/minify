package svg // import "github.com/tdewolff/minify/svg"

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/parse/svg"
	"github.com/tdewolff/parse/xml"
)

func TestSVG(t *testing.T) {
	var svgTests = []struct {
		svg      string
		expected string
	}{
		{`<!-- comment -->`, ``},
		{`<!DOCTYPE foo SYSTEM "Foo.dtd">`, ``},
		{`<?xml version="1.0" ?>`, ``},
		{`<style> <![CDATA[ x ]]> </style>`, `<style>x</style>`},
		{`<svg version="1.0"></svg>`, ``},
		{`<path x=" a "/>`, `<path x="a"/>`},
		{"<path x=\" a \n b \"/>", `<path x="a b"/>`},
		{`<path x="5.0px" y="0%"/>`, `<path x="5" y="0"/>`},
		{`<svg viewBox="5.0px 5px 240 0.10"><path/></svg>`, `<svg viewBox="5 5 240 .1"><path/></svg>`},
		{`<path d="M 100 100 L 300 100 L 200 100 z"/>`, `<path d="M100 100l200 0 100 0z"/>`},
		{`<path d="M100 -100M200 300z"/>`, `<path d="M100-100 200 300z"/>`},
		{`<path d="M0.5 0.6 M -100 0.5z"/>`, `<path d="M.5.6-100 .5z"/>`},
		{`<path d="M01.0 0.6 z"/>`, `<path d="M1 .6z"/>`},
		{`<path d="M20 20l-10-10z"/>`, `<path d="M20 20L10 10z"/>`},
		{`<?xml version="1.0" encoding="utf-8"?>`, ``},
		{`<svg viewbox="0 0 16 16"><path/></svg>`, `<svg viewbox="0 0 16 16"><path/></svg>`},
		{`<g></g>`, ``},
		{`<path fill="#ffffff"/>`, `<path fill="#fff"/>`},
		{`<line x1="5" y1="10" x2="20" y2="40"/>`, `<path d="M5 10L20 40z"/>`},
		{`<rect x="5" y="10" width="20" height="40"/>`, `<path d="M5 10h20v40H5z"/>`},
		{`<polygon points="1,2 3,4"/>`, `<path d="M1 2L3 4z"/>`},
		{`<polyline points="1,2 3,4"/>`, `<path d="M1 2L3 4"/>`},
		{`<svg contentStyleType="text/json ; charset=iso-8859-1"><style>{a : true}</style></svg>`, `<svg contentStyleType="text/json;charset=iso-8859-1"><style>{a : true}</style></svg>`},
		{`<metadata><dc:title /></metadata>`, ``},

		// from SVGO
		{`<!DOCTYPE bla><?xml?><!-- comment --><metadata/>`, ``},
	}

	m := minify.New()
	for _, tt := range svgTests {
		b := &bytes.Buffer{}
		assert.Nil(t, Minify(m, "image/svg+xml", b, bytes.NewBufferString(tt.svg)), "Minify must not return error in "+tt.svg)
		assert.Equal(t, tt.expected, b.String(), "Minify must give expected result in "+tt.svg)
	}
}

func TestGetAttribute(t *testing.T) {
	r := bytes.NewBufferString(`<rect x="0" y="1" width="2" height="3" rx="4" ry="5"/>`)
	attrTokenBuffer := make([]*svg.Token, 0, maxAttrLookup)
	l := xml.NewLexer(r)
	tb := svg.NewTokenBuffer(l)
	tb.Shift()
	getAttributes(&attrTokenBuffer, tb, svg.X, svg.Y, svg.Width, svg.Height, svg.Rx, svg.Ry)
	for i := 0; i < 6; i++ {
		assert.NotNil(t, attrTokenBuffer[i], "Attr is nil")
		j, _ := strconv.ParseInt(string(attrTokenBuffer[i].AttrVal), 10, 32)
		assert.Equal(t, i, int(j), "Attr data is bad")
	}
}

////////////////////////////////////////////////////////////////

func ExampleMinify() {
	m := minify.New()
	m.AddFunc("image/svg+xml", Minify)
	m.AddFunc("text/css", css.Minify)

	if err := m.Minify("image/svg+xml", os.Stdout, os.Stdin); err != nil {
		fmt.Println("minify.Minify:", err)
	}
}

////////////////////////////////////////////////////////////////

func BenchmarkGetAttributes(b *testing.B) {
	r := bytes.NewBufferString(`<rect x="0" y="1" width="2" height="3" rx="4" ry="5"/>`)
	attrTokenBuffer := make([]*svg.Token, 0, maxAttrLookup)
	l := xml.NewLexer(r)
	tb := svg.NewTokenBuffer(l)
	tb.Shift()
	tb.Peek(6)
	for i := 0; i < b.N; i++ {
		getAttributes(&attrTokenBuffer, tb, svg.X, svg.Y, svg.Width, svg.Height, svg.Rx, svg.Ry)
	}
}

func BenchmarkShortenPathData(b *testing.B) {
	pathDataBuffer := &PathData{}
	r := []byte("M8.64,223.948c0,0,143.468,3.431,185.777-181.808c2.673-11.702-1.23-20.154,1.316-33.146h16.287c0,0-3.14,17.248,1.095,30.848c21.392,68.692-4.179,242.343-204.227,196.59L8.64,223.948z")
	for i := 0; i < b.N; i++ {
		ShortenPathData(r, pathDataBuffer)
	}
}
