# nvidia-node-tagger
nvidia node feature discovery service for kubernetes gpu clutser. Strongly inspired by [NVIDIA GPU feature discovery](https://github.com/NVIDIA/gpu-feature-discovery)

&nbsp;
# Major Features
## Update Annotations
Add GPU informations to node's annotations

```
apiVersion: v1
kind: Node
metadata:
  annotations:
    nvidia-node-tagger/devices.0.memory.free: "22698"
    nvidia-node-tagger/devices.0.memory.total: "22698"
    nvidia-node-tagger/devices.0.memory.used: "0"
    nvidia-node-tagger/devices.0.name: Quadro RTX 6000
    nvidia-node-tagger/devices.0.uuid: GPU-0d773341-6454-5a4b-e7c6-bf4e14a309e1
    nvidia-node-tagger/devices.1.memory.free: "22698"
    nvidia-node-tagger/devices.1.memory.total: "22698"
    nvidia-node-tagger/devices.1.memory.used: "0"
    nvidia-node-tagger/devices.1.name: Quadro RTX 6000
    nvidia-node-tagger/devices.1.uuid: GPU-c3165ee7-f7f1-b585-81e0-ec20ba44bea3
    nvidia-node-tagger/devices.2.memory.free: "22698"
    nvidia-node-tagger/devices.2.memory.total: "22698"
    nvidia-node-tagger/devices.2.memory.used: "0"
    nvidia-node-tagger/devices.2.name: Quadro RTX 6000
    nvidia-node-tagger/devices.2.uuid: GPU-195e3c16-ad55-89f9-acf5-e9b3b3b511cf
    nvidia-node-tagger/sum.memory.free: "68094"
    nvidia-node-tagger/sum.memory.total: "68094"
    nvidia-node-tagger/sum.memory.used: "0"
    ...
```
&nbsp;

## Update Capacity
Add GPU memory capacity to node's capatcity 

```
apiVersion: v1
kind: Node
status:
  allocatable:
    nvidia-node-tagger/vram: "68094"
    nvidia-node-tagger/vram.devices.0: "22698"
    nvidia-node-tagger/vram.devices.1: "22698"
    nvidia-node-tagger/vram.devices.2: "22698"
    ...
  capacity:
    nvidia-node-tagger/vram: "68094"
    nvidia-node-tagger/vram.devices.0: "22698"
    nvidia-node-tagger/vram.devices.1: "22698"
    nvidia-node-tagger/vram.devices.2: "22698"
    ...
```
&nbsp;

This enables pods to be scheduled based vram capacity of nodes

```
apiVersion: v1
kind: Pod
metadata:
    name: gpu-schedule-demo
spec:
spec:
  containers:
  - name: gpu-schedule-demo
    image: busybox:latest
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
        nvidia-node-tagger/vram: "10000" # MiB
      limits:
        memory: "128Mi"
        cpu: "500m"
        nvidia-node-tagger/vram: "10000" # MiB
...

```