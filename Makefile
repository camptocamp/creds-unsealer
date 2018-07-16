DEPS = $(wildcard */*.go)
VERSION = $(shell git describe --always --dirty)

SHARED_LIBS = "libxml-2.0 augeas"
CGO_CFLAGS = $(shell pkg-config --cflags $(SHARED_LIBS))
CGO_LDFLAGS = $(shell pkg-config --static --libs $(SHARED_LIBS))

all: test creds-unsealer creds-unsealer.1

creds-unsealer: creds-unsealer.go $(DEPS)
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	CGO_ENABLED=1 GOOS=linux \
	  go build -a \
		  -ldflags='-linkmode external -extldflags -static -X main.version=$(VERSION)' \
	    -installsuffix cgo -o $@ $<
	strip $@

creds-unsealer.1: creds-unsealer
	./creds-unsealer -m > $@

lint:
	@ go get -v github.com/golang/lint/golint
	@for file in $$(git ls-files '*.go' | grep -v '_workspace/'); do \
		export output="$$(golint $${file} | grep -v 'type name will be used as docker.DockerInfo')"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

vet: creds-unsealer.go
	go vet $<

imports: creds-unsealer.go
	dep ensure
	goimports -d $<

test: imports lint vet
	go test -v ./...

coverage:
	rm -rf *.out
	go test -coverprofile=coverage.out
	for i in config engines handler metrics providers util volume orchestrators; do \
	 	go test -coverprofile=$$i.coverage.out github.com/camptocamp/creds-unsealer/$$i; \
		tail -n +2 $$i.coverage.out >> coverage.out; \
		done

clean:
	rm -f creds-unsealer creds-unsealer.1

.PHONY: all imports lint vet test coverage clean
