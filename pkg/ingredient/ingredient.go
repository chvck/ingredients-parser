// Package ingredient contains the base Ingredient type.
package ingredient

import (
	"fmt"
	"strings"
)

// Ingredient is the representation of an individual ingredient + info as a part of a recipe.
// e.g. 100 grams tomatoes, diced.
type Ingredient struct {
	unit string
	names []string
	quantity string
	notes []string
}

// SetUnit sets the unit aspect of the ingredient, e.g. gram.
func (i *Ingredient) SetUnit(unit string) *Ingredient {
	i.unit = unit
	return i
}

// AddName add a name for the ingredient, e.g. tomato.
func (i *Ingredient) AddName(name string) *Ingredient {
	i.names = append(i.names, name)
	return i
}

// SetQuantity sets the quantity of the ingredient, e.g. 100.
func (i *Ingredient) SetQuantity(quantity string) *Ingredient {
	i.quantity = quantity
	return i
}

// AddNotes add a note for the ingredient, e.g. diced.
func (i *Ingredient) AddNote(note string) *Ingredient {
	i.notes = append(i.notes, note)
	return i
}

// Unit gets the unit aspect of the ingredient, e.g. gram.
func (i Ingredient) Unit() string {
	return i.unit
}

// Name gets the name of the ingredient, e.g. tomato.
func (i Ingredient) Name() []string {
	return i.names
}

// Quantity gets the quantity of the ingredient, e.g. 100.
func (i Ingredient) Quantity() string {
	return i.quantity
}

// Notes gets any notes for the ingredient, e.g. diced.
func (i Ingredient) Notes() []string {
	return i.notes
}

func (i Ingredient) String() string {
	return fmt.Sprintf("%s %s %s, %s",i.quantity, i.unit, strings.Join(i.names, " "),
		strings.Join(i.notes, " "))
}