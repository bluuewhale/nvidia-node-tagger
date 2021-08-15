package gpu

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"

	"github.com/mitchellh/mapstructure"
)

var (
	KEYS = []string{
		"Index", "Name", "DriverVersion", "MemoryTotal", "MemoryUsed", "MemoryFree", "Temperature",
	}
)

func NewGpuStatList() (GpuStatList, error) {

	out, err := exec.Command(
		"nvidia-smi",
		"--query-gpu=index,name,driver_version,memory.total,memory.used,memory.free,temperature.gpu",
		"--format=csv,noheader,nounits",
	).Output()

	if err != nil {
		return GpuStatList{}, err
	}

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()

	if err != nil {
		return GpuStatList{}, err
	}

	fmt.Printf("%v\n", rows)

	inner := []GpuStat{}
	for _, row := range rows {
		gpuStat := newGpuStat(row)
		inner = append(inner, gpuStat)
	}

	return newGpuStatList(inner), nil
}

type GpuStatList struct {
	inner          []GpuStat
	MemoryTotalSum int `mapstructure:",omitempty"` // MiB
	MemoryUsedSum  int `mapstructure:",omitempty"` // MiB
	MemoryFreeSum  int `mapstructure:",omitempty"` // MiB
}

func newGpuStatList(inner []GpuStat) GpuStatList {
	gpuStatList := GpuStatList{}
	gpuStatList.inner = inner

	for _, gpuStat := range inner {
		gpuStatList.MemoryTotalSum += gpuStat.MemoryTotal
		gpuStatList.MemoryUsedSum += gpuStat.MemoryUsed
		gpuStatList.MemoryFreeSum += gpuStat.MemoryFree
	}

	return gpuStatList
}

func (g *GpuStatList) Iterate() []GpuStat {
	return g.inner
}

type GpuStat struct {
	Index         int     `mapstructure:",omitempty"`
	Name          string  `mapstructure:",omitempty"`
	DriverVersion string  `mapstructure:",omitempty"`
	MemoryTotal   int     `mapstructure:",omitempty"` // MiB
	MemoryUsed    int     `mapstructure:",omitempty"` // MiB
	MemoryFree    int     `mapstructure:",omitempty"` // MiB
	Temperature   float32 `mapstructure:",omitempty"`
}

func newGpuStat(data []string) GpuStat {
	gpuStat := GpuStat{}
	_map := map[string]string{}

	for i, v := range data {
		_map[KEYS[i]] = v
	}

	mapstructure.WeakDecode(_map, &gpuStat)
	return gpuStat
}
