package cgroups

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/containerd/cgroups/v3/cgroup2"
)

const (
	CPU_PERIOD = uint64(100000)
)

type cgroup struct {
	name      string
	resources *cgroup2.Resources
}

func (cg *cgroup) slice() string {
	return cg.name + ".slice"
}

func (cg *cgroup) path(file string) string {
	if file == "" {
		return filepath.Join("/sys/fs/cgroup", cg.slice())
	}
	return filepath.Join("/sys/fs/cgroup", cg.slice(), file)
}

func (cg *cgroup) root() string {
	return cg.path("")
}

func (cg *cgroup) exists() bool {
	_, err := os.Stat(cg.root())
	return err == nil
}

func (cg *cgroup) load() (*cgroup2.Manager, error) {
	return cgroup2.LoadSystemd("/", cg.slice())
}

func (cg *cgroup) create() error {
	if cg.exists() {
		return nil // cgroup already exists, ignoring
	}

	// dummy PID of -1 is used for creating a "general slice" to be used as a parent cgroup.
	// see https://github.com/containerd/cgroups/blob/1df78138f1e1e6ee593db155c6b369466f577651/v2/manager.go#L732-L735
	_, err := cgroup2.NewSystemd("/", cg.slice(), -1, cg.resources)
	return err
}

func (cg *cgroup) destroy() error {
	if !cg.exists() {
		return nil // it's already gone, let's ignore
	}

	m, err := cg.load()
	if err != nil {
		return nil // probably already unloaded
	}
	_ = m.DeleteSystemd()

	return nil
}

func (cg *cgroup) setMemoryMax(limit int64) {
	cg.resources.Memory.Max = &limit
}

func (cg *cgroup) setMemoryHigh(limit int64) {
	cg.resources.Memory.High = &limit
}

func (cg *cgroup) setMemorySwap(limit int64) {
	cg.resources.Memory.Swap = &limit
}

func (cg *cgroup) setPIDsMax(limit int64) {
	cg.resources.Pids.Max = limit
}

func (cg *cgroup) setCPUMax(quota int64, period uint64) {
	cg.resources.CPU.Max = cgroup2.NewCPUMax(&quota, &period)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func CreateCGroup(appName, name string, cpuQuota, memHardMax, memSoftMax, memSwap int64, pidsMax int64) (path string, err error) {
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	name = re.ReplaceAllString(name, "")
	cg := &cgroup{
		name: appName + name,
		resources: &cgroup2.Resources{
			CPU:    &cgroup2.CPU{},
			Memory: &cgroup2.Memory{},
			Pids:   &cgroup2.Pids{},
		},
	}

	if pidsMax < 16 {
		pidsMax = 16 // we need at least this many to enter the container
	}

	cg.setCPUMax(cpuQuota, CPU_PERIOD)
	cg.setPIDsMax(pidsMax)
	cg.setMemoryMax(memHardMax)
	cg.setMemoryHigh(memSoftMax)
	cg.setMemorySwap(memSwap)

	if err := cg.create(); err != nil {
		return "", err
	}

	return cg.root(), nil
}

func DestroyCGroup(appName, name string) error {
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	name = re.ReplaceAllString(name, "")
	cg := &cgroup{
		name: appName + name,
	}
	return cg.destroy()
}
