package executor

type Executor interface {
	ExecuteFromSource(source string) (output string, err error)
}
