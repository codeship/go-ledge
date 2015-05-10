.PHONY: \
	all \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	generate \
	build \
	install \
	lint \
	vet \
	errcheck \
	pretest \
	test \
	cov \
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
	go generate ./...

build: deps generate
	go build ./...

install: deps generate
	go install ./...

lint: testdeps generate
	go get -v github.com/golang/lint/golint
	golint ./...

vet: testdeps generate
	go get -v golang.org/x/tools/cmd/vet
	go vet ./...

errcheck:
	go get -v github.com/kisielk/errcheck
	errcheck ./...

pretest: lint vet errcheck

test: testdeps generate pretest
	go test -test.v ./...

cov: testdeps generate
	go get -v github.com/axw/gocov/gocov
	go get golang.org/x/tools/cmd/cover
	gocov test | gocov report

clean:
	go clean -i ./...
