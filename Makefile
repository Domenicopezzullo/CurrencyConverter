.PHONY: build build-linux build-windows build-macos clean

build: build-linux build-windows build-macos
	@echo "All builds complete"

build-linux:
	@mkdir -p build/linux
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/convert-currency main.go
	@echo "Linux build complete"

build-windows:
	@mkdir -p build/windows
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/windows/convert-currency.exe main.go
	@echo "Windows build complete"

build-macos:
	@mkdir -p build/macos
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/macos/convert-currency main.go
	@echo "macOS build complete"

clean: 
	@rm -rf build/*