package executor

import (
	"io/ioutil"
	"os/exec"
)

type ExecutorCPP struct{}

func NewExecutorCPP() *ExecutorCPP {
	return &ExecutorCPP{}
}

func (e *ExecutorCPP) ExecuteFromSource(source string) (output string, err error) {

	err = ioutil.WriteFile("source.cpp", []byte(source), 0644)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("gcc", "-lstdc++", "source.cpp")
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	cmd = exec.Command("chmod", "+x", "a.out")
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("./a.out").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
