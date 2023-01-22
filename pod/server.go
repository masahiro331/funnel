package pod

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/masahiro331/funnel/utils"
	"io"
	"net/http"
	"os/exec"
	"sync"
)

var (
	m sync.Map
)

func Exec(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	e := ExecRequest{}
	if err := json.NewDecoder(request.Body).Decode(&e); err != nil {
		utils.Error(fmt.Errorf("decode exec request error: %v", err), writer, http.StatusBadRequest)
		return
	}

	if e.TaskId == "" {
		utils.Error(errors.New("TaskId is required"), writer, http.StatusBadRequest)
		return
	}

	_, ok := m.Load(e.TaskId)
	if ok {
		utils.Error(errors.New("TaskId is exist"), writer, http.StatusBadRequest)
		return
	}

	go func(taskId string) {
		m.Store(taskId, Result{
			Status: Status{State: Running, TaskId: taskId},
		})
		buf := bytes.NewBuffer(nil)
		cmd := exec.Command(e.Name, e.Args...)
		cmd.Stdout = buf
		cmd.Stderr = buf
		cmd.Run()
		m.Store(taskId, Result{
			Status:  Status{State: Done, TaskId: taskId},
			Content: buf,
		})
	}(e.TaskId)

}

func Statuses(writer http.ResponseWriter, _ *http.Request) {
	var statuses []Status
	m.Range(func(k, v any) bool {
		result := v.(Result)
		statuses = append(statuses, result.Status)
		return true
	})
	if err := json.NewEncoder(writer).Encode(&statuses); err != nil {
		utils.Error(fmt.Errorf("decode statuses request error: %v", err), writer, http.StatusInternalServerError)
		return
	}
}

func Pull(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	e := PullRequest{}
	if err := json.NewDecoder(request.Body).Decode(&e); err != nil {
		utils.Error(fmt.Errorf("decode pull request error: %v", err), writer, http.StatusBadRequest)
		return
	}
	v, exist := m.LoadAndDelete(e.TaskId)
	if !exist {
		utils.Error(errors.New("does not taskId"), writer, http.StatusNotFound)
		return
	}
	result := v.(Result)
	if result.Status.State == Done {
		io.Copy(writer, result.Content)
	} else {
		// TODO
	}
}

func Healthz(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, "OK")
}
