package observe

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

var (
	ErrInvalidIPStr = errors.New("invalid ip address string")
	ErrInvalidIP    = errors.New("invalid ip address")
)

func AnonymizeIP(remoteAddr string) (string, error) {
	ipStr, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ipStr = remoteAddr
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", ErrInvalidIPStr
	}

	var mask net.IPMask

	if ip.To4() != nil {
		mask = net.CIDRMask(20, 32)
	} else if ip.To16() != nil {
		mask = net.CIDRMask(32, 128)
	} else {
		return "", ErrInvalidIP
	}

	return ip.Mask(mask).String(), nil
}

func GetAnonymizedIP(req *http.Request) (string, error) {
	ip := req.RemoteAddr
	fwd := req.Header.Get("X-Forwarded-For")

	if fwd != "" {
		ips := strings.Split(fwd, ", ")
		ip = ips[0]
	}

	return AnonymizeIP(ip)
}
