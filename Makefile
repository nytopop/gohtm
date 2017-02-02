default: build

test: 
	go test -v ./...

build: test
	go build ./...
	go build

run: build
	./gohtm

cpu: build run
	go tool pprof -web cpuprofile

clean:
	rm cpuprofile gohtm
