package crfpp

import (
	"os/exec"
	"syscall"
	"bytes"
	"fmt"
)

// ITagger is the interface that wraps crfpp functionality.
type ITagger interface {
	Test (modelFilePath string, dataFilePath string) (string, error)
	Learn (templateFilePath string, trainingFilePath string, modeFilePath string) (string, error)
}

// Tagger is the struct that wraps crfpp functionality.
type Tagger struct {

}

var execCommand = exec.Command

// Test runs crf_test on the supplied data
func (c Tagger) Test(modelFilePath string, dataFilePath string) (string, error) {
	cmd := execCommand("crf_test", "-v", "1", "-m", modelFilePath, dataFilePath)
	crfppOutput, err := execute(cmd)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	return crfppOutput, nil
}

// Learn runs crf_learn on the supplied data
func (c Tagger) Learn(templateFilePath string, trainingFilePath string, modeFilePath string) (string, error) {
	cmd := execCommand("crf_learn", templateFilePath, trainingFilePath)
	crfppOutput, err := execute(cmd)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	return crfppOutput, nil
}

// execute executes the given command and returns the output from Stdout.
func execute(cmd *exec.Cmd) (string, error) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	cmd.Stderr = cmdOutput

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s, %s", err.Error(), cmdOutput.Bytes())
	}

	return string(cmdOutput.Bytes()), nil
}