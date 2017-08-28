SHELL = bash
GOTOOLS = \
	golang.org/x/tools/cmd/cover \
	github.com/mattn/goveralls \
	github.com/hashicorp/consul \
	github.com/denisenkom/go-mssqldb \
	github.com/axw/gocov/gocov \
	github.com/stretchr/testify/assert \
	gopkg.in/matm/v1/gocov-html \
	github.com/urfave/cli

GOTAGS ?= consul-mirror
GOFILES ?= $(shell go list ./... | grep -v /vendor/)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Get the git commit
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(shell git describe --tags --always)
GIT_IMPORT=github.com/hashicorp/consul/version
GOLDFLAGS=-X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE)

export GOLDFLAGS

# all builds binaries for all targets
all: test bin

bin: tools
	go build .

# debug creates a gdb_sandbox to debug the code
debug:
	mkdir -p pkg/$(GOOS)_$(GOARCH)/ bin/
	go install -ldflags '$(GOLDFLAGS)' -tags '$(GOTAGS)'
	cp $(GOPATH)/bin/consul bin/
	cp $(GOPATH)/bin/consul pkg/$(GOOS)_$(GOARCH)

cov:
	gocov test $(GOFILES) | gocov-html > /tmp/coverage.html
	open /tmp/coverage.html

test: tools vet
	go test -tags '$(GOTAGS)' -i ./...
	go test $(GOTEST_FLAGS) -tags '$(GOTAGS)' -timeout 7m -v ./... 2>&1 >test$(GOTEST_FLAGS).log ; echo $$? > exit-code
	@echo "Exit code: `cat exit-code`" >> test$(GOTEST_FLAGS).log
	@echo "----"
	@grep -A5 'DATA RACE' test.log || true
	@grep -A10 'panic: test timed out' test.log || true
	@grep '^PASS' test.log | uniq || true
	@grep -A1 -- '--- FAIL:' test.log || true
	@grep '^FAIL' test.log || true
	@test "$$TRAVIS" == "true" && cat test.log || true
	@exit $$(cat exit-code)

cover:
	go test $(GOFILES) --cover

format:
	@echo "--> Running go fmt"
	@go fmt $(GOFILES)

vet:
	@echo "--> Running go vet"
	@go vet $(GOFILES); if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

tools:
	go get -u -v $(GOTOOLS)

.PHONY: all ci bin dev dist cov test cover format vet ui static-assets tools
