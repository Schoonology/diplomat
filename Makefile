.PHONY := format watch clean e2e test

main: main.go
	go build -o main

format:
	go fmt main.go

watch:
	rg --files | entr -rc sh -c "make format && make main && make test && make e2e"

clean:
	rm -f *.go.* */*.go.*
	rm -f main

e2e:
	./main test.txt httpbin.org:80

test:
	go test
