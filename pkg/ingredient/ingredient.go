package ingredient

type Ingredient struct {
	unit string
	name string
	quantity string
	notes string
}

func (i *Ingredient) SetUnit(unit string) *Ingredient {
	i.unit = unit
	return i
}

func (i *Ingredient) SetName(name string) *Ingredient {
	i.name = name
	return i
}

func (i *Ingredient) SetQuantity(quantity string) *Ingredient {
	i.quantity = quantity
	return i
}

func (i *Ingredient) SetNotes(notes string) *Ingredient {
	i.notes = notes
	return i
}

func (i *Ingredient) Unit() string {
	return i.unit
}

func (i *Ingredient) Name() string {
	return i.name
}

func (i *Ingredient) Quantity() string {
	return i.quantity
}

func (i *Ingredient) Notes() string {
	return i.notes
}
