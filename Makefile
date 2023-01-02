OS=linux
ARCH=amd64

GOENV=env GOOS=$(OS) GOARCH=$(ARCH)

TEST_FLAGS =

GOCMD=go
GOBUILD=$(GOENV) $(GOCMD) build
GOTEST=$(GOCMD) test -coverprofile=coverage-report.out $(TEST_FLAGS)
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOSEC=gosec -fmt=json -out=gl-sast-report.json
GOCOVERTOOLS=$(GOCMD) tool cover -html=coverage-report.out -o coverage-report.html
GOCLEAN=$(GOCMD) clean
GOLIST=$(GOCMD) list

PROJECT_PATH=github.com/sungup/t-fio

BUILD_LDFLAG="-extldflags=-static"

SOURCES=$(shell find . -name "*.go")

# build:

format:
	@$(GOFMT) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)

test: format
	@$(GOVET) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)
	@$(GOTEST) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)
	@$(GOCOVERTOOLS)
	@$(GOSEC) ./...

clean:
	@$(GOCLEAN)
	@$(GOCLEAN) -testcache
	# TODO add remove output files
	# @rm -f $(TARGETS)