package parser

import (
	"encoding/json"
	"github.com/chvck/ingredients-parser/pkg/ingredient"
	"fmt"
	"strconv"
	"regexp"
	"strings"
	"io/ioutil"
	"os"
	"errors"
	"github.com/kljensen/snowball"
	"github.com/chvck/ingredients-parser/internal/crfpp"
)

// crfppParser parses a set of ingredients by using crfpp (https://github.com/taku910/crfpp).
// This parser implementation requires crf++ to be installed and is expensive
// due to calling out to crf_test which writes a file which must then be read and interpreted.
type crfppParser struct {
	config config
	tagger crfpp.ITagger
}

// config is a set of configuration values for use during parsing
type config struct {
	ModelFilePath string `json:"modelfilepath"`
	Unit          string `json:"unit"`
	Name          string `json:"name"`
	Quantity      string `json:"quantity"`
}

// setConfig sets the parser up ready for use. crfppParser expects the config to contain
// the path to the model file.
func (p *crfppParser) setConfig(data []byte) error {
	var cf config
	if err := json.Unmarshal(data, &cf); err != nil {
		return err
	}

	p.config = cf
	p.tagger = crfpp.Tagger{}

	return nil
}

// isConfigured returns whether not the parse has been setup
func (p crfppParser) isConfigured() bool {
	return p.config.ModelFilePath != ""
}

// Parse accepts a human readable list of ingredients and returns a slice of Ingredient.
// Example Input: 1 cup sugar\n500 grams flour
func (p crfppParser) Parse(ingredientsStr string) ([]ingredient.Ingredient, error) {
	if !p.isConfigured() {
		return nil, errors.New("parser has not been configured")
	}

	// we have to convert the ingredients string to a format the crfpp recognises
	formatted := toCrfppFormat(ingredientsStr)

	filename := "tmp"
	if err := ioutil.WriteFile(filename, []byte(formatted), 0666); err != nil {
		return nil, err
	}
	defer os.Remove(filename)

	crfppOutput, err := p.tagger.Test(p.config.ModelFilePath, filename)
	if err != nil {
		return nil, err
	}

	return p.createIngredientsFromCrfpp(crfppOutput), nil
}

/*
createIngredientsFromCrfpp converts a crfpp output into a slice of Ingredient.
Example Input:
# 0.950492
1	I1	L20	NoCAP	NoPAREN	B-QTY/0.979304
cup	I2	L20	NoCAP	NoPAREN	B-UNIT/0.978106
sugar	I3	L20	NoCAP	NoPAREN	B-NAME/0.984194

 */
func (p crfppParser) createIngredientsFromCrfpp(crfppOutput string) []ingredient.Ingredient {
	var ing ingredient.Ingredient
	var ingredients []ingredient.Ingredient
	newIng := false
	unit := p.unit()
	quantity := p.quantity()
	name := p.name()
	re := regexp.MustCompile(`^[BI]\-`)
	for _, line := range strings.Split(crfppOutput, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}

		if line == "" && newIng {
			ingredients = append(ingredients, ing)
			ing = ingredient.Ingredient{}
			newIng = false
		} else {
			newIng = true
			columns := strings.Split(strings.Trim(line, " "), "\t")
			token := strings.Trim(columns[0], " ")
			token = unclumpFractions(token)

			split := strings.Split(columns[len(columns)-1], "/")
			tag := re.ReplaceAllString(split[0], "")
			tag = strings.ToLower(tag)
			//confidence := split[1]

			if tag == unit {
				singled, err := singularize(token)
				if err != nil {
					singled = token
				}
				ing.SetUnit(singled)
			} else if tag == quantity {
				ing.SetQuantity(token)
			} else if tag == name {
				ing.AddName(token)
			} else {
				ing.AddNote(token)
			}
		}
	}

	return ingredients
}

// unit returns the unit value in config, or "unit" if nothing set
func (p crfppParser) unit() string {
	unit := p.config.Unit
	if unit == "" {
		unit = "unit"
	}

	return unit
}

// name returns the name value in config, or "unit" if nothing set
func (p crfppParser) name() string {
	name := p.config.Name
	if name == "" {
		name = "name"
	}

	return name
}

// quantity returns the quantity value in config, or "unit" if nothing set
func (p crfppParser) quantity() string {
	quantity := p.config.Quantity
	if quantity == "" {
		quantity = "qty"
	}

	return quantity
}

/*
toCrfppFormat converts a human readable list of ingredients to an input that crfpp expects.
Example Input: 1 cup sugar\n500 grams flour
Example output:
# 0.950492
1	I1	L20	NoCAP	NoPAREN	B-QTY
cup	I2	L20	NoCAP	NoPAREN	B-UNIT
sugar	I3	L20	NoCAP	NoPAREN	B-NAME

 */
func toCrfppFormat(ingredients string) string {
	split := strings.Split(ingredients, "\n")
	re := regexp.MustCompile("<[^<]+?>")
	var parsed []string
	for _, line := range split {
		cleaned := re.ReplaceAllString(line, "")
		tokens := tokenise(cleaned)

		for i, token := range tokens {
			features := getFeatures(token, i+1, tokens)
			combined := strings.Join(append([]string{token}, features...), "\t")
			parsed = append(parsed, combined)
		}
		parsed = append(parsed, "")
	}

	return strings.Join(parsed, "\n")
}

// getFeatures returns a list of features for a given token.
func getFeatures(token string, index int, tokens []string) []string {
	length := len(tokens)
	caps := "No"
	if startsWithCapital(token) {
		caps = "Yes"
	}

	return []string{
		fmt.Sprintf("I%s", strconv.Itoa(index)),
		fmt.Sprintf("L%s", bucketLength(length)),
		caps + "CAP",
		"NoPAREN",
	}
}

// startsWithCapital returns true if a given token starts with a capital letter.
func startsWithCapital(s string) bool {
	re := regexp.MustCompile(`^[A-Z]`)
	return re.MatchString(s)
}

// bucketLength buckets the length of the ingredient into 6 buckets.
func bucketLength(length int) string {
	buckets := [5]int{4, 8, 12, 16, 20}
	bucketed := "X"

	for _, bucket := range buckets {
		if length < bucket {
			bucketed = strconv.Itoa(bucket)
		}
	}

	return bucketed
}

// Tokenise on parenthesis, punctuation, spaces and slashes (that aren't part of a fraction).
func tokenise(s string) []string {
	tokenised := clumpFractions(s)
	tokenised = hideFractions(tokenised)

	tokenised = strings.Replace(tokenised, "/", " ", -1)
	//// Make fractions look like fractions again
	tokenised = unhideFractions(tokenised)

	re := regexp.MustCompile(`[(\d+!)?\d\/\d]+|[A-Za-z*?!']+|[(),]`)
	split := re.FindAllString(tokenised, -1)
	return split
}

/*
clumpFractions replaces the whitespace between the integer and fractional part of a quantity
with a dollar sign, so it's interpreted as a single token. The rest of the string is left alone.
	clumpFractions("aaa 1 2/3 bbb")
	# => "aaa 1$2/3 bbb"
 */
func clumpFractions(s string) string {
	re := regexp.MustCompile(`(\d+)\s+(\d)\/(\d)`)
	return re.ReplaceAllString(s, "$1!$2/$3")
}

// unclumpFractions replaces $ in fractions with spaces. The reverse of clumpFractions.
func unclumpFractions(s string) string {
	re := regexp.MustCompile(`(\d+)!(\d)\/(\d)`)
	return re.ReplaceAllString(s, "$1 $2/$3")
}

// hideFractions replaces the slash in fractions with a £.
func hideFractions(s string) string {
	re := regexp.MustCompile(`(\d+!)?(\d)\/(\d)`)
	return re.ReplaceAllString(s, "$1$2£$3")

}

// unhideFractions replaces the £ in fractions with a /.
func unhideFractions(s string) string {
	re := regexp.MustCompile(`(\d+)£(\d)`)
	return re.ReplaceAllString(s, "$1/$2")
}

// singularize normalizes plurals into singulars
func singularize(s string) (string, error) {
	return snowball.Stem(s, "english", true)
}
