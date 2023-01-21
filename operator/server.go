package operator

import (
	"encoding/json"
	"github.com/masahiro331/funnel/operator/factory"
	"github.com/masahiro331/funnel/operator/operation"
	"github.com/masahiro331/funnel/operator/operation/netscan"
	"net/http"

	"github.com/masahiro331/funnel/utils"
)

const (
	NmapOperation = "Nmap Operation"
)

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
