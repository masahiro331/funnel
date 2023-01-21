package pod

import (
	"github.com/masahiro331/funnel/pod"
	"net/http"
)

func Run() {
	http.HandleFunc("/exec", pod.Exec)
	http.HandleFunc("/pull", pod.Pull)
	http.HandleFunc("/statuses", pod.Statuses)
	http.HandleFunc("/healthz", pod.Healthz)
	http.ListenAndServe(":6332", nil)
}
