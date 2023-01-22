package netscan

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/masahiro331/funnel/operator/factory"
	"github.com/masahiro331/funnel/operator/operation"
	p "github.com/masahiro331/funnel/pod"
)

const (
	parallel = 10
)

func NewNmapOperation(name string, parallels int, hosts []string) operation.Operation {
	var tasks []operation.Task
	for _, h := range hosts {
		tasks = append(tasks, operation.Task{
			Id:      strings.Join([]string{"task", "nmap", string(time.Now().UnixMicro())}, "-"),
			Command: "nmap",
			Args:    []string{"-sU", h}, // UDP scan
		})
		tasks = append(tasks, operation.Task{
			Id:      strings.Join([]string{"task", "nmap", string(time.Now().UnixMicro())}, "-"),
			Command: "nmap",
			Args:    []string{"-sS", h}, // TCP scan
		})
	}
	return &NmapOperation{
		Id:        strings.Join([]string{"operation", "nmap", string(time.Now().UnixMicro())}, "-"),
		name:      name,
		tasks:     tasks,
		Parallels: parallels,
	}
}

func (n *NmapOperation) Name() string {
	return n.name
}

func (n *NmapOperation) Init(f factory.Factory) (err error) {
	n.pods, err = f.Create(n.Parallels)
	if err != nil {
		return err
	}
	n.state = operation.Initialize
	go func() {
		var flag bool
		for {
			for _, pod := range n.pods {
				flag, err = pod.Ready()
				if err != nil {
					log.Printf("error: %s", err.Error())
				}

				if !flag {
					break
				}
			}
			if flag {
				n.state = operation.Running
				break
			}
		}
	}()
	return nil
}

type NmapOperation struct {
	Id        string
	name      string
	Parallels int
	state     operation.State
	pods      []factory.Pod
	tasks     []operation.Task
}

func (n *NmapOperation) Action() error {
	if len(n.tasks) == 0 {
		n.state = operation.Finish
	}
	for _, pod := range n.pods {
		statuses, err := pod.Statuses()
		if err != nil {
			return err
		}
		if len(statuses) < parallel {
			for i := 0; i < parallel-len(statuses); i++ {
				if len(n.tasks) == 0 {
					break
				}
				task := n.tasks[0]
				n.tasks = n.tasks[1:]
				pod.Exec(task.Id, strings.Join(task.Args, " "), task.Args)
			}
		}
		for _, status := range statuses {
			if status.State != p.Done {
				continue
			}
			result, err := pod.Pull(status.TaskId)
			if err != nil {
				return err
			}
			fmt.Println(result)
		}
	}

	return nil
}

func (n *NmapOperation) State() operation.State {
	return n.state
}

func (n *NmapOperation) Pods() []factory.Pod {
	return n.pods
}
