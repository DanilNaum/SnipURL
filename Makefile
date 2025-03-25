mockgen:
	go install github.com/matryer/moq@v0.4.0 && \
	go generate ./...

coverage:
	go install github.com/matryer/moq@v0.4.0 && \
	go generate ./... && \
	go test -v -p 1 -cover -coverprofile="coverage.out" ./... && \
	go tool cover -html="./coverage.out" -o "coverage.html"
