package internal

type ExecuteResult struct {
	StatusCode int
	Body       map[string]any
}

type Executor interface {
	Execute(call Call) (*ExecuteResult, error)
}
