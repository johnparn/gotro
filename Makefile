build:
	mkdir -p build
	rm -rf build/*
	go build -ldflags="-s -w" -o build/demo main.go
	# CGO_ENABLED=0 CC=gcc GOOS=linux GOARCH=amd64 go build -tags static -ldflags "-s -w" -o build/demo main.go

run:
	./build/demo

install:
	sudo apt install libsdl2-mixer-dev libsdl2-gfx-dev

.PHONY: all test clean build install
