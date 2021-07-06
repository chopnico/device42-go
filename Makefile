build:
	go mod vendor
	go fmt ./...
	go build -o build/device42/device42 cmd/device42/device42.go
clean:
	rm -rf build
	go clean
