// Package ingredient contains the base Ingredient type.
package ingredient

// Ingredient is the representation of an individual ingredient + info as a part of a recipe.
// e.g. 100 grams tomatoes, diced.
type Ingredient struct {
	unit string
	name string
	quantity string
	notes string
}

// SetUnit sets the unit aspect of the ingredient, e.g. gram.
func (i *Ingredient) SetUnit(unit string) *Ingredient {
	i.unit = unit
	return i
}

// SetName sets the name of the ingredient, e.g. tomato.
func (i *Ingredient) SetName(name string) *Ingredient {
	i.name = name
	return i
}

// SetQuantity sets the quantity of the ingredient, e.g. 100.
func (i *Ingredient) SetQuantity(quantity string) *Ingredient {
	i.quantity = quantity
	return i
}

// SetNotes sets any notes for the ingredient, e.g. diced.
func (i *Ingredient) SetNotes(notes string) *Ingredient {
	i.notes = notes
	return i
}

// Unit gets the unit aspect of the ingredient, e.g. gram.
func (i *Ingredient) Unit() string {
	return i.unit
}

// Name gets the name of the ingredient, e.g. tomato.
func (i *Ingredient) Name() string {
	return i.name
}

// Quantity gets the quantity of the ingredient, e.g. 100.
func (i *Ingredient) Quantity() string {
	return i.quantity
}

// Notes gets any notes for the ingredient, e.g. diced.
func (i *Ingredient) Notes() string {
	return i.notes
}
