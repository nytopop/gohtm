default: build

all: build test run
	echo All Success!

test: 
	go test -v ./...

build:
	go build ./...
	go build

run: build
	./gohtm

cpu: build run
	go tool pprof -web cpuprofile

clean:
	rm cpuprofile gohtm
