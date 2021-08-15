FROM golang:1.16 as build

RUN apt-get update 
RUN apt-get install -y build-essential

WORKDIR /build
COPY go.mod go.sum Makefile ./
RUN go mod download

COPY ./src ./src
COPY ./pkg ./pkg

RUN make build

WORKDIR /dist
RUN cp /build/bin/nvidia-node-tagger .

# ================================================
FROM nvcr.io/nvidia/cuda:11.0-base-ubi8
LABEL name="NVIDIA GPU Node Tagger"
LABEL description="Discovery GPU features and add labels to nodes in the cluster"
COPY ./LICENSE /license/LICENSE

COPY --from=build /dist/nvidia-node-tagger .

ENTRYPOINT ["/nvidia-node-tagger"]
