all: build

build:
	go build -o wikicmd
	mkdir -p bin
	mv ./wikicmd ./bin/.

docs_build:
	bash scripts/docs_build.sh
