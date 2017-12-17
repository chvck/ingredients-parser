package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"os"
	"fmt"
	"github.com/chvck/ingredients-parser/pkg/ingredient"
)

func fakeExecCommand(funct string) func(command string, args...string) *exec.Cmd {
	return func(command string, args...string) *exec.Cmd {
		cs := []string{"-test.run="+funct, "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
}

func TestCrfppParser_Parse(t *testing.T) {
	execCommand = fakeExecCommand("TestHelperProcessSucceed")
	defer func(){ execCommand = exec.Command }()
	modelPath := "/path/to/model"
	config := &config{}
	config.ModelFilePath = modelPath

	parser := &CrfppParser{}
	parser.config = *config

	ingredients, err := parser.Parse("1 1/2 teaspoon fresh thyme leaves, finely chopped\n" +
		"2 tablespoons sherry vinegar\n2 tablespoons extra-virgin olive oil")

	assert.Nil(t, err)
	assert.Equal(t, 3, len(ingredients))
	assertIngredient(t, ingredients[0], []string{"thyme", "leaves"}, "teaspoon", "1 1/2",
	[]string{"fresh", ",", "finely", "chopped"})
	assertIngredient(t, ingredients[1], []string{"sherry", "vinegar"}, "tablespoons", "2",
		[]string{""})
	assertIngredient(t, ingredients[2], []string{"olive", "oil"}, "tablespoons", "2",
		[]string{"extra-virgin"})
}

func TestCrfppParser_ParseError(t *testing.T) {
	execCommand = fakeExecCommand("TestHelperProcessFail")
	defer func(){ execCommand = exec.Command }()
	modelPath := "/path/to/model"
	config := &config{}
	config.ModelFilePath = modelPath

	parser := &CrfppParser{}
	parser.config = *config

	ingredients, err := parser.Parse("1 1/2 teaspoon fresh thyme leaves, finely chopped")

	assert.Error(t, err)
	assert.Nil(t, ingredients)
}

// toCrfppFormat is really difficult to test as a part of testing the Parser function
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


func TestHelperProcessSucceed(t *testing.T){
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprint(os.Stdout,`# 0.148792
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
`)
	os.Exit(0)
}


func TestHelperProcessFail(t *testing.T){
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprint(os.Stdout, "feature_index.cpp(193) [mmap_.open(model_filename)] mmap.h(153) " +
		"[(fd = ::open(filename, flag | O_BINARY)) >= 0] open failed: file")
	os.Exit(-1)
}

func assertIngredient(t *testing.T, ingredient ingredient.Ingredient, n []string, u string, q string, no []string) {
	assert.Equal(t, n, ingredient.Name())
	assert.Equal(t, u, ingredient.Unit())
	assert.Equal(t, q, ingredient.Quantity())
	assert.Equal(t, no, ingredient.Notes())
}
