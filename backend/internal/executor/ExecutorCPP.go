package executor

import "io/ioutil"

type ExecutorCPP struct{}

func NewExecutorCPP() *ExecutorCPP {
	return &ExecutorCPP{}
}

func (e *ExecutorCPP) ExecuteFromSource(source string) (output string, err error) {

	err = ioutil.WriteFile("source.cpp", []byte(source), 0644)
	if err != nil {
		return "", err
	}

	return "", nil
}
