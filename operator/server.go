package operator

import (
	"encoding/json"
	"github.com/masahiro331/funnel/operator/factory"
	"github.com/masahiro331/funnel/operator/operation"
	"github.com/masahiro331/funnel/operator/operation/netscan"
	"github.com/masahiro331/funnel/operator/report"
	"net/http"

	"github.com/masahiro331/funnel/utils"
)

const (
	NmapOperation = "Nmap Operation"
)

func Operations(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var r []OperationResponse
	for _, o := range operation.Operations {
		r = append(r, OperationResponse{
			OperationId: o.Name(),
			State:       string(o.State()),
		})
	}
	b, err := json.Marshal(&r)
	if err != nil {
		utils.Error(err, writer, http.StatusInternalServerError)
	}
	writer.Write(b)
}

func Result(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	r := ResultRequest{}
	if err := json.NewDecoder(request.Body).Decode(&r); err != nil {
		utils.Error(err, writer, http.StatusBadRequest)
		return
	}
	reports := report.MemStore.Pop(r.OperationId)
	b, err := json.Marshal(&reports)
	if err != nil {
		utils.Error(err, writer, http.StatusInternalServerError)
	}
	writer.Write(b)
}

func NetworkScan(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	n := NmapScanRequest{}
	if err := json.NewDecoder(request.Body).Decode(&n); err != nil {
		utils.Error(err, writer, http.StatusBadRequest)
		return
	}

	var targets []string
	for _, cdir := range n.CIDRBlocks {
		hosts, err := utils.Hosts(cdir)
		if err != nil {
			utils.Error(err, writer, http.StatusBadRequest)
		}

		targets = append(targets, hosts...)
	}

	nmapOperation := netscan.NewNmapOperation(NmapOperation, n.Parallel, targets)
	err := nmapOperation.Init(&factory.EC2Factory{})
	if err != nil {
		utils.Error(err, writer, http.StatusInternalServerError)
	}

	operation.Register(nmapOperation)
}
