.PHONY : all build test clean

all : clean test build

install:
	cp bin/index_server /usr/local/bin

build:
	mkdir bin
	go build -o bin/index_server *.go

clean:
	rm -rf bin/

test:
	go test -v -bench . -cover
