package operator

import (
	"fmt"
	"github.com/masahiro331/funnel/operator"
	"net/http"
)

func Run() {
	http.HandleFunc("/netscan", operator.NetworkScan)
	http.HandleFunc("/operations", operator.Operations)
	http.HandleFunc("/result", operator.Result)

	fmt.Println("listen...", "0.0.0.0:6331")
	http.ListenAndServe("0.0.0.0:6331", nil)
}
