package operator

type NmapScanRequest struct {
	Parallel   int
	CIDRBlocks []string
}
