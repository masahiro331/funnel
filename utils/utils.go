package utils

import (
	"encoding/json"
	"net"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

func Error(err error, writer http.ResponseWriter, status int) {
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(&ErrorResponse{
		Error: err.Error(),
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
