default: build

all: build test run
	echo Successful build!

test: 
	go test -v ./...

build:
	go build ./...
