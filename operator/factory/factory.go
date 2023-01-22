package factory

import "github.com/masahiro331/funnel/pod"

type Factory interface {
	Create(number int) ([]Pod, error)
}

type Pod interface {
	Name() string
	Ready() (bool, error)
	Delete() error
	Target() string

	Exec(taskId, name string, args []string) error
	Statuses() ([]pod.Status, error)
	Pull(taskId string) ([]byte, error)
}
