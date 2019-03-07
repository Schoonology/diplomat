.PHONY := format watch clean e2e test generate
PKGS=$(shell go list -deps | grep http-assertion-tool)

main: main.go */*.go
	go build -o main

format:
	go fmt $(PKGS)

bin/mockery:
	GOBIN=`pwd`/bin go get github.com/vektra/mockery/.../

bin/templify:
	GOBIN=`pwd`/bin go get github.com/wlbr/templify

generate: bin/mockery bin/templify
	bin/mockery -all
	go generate ./...

watch:
	rg --files | entr -rc sh -c "make format && make main && make test && make e2e"

clean:
	rm -f *.go.* */*.go.*
	rm -f main

e2e:
	./main fixtures/test1.txt http://httpbin.org
	./main fixtures/test2.txt http://httpbin.org
	./main fixtures/test3.txt http://httpbin.org
	./main fixtures/test-post.txt http://httpbin.org	
	./main fixtures/test.markdown http://httpbin.org

test:
	go test $(PKGS)
