package parsers_test

import (
	"testing"

	"github.com/schoonology/diplomat/parsers"
	"github.com/stretchr/testify/assert"
)

func TestTextInput(t *testing.T) {
	parser := parsers.GetParserFromFileName("test.txt")

	switch parser.(type) {
	case *parsers.PlainTextParser:
	default:
		t.Fail()
	}
}

func TestMdInput(t *testing.T) {
	parser := parsers.GetParserFromFileName("test.md")

	switch v := parser.(type) {
	case *parsers.Markdown:
	default:
		assert.Failf(t, "Incorrect Type", "Received type %T.", v)
	}
}

func TestMarkdownInput(t *testing.T) {
	parser := parsers.GetParserFromFileName("test.markdown")

	switch v := parser.(type) {
	case *parsers.Markdown:
	default:
		assert.Failf(t, "Incorrect Type", "Received type %T.", v)
	}
}

func TestInvalidInput(t *testing.T) {
	parser := parsers.GetParserFromFileName("test.horribad")

	assert.Nil(t, parser)
}
