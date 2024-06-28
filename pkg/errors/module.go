package errors

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wangweihong/gotoolbox/pkg/netutil"
)

var currentModule ModuleGetter

type ModuleGetter interface {
	PID() int
	Host() string
	Name() string
	String() string
}

func NewModuleGetter(module, host string, pid int) ModuleGetter {
	return simpleModule{
		name: module,
		host: host,
		pid:  pid,
	}
}

type simpleModule struct {
	name string
	host string
	pid  int
}

func (s simpleModule) PID() int {
	return s.pid
}

func (s simpleModule) Host() string {
	return s.host
}

func (s simpleModule) Name() string {
	return s.name
}

func (s simpleModule) String() string {
	return fmt.Sprintf("host:%s,pid:%d,module:%s", s.host, s.pid, s.name)
}

func GetModuleInfo() ServiceInfo {
	return ServiceInfo{
		Host: currentModule.Host(),
		Pid:  currentModule.PID(),
		Name: currentModule.Name(),
	}
}

func UpdateModuleInfo(getter ModuleGetter) {
	currentModule = getter
}

func ModuleString() string {
	return currentModule.String()
}

//nolint:gochecknoinits
func init() {
	moduleIP := "127.0.0.1"

	ipList, _ := netutil.GetIPAddrs(false)
	for _, ip := range ipList {
		if ip != "127.0.0.1" {
			moduleIP = ip
			break
		}
	}

	currentModule = &simpleModule{
		name: filepath.Base(os.Args[0]),
		host: moduleIP,
		pid:  os.Getpid(),
	}
}
