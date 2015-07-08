.PHONY: \
	all \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	build \
	install \
	lint \
	vet \
	errcheck \
	pretest \
	test \
	cov \
	checkproto \
	proto \
	clean

all: test

deps:
	go get -d -v ./...

updatedeps:
	go get -d -v -u -f ./...

testdeps:
	go get -d -v -t ./...

updatetestdeps:
	go get -d -v -t -u -f ./...

generate:
	go ./...

build: deps
	go build ./...

install: deps
	go install ./...

lint: testdeps
	go get -v github.com/golang/lint/golint
	for file in $$(git ls-files '*.go' | grep -v '\.pb\.go'); do \
		golint $$file; \
	done

vet: testdeps
	go get -v golang.org/x/tools/cmd/vet
	go vet ./...

errcheck: testdeps
	go get -v github.com/kisielk/errcheck
	errcheck ./...

pretest: lint vet errcheck

test: testdeps pretest
	go test -test.v ./...

cov: testdeps
	go get -v github.com/axw/gocov/gocov
	go get golang.org/x/tools/cmd/cover
	gocov test | gocov report

checkproto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ if [ "$$(protoc --version)" != "libprotoc 3.0.0" ]; then \
	  echo "error: proto 3 must be installed" >&2; \
		exit 1; \
	fi

proto: checkproto
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I . --go_out=. ledge.proto

clean:
	go clean -i ./...
