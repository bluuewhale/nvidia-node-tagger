build:
	goos=linux goarch=amd64 cgo_enabled=0 go build -installsuffix cgo -o  ./bin/nvidia-node-tagger ./src/main.go
