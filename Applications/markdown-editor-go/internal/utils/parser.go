package utils

import (
	"bytes"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// MarkdownParser provides enhanced Markdown parsing utilities
type MarkdownParser struct {
	extensions parser.Extensions
	htmlFlags  html.Flags
}

// NewMarkdownParser creates a new parser with default settings
func NewMarkdownParser() *MarkdownParser {
	// Default parser extensions
	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoEmptyLineBeforeBlock |
		parser.Tables |
		parser.FencedCode |
		parser.Strikethrough

	// Default HTML renderer flags
	htmlFlags := html.CommonFlags |
		html.HrefTargetBlank |
		html.CompletePage |
		html.FootnoteReturnLinks

	return &MarkdownParser{
		extensions: extensions,
		htmlFlags:  htmlFlags,
	}
}

// MarkdownToHTML converts markdown text to HTML
func (p *MarkdownParser) MarkdownToHTML(md string) string {
	// Create a new parser with the defined extensions
	parser := parser.NewWithExtensions(p.extensions)

	// Parse the markdown input
	node := parser.Parse([]byte(md))

	// Create HTML renderer with our flags
	renderer := html.NewRenderer(html.RendererOptions{
		Flags: p.htmlFlags,
	})

	// Generate HTML output
	html := markdown.Render(node, renderer)

	return string(html)
}

// ExtractTOC extracts a table of contents from markdown
func (p *MarkdownParser) ExtractTOC(md string) string {
	parser := parser.NewWithExtensions(p.extensions)
	node := parser.Parse([]byte(md))

	var headers []string
	ast.WalkFunc(node, func(node ast.Node, entering bool) ast.WalkStatus {
		if heading, ok := node.(*ast.Heading); ok && entering {
			level := heading.Level
			text := renderHeadingText(heading)

			// Create a slug for the heading
			slug := slugify(text)

			// Create list item with proper indentation based on level
			indent := strings.Repeat("  ", level-1)
			headers = append(headers, indent+"- ["+text+"](#"+slug+")")
		}
		return ast.GoToNext
	})

	return strings.Join(headers, "\n")
}

// ExtractHeadings extracts all headings from markdown
func (p *MarkdownParser) ExtractHeadings(md string) []map[string]interface{} {
	parser := parser.NewWithExtensions(p.extensions)
	node := parser.Parse([]byte(md))

	var headings []map[string]interface{}

	ast.WalkFunc(node, func(node ast.Node, entering bool) ast.WalkStatus {
		if heading, ok := node.(*ast.Heading); ok && entering {
			text := renderHeadingText(heading)
			slug := slugify(text)

			headings = append(headings, map[string]interface{}{
				"level": heading.Level,
				"text":  text,
				"slug":  slug,
			})
		}
		return ast.GoToNext
	})

	return headings
}

// WordCount counts words in markdown text (ignoring code blocks and metadata)
func (p *MarkdownParser) WordCount(md string) int {
	// Remove code blocks
	codeBlockRegex := "```[\\s\\S]*?```"
	md = strings.ReplaceAll(md, codeBlockRegex, "")

	// Remove inline code
	inlineCodeRegex := "`[^`]*`"
	md = strings.ReplaceAll(md, inlineCodeRegex, "")

	// Remove HTML tags
	htmlTagRegex := "<[^>]*>"
	md = strings.ReplaceAll(md, htmlTagRegex, "")

	// Remove URLs
	urlRegex := "\\[([^\\[]+)\\]\\([^\\)]+\\)"
	md = strings.ReplaceAll(md, urlRegex, "$1")

	// Count words by splitting on whitespace
	words := strings.Fields(md)
	return len(words)
}

// Helper function to render heading text
func renderHeadingText(heading *ast.Heading) string {
	var buf bytes.Buffer
	for _, child := range heading.Children {
		if text, ok := child.(*ast.Text); ok {
			buf.Write(text.Literal)
		}
	}
	return buf.String()
}

// Slugify creates a URL-friendly slug from a string
func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove special characters
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)

	return s
}
