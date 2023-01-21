package pod

import "bytes"

type Result struct {
	Status  Status
	Content *bytes.Buffer
}

const (
	Running = "Running"
	Done    = "Done"
)

type Status struct {
	TaskId string
	State  string
}

type PullRequest struct {
	TaskId string
}

type ExecRequest struct {
	TaskId string
	Name   string
	Args   []string
}
