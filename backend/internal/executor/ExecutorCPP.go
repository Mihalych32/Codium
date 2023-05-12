package executor

import "fmt"

type ExecutorCPP struct{}

func NewExecutorCPP() *ExecutorCPP {
	return &ExecutorCPP{}
}

func (e *ExecutorCPP) ExecuteFromSource(source string) (output string, err error) {

	fmt.Println("Executing CPP source code:")
	fmt.Println(source)

	return "", nil
}
