.PHONY: all clean clean-coverage install install-deps install-tools test test-verbose test-with-coverage

export ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
export PKG := github.com/maksimov/liconv
export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

$(eval $(ARGS):;@:) # turn arguments into do-nothing targets
export ARGS

ifdef ARGS
	PKG_TEST := $(ARGS)
else
	PKG_TEST := $(PKG)/...
endif

all: install-tools install-deps install test

clean:
	go clean -i $(PKG)/...
	go clean -i -race $(PKG)/...
clean-coverage:
	find $(ROOT_DIR) | grep .coverprofile | xargs rm
install:
	go install -v $(PKG)/...
install-deps:
	go get -t -v $(PKG)/...
	go build -v $(PKG)/...
install-tools:
	# Install code coverage tools
	go get -u -v github.com/onsi/ginkgo/ginkgo/...
test:
	go test -race $(PKG_TEST)
test-verbose:
	go test -race -v $(PKG_TEST)
test-with-coverage:
	ginkgo -r -cover -race