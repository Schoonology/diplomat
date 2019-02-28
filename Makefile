.PHONY := format watch clean e2e test generate
PKGS=$(shell go list -deps | grep http-assertion-tool)

main: main.go */*.go
	go build -o main

format:
	go fmt $(PKGS)

bin/mockery:
	GOBIN=`pwd`/bin go get github.com/vektra/mockery/.../

generate: bin/mockery
	rm -rf mocks
	bin/mockery -all

watch:
	rg --files | entr -rc sh -c "make format && make main && make test && make e2e"

clean:
	rm -f *.go.* */*.go.*
	rm -f main

e2e:
	./main test.txt http://httpbin.org

test:
	go test $(PKGS)
