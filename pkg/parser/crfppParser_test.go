package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/chvck/ingredients-parser/pkg/ingredient"
	"errors"
)

type FakeTaggerSuccess struct {
}

type FakeTaggerFail struct {
}

func (f FakeTaggerSuccess) Test(modelFilePath string, dataFilePath string) (string, error) {
	return `# 0.148792
1!1/2	I1	L20	NoCAP	NoPAREN	B-QTY/0.522172
teaspoon	I2	L20	NoCAP	NoPAREN	B-UNIT/0.978515
fresh	I3	L20	NoCAP	NoPAREN	B-COMMENT/0.738304
thyme	I4	L20	NoCAP	NoPAREN	B-NAME/0.949495
leaves	I5	L20	NoCAP	NoPAREN	I-NAME/0.945802
,	I6	L20	NoCAP	NoPAREN	B-COMMENT/0.389967
finely	I7	L20	NoCAP	NoPAREN	I-COMMENT/0.568856
chopped	I8	L20	NoCAP	NoPAREN	I-COMMENT/0.759612


# 0.336084
2	I1	L8	NoCAP	NoPAREN	B-QTY/0.984317
tablespoons	I2	L8	NoCAP	NoPAREN	B-UNIT/0.638708
sherry	I3	L8	NoCAP	NoPAREN	I-NAME/0.601504
vinegar	I4	L8	NoCAP	NoPAREN	I-NAME/0.939174

# 0.860110
2	I1	L8	NoCAP	NoPAREN	B-QTY/0.998181
tablespoons	I2	L8	NoCAP	NoPAREN	B-UNIT/0.998844
extra-virgin	I3	L8	NoCAP	NoPAREN	B-COMMENT/0.864063
olive	I4	L8	NoCAP	NoPAREN	B-NAME/0.884804
oil	I5	L8	NoCAP	NoPAREN	I-NAME/0.997967
`, nil
}

func (f FakeTaggerFail) Test(modelFilePath string, dataFilePath string) (string, error) {
	return "", errors.New("feature_index.cpp(193) [mmap_.open(model_filename)] mmap.h(153) " +
		"[(fd = ::open(filename, flag | O_BINARY)) >= 0] open failed: file")
}

func TestCrfppParser_Parse(t *testing.T) {
	modelPath := "/path/to/model"
	config := &config{}
	config.ModelFilePath = modelPath
	config.Unit = "unit"
	config.Name = "name"
	config.Quantity = "qty"

	p := &crfppParser{}
	p.config = *config
	p.tagger = FakeTaggerSuccess{}

	ingredients, err := p.Parse("1 1/2 teaspoon fresh thyme leaves, finely chopped\n" +
		"2 tablespoons sherry vinegar\n2 tablespoons extra-virgin olive oil")

	assert.Nil(t, err)
	assert.Equal(t, 3, len(ingredients))
	assertIngredient(t, ingredients[0], []string{"thyme", "leaves"}, "teaspoon", "1 1/2",
		[]string{"fresh", ",", "finely", "chopped"})
	assertIngredient(t, ingredients[1], []string{"sherry", "vinegar"}, "tablespoon", "2",
		[]string{""})
	assertIngredient(t, ingredients[2], []string{"olive", "oil"}, "tablespoon", "2",
		[]string{"extra-virgin"})
}

func TestCrfppParser_ParseError(t *testing.T) {
	modelPath := "/path/to/model"
	config := &config{}
	config.ModelFilePath = modelPath
	config.Unit = "unit"
	config.Name = "name"
	config.Quantity = "qty"

	p := &crfppParser{}
	p.config = *config
	p.tagger = FakeTaggerFail{}

	ingredients, err := p.Parse("1 1/2 teaspoon fresh thyme leaves, finely chopped")

	assert.Error(t, err)
	assert.Nil(t, ingredients)
}

// toCrfppFormat is really difficult to test as a part of testing the IParser function
func TestCrfppParser_toCrfppFormat(t *testing.T) {
	output := toCrfppFormat("1 1/2 cup sugar\n500 grams flour")
	expected := `1!1/2	I1	L20	NoCAP	NoPAREN
cup	I2	L20	NoCAP	NoPAREN
sugar	I3	L20	NoCAP	NoPAREN

500	I1	L20	NoCAP	NoPAREN
grams	I2	L20	NoCAP	NoPAREN
flour	I3	L20	NoCAP	NoPAREN
`

	assert.Equal(t, expected, output)
}

func TestCrfppParser_setConfig(t *testing.T) {
	p := &crfppParser{}
	err := p.setConfig([]byte(`{"modelfilepath": "/path/to/file", "unit": "unit",
		"quantity": "qty", "name": "name"}`))

	assert.Nil(t, err)
	assert.Equal(t, "/path/to/file", p.config.ModelFilePath)
}

func TestCrfppParser_isConfigured(t *testing.T) {
	p := &crfppParser{}
	modelPath := "/path/to/model"
	config := &config{}
	config.ModelFilePath = modelPath
	config.Unit = "unit"
	config.Name = "name"
	config.Quantity = "qty"

	p.config = *config

	assert.True(t, p.isConfigured())
}

func TestCrfppParser_isConfiguredFalse(t *testing.T) {
	p := &crfppParser{}
	assert.False(t, p.isConfigured())
}

func assertIngredient(t *testing.T, ingredient ingredient.Ingredient, n []string, u string, q string, no []string) {
	assert.Equal(t, n, ingredient.Name())
	assert.Equal(t, u, ingredient.Unit())
	assert.Equal(t, q, ingredient.Quantity())
	assert.Equal(t, no, ingredient.Notes())
}
