########################################################################################

.PHONY = deps fmt test

########################################################################################

deps:
	go get -v pkg.re/essentialkaos/ek.v7

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

test:
	go get -v pkg.re/check.v1
	go test -covermode=count .
