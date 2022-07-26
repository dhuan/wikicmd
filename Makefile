all: build

build:
	go build -o wikicmd
	mkdir -p bin
	mv ./wikicmd ./bin/.

docs_build:
	bash scripts/docs_build.sh

test_unit:
	go test -v ./tests/unit/*_test.go

test_e2e:
	go test -v ./tests/e2e/*_test.go

prepare_e2e:
	bash scripts/setup_mock.sh
	bash scripts/setup_fakevim.sh

test: test_unit test_e2e
