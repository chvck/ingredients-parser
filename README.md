Ingredients Parser
==================

Ingredients Parser is a golang library designed to parse human readable lists of ingredients. The idea behind
this is to normalize the quantities/names of ingredients so that lists on ingredients can easily be accumulated
across recipes.

Parsers
-------
There is presently only one parser type; crfpp. The parser requires that crfpp be installed and that `crf_test`be
accessible from command line. 

To create a parser use `parser.NewParser` passing in a json byte array containing the path to the crfpp model file
as `modelfilepath` and the type of parser, in this case `crfppParser` as `parsertype`.
e.g. `parser.NewParser([]byte({"parsertype": "crfppParser", "modelfilepath": "/path/to/model"}))`.

The config passed to the crfpp parser also accepts `unit`, `quantity` and `name` fields. The values of these correspond
to the features in the crfpp output. If these aren't present in the config then they will default to the following values:

`unit`: `unit`

`quantity`: `qty`

`name`: `name`

Any feature with that name outside of these three will stored in the Ingredient as a `note`.

Functionality to assist in creating a model file coming soon.

Credits
-------
[crfpp](https://github.com/taku910/crfpp)

Implementation of crfpp parser heavily based off of [NYTimes ingredient-phrase-tagger](https://github.com/NYTimes/ingredient-phrase-tagger/)