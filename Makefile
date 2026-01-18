build: build-windows build-linux

build-windows:
	GOOS=windows GOARCH=amd64 go build -o build/chroma.exe
build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/chroma

.PHONY: build build-windows build-linux
