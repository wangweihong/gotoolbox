package netutil

import (
	"net"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

func GetDefaultInterface() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			if strings.HasPrefix(iface.Name, "e") ||
				strings.HasPrefix(iface.Name, "br") {
				return &iface, nil
			}
		}
	}

	return nil, errors.New("no suitable network interface found")
}
