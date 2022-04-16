all: build

build:
	go build -o wikicmd
	mkdir -p bin
	mv ./wikicmd ./bin/.
