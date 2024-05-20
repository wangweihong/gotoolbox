package errors

import (
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/wangweihong/gotoolbox/sets"
)

var CurrentModule ModuleGetter

type ModuleGetter interface {
	PID() int
	IP() string
	Name() string
	String() string
}

type simpleModule struct {
	name string
	ip   string
	pid  int
}

func (s simpleModule) PID() int {
	return s.pid
}

func (s simpleModule) IP() string {
	return s.ip
}

func (s simpleModule) Name() string {
	return s.name
}

func (s simpleModule) String() string {
	return fmt.Sprintf("host:%s,pid:%d,module:%s", s.IP(), s.PID(), s.Name())
}

func Caller() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	funcName := "???"
	if fp := runtime.FuncForPC(pc); fp != nil {
		funcName = fp.Name()
	}

	dir, filename := path.Split(file)
	// show package name for error stack
	if dir != "" {
		parent := filepath.Base(dir)
		filename = filepath.Join(parent, filename)
	}

	fileList := strings.Split(funcName, ".")
	funcName = fileList[len(fileList)-1]
	format := "file:" + filename + ",func:" + funcName + ",line:" + strconv.FormatInt(int64(line), 10)
	return format
}

func UpdateModuleInfo(getter ModuleGetter) {
	currentModule = getter
}

func getIPAddrs(wantIpv6 bool) ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, iface := range ifaces {
		if sets.NewString(iface.Name).HasAnyPrefix("e", "br") {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok {
					if p4 := ipnet.IP.To4(); len(p4) == net.IPv4len {
						ips = append(ips, ipnet.IP.String())
					} else if len(ipnet.IP) == net.IPv6len && wantIpv6 {
						ips = append(ips, ipnet.IP.String())
					}
				}
			}
		}
	}

	return ips, nil
}

// nolint:gochecknoinits
func init() {
	moduleIP := "127.0.0.1"
	ipList, _ := getIPAddrs(false)
	for _, ip := range ipList {
		if ip != "127.0.0.1" {
			moduleIP = ip
			break
		}
	}

	currentModule = &simpleModule{
		name: filepath.Base(os.Args[0]),
		ip:   moduleIP,
		pid:  os.Getpid(),
	}
}
