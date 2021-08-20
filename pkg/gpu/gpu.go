package gpu

import (
	"errors"
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

const (
	MiB = 1048576
)

type GpuDeviceInfos struct {
	Devices   map[string]*GpuDeviceInfo `json:"devices"`
	MemorySum Memory                    `json:"memory"`
}

func NewGpuDeviceInfos() (*GpuDeviceInfos, error) {
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

	deviceInfos := map[string]*GpuDeviceInfo{}
	for i, device := range devices {
		gpuDevice, err := newGpuDeviceInfo(device)
		if err != nil {
			return nil, err
		}

		deviceInfos[fmt.Sprint(i)] = gpuDevice
	}

	gpuDeviceInfos := &GpuDeviceInfos{
		Devices:   deviceInfos,
		MemorySum: Memory{},
	}

	for _, gpuDevice := range deviceInfos {
		gpuDeviceInfos.MemorySum.Total += gpuDevice.Memory.Total
		gpuDeviceInfos.MemorySum.Used += gpuDevice.Memory.Used
		gpuDeviceInfos.MemorySum.Free += gpuDevice.Memory.Free
	}

	return gpuDeviceInfos, nil
}

type GpuDeviceInfo struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Memory Memory `json:"memory"`
}

func newGpuDeviceInfo(device *nvml.Device) (*GpuDeviceInfo, error) {
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

	gpuDevice := &GpuDeviceInfo{
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

// ==============================================================

type GpuVramCapacity struct {
	Devices map[string]uint64 `json:"vram.devices"`
	Sum     uint64            `json:"vram"`
}

func NewGpuVramCapacity() (*GpuVramCapacity, error) {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return nil, errors.New(nvml.ErrorString(ret))
	}
	defer nvml.Shutdown()

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return nil, errors.New(nvml.ErrorString(ret))
	}

	gpuVramCapacity := &GpuVramCapacity{
		Devices: make(map[string]uint64),
		Sum:     0,
	}

	for i := 0; i < count; i++ {
		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			return nil, errors.New(nvml.ErrorString(ret))
		}

		memory, ret := device.GetMemoryInfo()
		if ret != nvml.SUCCESS {
			return nil, errors.New(nvml.ErrorString(ret))
		}

		gpuVramCapacity.Devices[fmt.Sprint(i)] = memory.Total / MiB
		gpuVramCapacity.Sum += memory.Total / MiB
	}

	return gpuVramCapacity, nil
}
