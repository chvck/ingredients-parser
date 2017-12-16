package parser

import "github.com/chvck/ingredients-parser/pkg/ingredient"

type Parser interface {
	SetConfig(config []byte) error
	ParseIngredients(ingredients string) ([]ingredient.Ingredient, error)
}