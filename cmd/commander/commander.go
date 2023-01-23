package commander

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/masahiro331/funnel/operator"
	"github.com/masahiro331/funnel/operator/report"
	"log"
	"net/http"
)

func Run(host string, operationId string) {
	monitor(host, operationId)
}

func monitor(host string, operationId string) {

	for {
		r := operator.ResultRequest{
			OperationId: operationId,
		}
		b, err := json.Marshal(&r)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post(fmt.Sprintf("%s/result", host), "", bytes.NewReader(b))
		if err != nil {
			log.Fatal(err)
		}

		var reports []report.Report
		if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
			log.Fatal(err)
		}
		for _, repo := range reports {
			fmt.Println(string(repo.Content))
		}
		resp.Body.Close()
	}

}
