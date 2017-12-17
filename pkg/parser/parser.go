// Package parser provides a way to take a string set of ingredients and turn them into
// an array of Ingredient.
package parser

import (
	"github.com/chvck/ingredients-parser/pkg/ingredient"
	"fmt"
	"encoding/json"
)

// Parser is the interface that wraps the basic Parse method.
type Parser interface {
	isConfigured() bool
	Parse(ingredients string) ([]ingredient.Ingredient, error)
}

type parserConfig struct {
	ParserType string `json:"parsertype"`
}

func NewParser(data []byte) (Parser, error) {
	var cf parserConfig
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, err
	}
	parser, err := stringToStruct(cf.ParserType, data)
	if err != nil {
		return nil, err
	}

	return parser, nil
}

func stringToStruct(name string, data []byte) (Parser, error) {
	switch name {
	case "crfppParser":
		parser := crfppParser{}
		if err := parser.setConfig(data); err != nil {
			return nil, err
		}
		return parser, nil
	default:
		return nil, fmt.Errorf("%s is not a known struct name", name)
	}
}