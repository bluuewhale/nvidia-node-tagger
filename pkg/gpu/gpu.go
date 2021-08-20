package gpu

import (
	"errors"
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

const (
	MiB = 1048576
)

type GpuDeviceList struct {
	Devices   map[string]*GpuDevice `json:"devices"`
	MemorySum Memory                `json:"sum.memory"`
}

func NewGpuDeviceList() (*GpuDeviceList, error) {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return nil, errors.New(nvml.ErrorString(ret))
	}
	defer nvml.Shutdown()

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return nil, errors.New(nvml.ErrorString(ret))
	}

	devices := []*nvml.Device{}

	for i := 0; i < count; i++ {
		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			return nil, errors.New(nvml.ErrorString(ret))
		}

		devices = append(devices, &device)
	}

	return newGpuDeviceList(devices)
}

func newGpuDeviceList(devices []*nvml.Device) (*GpuDeviceList, error) {
	deviceList := map[string]*GpuDevice{}
	for i, device := range devices {
		gpuDevice, err := newGpuDevice(device)
		if err != nil {
			return nil, err
		}

		deviceList[fmt.Sprint(i)] = gpuDevice
	}

	gpuDeviceList := &GpuDeviceList{
		Devices:   deviceList,
		MemorySum: Memory{},
	}

	for _, gpuDevice := range deviceList {
		gpuDeviceList.MemorySum.Total += gpuDevice.Memory.Total
		gpuDeviceList.MemorySum.Used += gpuDevice.Memory.Used
		gpuDeviceList.MemorySum.Free += gpuDevice.Memory.Free
	}

	return gpuDeviceList, nil
}

type GpuDevice struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Memory Memory `json:"memory"`
}

func newGpuDevice(device *nvml.Device) (*GpuDevice, error) {
	uuid, ret := device.GetUUID()
	if ret != nvml.SUCCESS {
		return nil, errors.New(nvml.ErrorString(ret))
	}

	name, ret := device.GetName()
	if ret != nvml.SUCCESS {
		return nil, errors.New(nvml.ErrorString(ret))
	}

	memory, ret := device.GetMemoryInfo()
	if ret != nvml.SUCCESS {
		return nil, errors.New(nvml.ErrorString(ret))
	}

	gpuDevice := &GpuDevice{
		UUID: uuid,
		Name: name,
		Memory: Memory{
			Total: memory.Total / MiB,
			Used:  memory.Used / MiB,
			Free:  memory.Free / MiB,
		},
	}
	return gpuDevice, nil
}

type Memory struct {
	Total uint64 `json:"total"` // MiB
	Used  uint64 `json:"used"`  // MiB
	Free  uint64 `json:"free"`  // MiB
}
