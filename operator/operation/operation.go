package operation

import (
	"github.com/masahiro331/funnel/operator/factory"
	"log"
	"time"
)

var (
	Operations []Operation
)

type State string

const (
	Initialize State = "Initialize"
	Running    State = "Running"
	Finish     State = "Finish"
)

func init() {
	go func() {
		for {
			for i, o := range Operations {
				switch o.State() {
				case Initialize:
					continue
				case Running:
					err := o.Action()
					if err != nil {
						log.Printf("error: %s", err.Error())
					}

				case Finish:
					for _, pod := range o.Pods() {
						err := pod.Delete()
						if err != nil {
							log.Printf("error: %s", err.Error())
						}
					}
					Operations = append(Operations[:i], Operations[i+1:]...)
				}
			}
			time.Sleep(time.Second)
		}
	}()
}

type Operation interface {
	Name() string
	Init(f factory.Factory) error
	State() State
	Pods() []factory.Pod
	Action() error
}

func Register(o Operation) {
	Operations = append(Operations, o)
}

type Task struct {
	Id      string
	Command string
	Args    []string
}
