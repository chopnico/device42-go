build:
	go mod vendor
	go fmt ./...
	GOOS=linux GOARCH=amd64 go build -o build/device42/device42_linux_amd64 cmd/device42/device42.go
	GOOS=windows GOARCH=amd64 go build -o build/device42/device42_windows_amd64.exe cmd/device42/device42.go
clean:
	rm -rf build
	go clean
