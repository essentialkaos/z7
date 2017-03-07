########################################################################################

.PHONY = fmt deps deps-test test

########################################################################################

deps:
	go get -v pkg.re/essentialkaos/ek.v7

deps-test:
	go get -v pkg.re/check.v1
	go get -v pkg.re/essentialkaos/ek.v7

test:
	go test -covermode=count .

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

########################################################################################

