package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCrfppParser_newParser(t *testing.T) {
	parser, err := NewParser([]byte(`{"parsertype":"crfppParser", "modelfilepath": "/path/to/file"}`))

	assert.Nil(t, err)
	assert.IsType(t, crfppParser{}, parser)
	assert.True(t, parser.isConfigured())
}
