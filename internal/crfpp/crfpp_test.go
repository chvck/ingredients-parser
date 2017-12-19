package crfpp

import (
	"testing"
	"os/exec"
	"os"
	"github.com/stretchr/testify/assert"
	"fmt"
)

var successOutput = `# 0.148792
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
`

var failOutput = "feature_index.cpp(193) [mmap_.open(model_filename)] mmap.h(153) " +
	"[(fd = ::open(filename, flag | O_BINARY)) >= 0] open failed: file"

func fakeExecCommand(funct string) func(command string, args...string) *exec.Cmd {
	return func(command string, args...string) *exec.Cmd {
		cs := []string{"-test.run="+funct, "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
}

func TestTagger_TestSuccess(t *testing.T) {
	execCommand = fakeExecCommand("TestHelperProcessSucceed")
	defer func(){ execCommand = exec.Command }()

	tagger := Tagger{}
	output, err := tagger.Test("/path/to/model", "/path/to/data")

	assert.Nil(t, err)
	assert.Equal(t, successOutput, output)
}

func TestTagger_TestFail(t *testing.T) {
	execCommand = fakeExecCommand("TestHelperProcessFail")
	defer func(){ execCommand = exec.Command }()

	tagger := Tagger{}
	output, err := tagger.Test("/path/to/model", "/path/to/data")

	assert.Equal(t, "", output)
	assert.Error(t, err)
}

func TestHelperProcessSucceed(t *testing.T){
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprint(os.Stdout, successOutput)
	os.Exit(0)
}


func TestHelperProcessFail(t *testing.T){
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprint(os.Stdout, failOutput)
	os.Exit(-1)
}
