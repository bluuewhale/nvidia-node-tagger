VERSION=0.0.1

build:
	goos=linux goarch=amd64 cgo_enabled=0 go build -installsuffix cgo -o  ./bin/nvidia-node-tagger ./src/main.go

package:
	docker build -t koko8624/nvidia-node-tagger:$(VERSION) .

push:
	docker push koko8624/nvidia-node-tagger:$(VERSION)
