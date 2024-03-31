BINARY = myapp

all:
	make build

build:
	go build -o $(BINARY) main.go