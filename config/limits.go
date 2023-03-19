package config

import (
	"fmt"
	"strings"

	"github.com/toxyl/devbox/cgroups"
	"github.com/toxyl/devbox/host"
	"github.com/toxyl/glog"
)

type LimitsMemory struct {
	Hard float64 `mapstructure:"hard"`
	Soft float64 `mapstructure:"soft"`
	Swap float64 `mapstructure:"swap"`
}

type Limits struct {
	CPU  float64      `mapstructure:"cpu"`
	Mem  LimitsMemory `mapstructure:"mem"`
	PIDs int64        `mapstructure:"pids"`
}

func (dlc *Limits) GetCPUQuota() int64 {
	if dlc == nil {
		panic("there are no CPU limits defined")
	}
	return int64(float64(cgroups.CPU_PERIOD) * dlc.CPU)
}

func (dlc *Limits) GetMemoryLimits() (memTotal, memHardMax, memSoftMax, memSwap int64) {
	if dlc == nil {
		panic("there are no memory limits defined")
	}
	mem := float64(host.TotalMemory())
	mhx := mem * dlc.Mem.Hard
	msm := mhx * dlc.Mem.Soft
	ms := mhx * dlc.Mem.Swap
	return int64(mem), int64(mhx), int64(msm), int64(ms)
}

func (dlc *Limits) String() string {
	cpuQuota := dlc.GetCPUQuota()
	memTotal, memHardMax, memSoftMax, memSwap := dlc.GetMemoryLimits()
	return fmt.Sprintf(
		"Limits: %s CPU (quota: %s), %s RAM (%s soft / %s hard / %s total, %s swap), %s PIDs",
		strings.ReplaceAll(glog.Auto(dlc.CPU), "%", "%%"),
		glog.Auto(cpuQuota),
		strings.ReplaceAll(glog.Auto(dlc.Mem.Hard), "%", "%%"),
		glog.HumanReadableBytesIEC(memSoftMax),
		glog.HumanReadableBytesIEC(memHardMax),
		glog.HumanReadableBytesIEC(memTotal),
		glog.HumanReadableBytesIEC(memSwap),
		glog.Int(dlc.PIDs),
	)

}
