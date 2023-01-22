package operation

import (
	"github.com/masahiro331/funnel/operator/factory"
	"github.com/masahiro331/funnel/operator/report"
	"log"
)

var (
	Operations []Operation
)

type State string

const (
	Initialize State = "Initialize"
	Running    State = "Running"
	Finish     State = "Finish"
	Done       State = "Done"
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
					err := o.Done()
					if err != nil {
						log.Printf("error: %s", err.Error())
					}
				case Done:
					n := report.MemStore.Len(o.Name())
					if n == 0 {
						Operations = append(Operations[:i], Operations[i+1:]...)
					}
				}
			}
		}
	}()
}

type Operation interface {
	Name() string
	Init(f factory.Factory) error
	State() State
	Pods() []factory.Pod
	Action() error
	Done() error
}

func Register(o Operation) {
	Operations = append(Operations, o)
}

type Task struct {
	Id      string
	Command string
	Args    []string
}
