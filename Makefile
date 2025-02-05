build:
	mkdir -p bin
	go build -buildvcs=false -o bin/main ./src
