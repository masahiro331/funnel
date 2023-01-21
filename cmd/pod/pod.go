package pod

import (
	"fmt"
	"github.com/masahiro331/funnel/pod"
	"net/http"
)

func Run() {
	http.HandleFunc("/exec", pod.Exec)
	http.HandleFunc("/pull", pod.Pull)
	http.HandleFunc("/statuses", pod.Statuses)
	http.HandleFunc("/healthz", pod.Healthz)
	fmt.Println("listen...", "0.0.0.0:6332")
	http.ListenAndServe("0.0.0.0:6332", nil)
}
