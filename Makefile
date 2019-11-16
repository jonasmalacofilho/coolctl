all: build

build:
	@go build

install:
	@go install

commit: lint vet format test

lint:
	@golint -set_exit_status ./...

vet:
	@go vet ./...

format:
	@goimports -d .

test:
	@go test ./... -v -coverprofile .coverage.txt
	@go tool cover -func .coverage.txt

coverage: test
	@go tool cover -html=.coverage.txt

dep:
	@go get -v -d ./...

cyclo:
	@gocyclo -over 15 .