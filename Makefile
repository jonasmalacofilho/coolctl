BINARY = coolctl
GOARCH = amd64

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = -ldflags "-X main.Version=${BRANCH}:${COMMIT}"

GOCMD = GO111MODULE=on go

all: build

.PHONY: build
build:
	${GOCMD} build ${LDFLAGS} -o ./bin/${BINARY} ./main.go

.PHONY: linux
linux:
	GOOS=linux GOARCH=${GOARCH} ${GOCMD} build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} .

.PHONY: macos
macos:
	GOOS=darwin GOARCH=${GOARCH} ${GOCMD} build ${LDFLAGS} -o ${BINARY}-macos-${GOARCH} .

.PHONY: windows
windows:
	GOOS=windows GOARCH=${GOARCH} ${GOCMD} build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe .

cross: linux macos windows

install:
	${GOCMD} install

commit: lint vet format test

lint:
	golint -set_exit_status ./...

vet:
	${GOCMD} vet ./...

format:
	goimports -d .

test:
	${GOCMD} test ./... -v -coverprofile .coverage.txt
	${GOCMD} tool cover -func .coverage.txt

coverage: test
	${GOCMD} tool cover -html=.coverage.txt

dep:
	${GOCMD} get -v -d ./...

cyclo:
	@gocyclo -over 15 .

tidy:
	${GOCMD} mod tidy