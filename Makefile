test: clean
	go test -coverprofile cover.out

coverage: test
	go tool cover -html cover.out -o coverage.html
	open coverage.html

clean:
	go clean
	rm -f -- coverage.html
	rm -f -- cover.out