package operator

type NmapScanRequest struct {
	Parallel   int
	CIDRBlocks []string
}

type ResultRequest struct {
	OperationId string
}

type OperationResponse struct {
	OperationId string
	State       string
}
