.PHONY: run build clean

run:
	go run cmd/unisat-monitor/main.go

build:
	@echo "Building the binary executable..."
	go build -o bin/unisat-monitor cmd/unisat-monitor/main.go

clean:
	@echo "Cleaning up bin directory..."
	rm -rf bin/
