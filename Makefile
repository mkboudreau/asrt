all: deps build

build:
	go fmt ./...
	go build
test:
	go test ./...
test-report: report-dir
	go test ./... -v | go-junit-report > reports/tests.xml
deps:
	go get -d -t
imports:
	go get -v $(go list -f '{{range .Imports}}{{ . }} {{end}}')
tools:
	go get github.com/jstemmer/go-junit-report
	go get github.com/golang/lint/golint
	go get github.com/ryancox/gobench2plot
clean:
	rm -rf reports/
bench:
	go test ./... -bench=.
bench-report: report-dir
	go test ./... -bench=. -benchmem | gobench2plot > reports/benchmarks.xml
check: report-dir
	golint ./... > reports/lint.txt
	go vet ./...  2> reports/vet.txt
report-dir:
	mkdir -p reports
ci: deps tools test-report cover bench-report check

