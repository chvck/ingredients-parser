// Package parser provides a way to take a string set of ingredients and turn them into
// an array of Ingredient.
package parser

import "github.com/chvck/ingredients-parser/pkg/ingredient"

// Parser is the interface that wraps the basic Parse method.
type Parser interface {
	SetConfig(config []byte) error
	Parse(ingredients string) ([]ingredient.Ingredient, error)
}