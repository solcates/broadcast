package broadcast

import (
	"net"
	"strings"
)

var netInterface = net.InterfaceAddrs

func init() {
	netInterface = net.InterfaceAddrs
}

func getLocalIPs() (ips map[string]struct{}) {
	ips = make(map[string]struct{})
	ips["localhost"] = struct{}{}
	ips["127.0.0.1"] = struct{}{}
	addrs, err := netInterface()
	if err != nil {
		return
	}
	for _, address := range addrs {
		ips[strings.Split(address.String(), "/")[0]] = struct{}{}
	}
	return
}
