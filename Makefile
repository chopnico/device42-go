build:
	go mod download
	go build -o build/device42/device42 cmd/device42/main.go
clean:
	rm -rf build
	go clean
