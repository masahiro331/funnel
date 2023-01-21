package factory

import (
	"github.com/masahiro331/funnel/pod"
	"strconv"
)

var (
	_ Factory = &MockEC2Factory{}
	_ Pod     = &MockEC2{}
)

type MockEC2Factory struct {
}

func (m MockEC2Factory) Create(number int) ([]Pod, error) {
	var pods []Pod
	for i := 0; i < number; i++ {
		pods = append(pods, MockEC2{
			name: strconv.Itoa(i),
		})
	}
	return pods, nil
}

type MockEC2 struct {
	name string
}

func (m MockEC2) Target() string {
	//TODO implement me
	panic("implement me")
}

func (m MockEC2) Exec(taskId, name string, args []string) error {
	//TODO implement me
	panic("implement me")
}

func (m MockEC2) Statuses() ([]pod.Status, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockEC2) Pull(taskId string) (pod.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockEC2) Name() string {
	return m.name
}

func (m MockEC2) Ready() (bool, error) {
	return true, nil
}

func (m MockEC2) Delete() error {
	return nil
}
