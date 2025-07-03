mockgen:
	go install github.com/matryer/moq@v0.4.0 && \
	go generate ./...

coverage:
	go install github.com/matryer/moq@v0.4.0 && \
	go generate ./... && \
	go test -v -p 1 -cover -coverprofile="coverage.out" ./... && \
	go tool cover -html="./coverage.out" -o "coverage.html"

pprof:
	curl -sK -v http://localhost:8080/debug/pprof/heap > heap.out
	go tool pprof -http=":9090" -seconds=30 heap.out 

bench_coverage:
	go test -benchmem  -cover -coverprofile="bench_coverage.out" ./... && \
	go tool cover -html="./bench_coverage.out" -o "bench_coverage.html"

multichecker_test:
	go build -o multichecker  ./cmd/multichecker/main.go
	./multichecker ./cmd/shortener/... 
	@echo "Проверка ./cmd/shortener/... завершена успешно" 
	./multichecker ./internal/... 
	@echo "Проверка ./internal/...  завершена успешно" 
	./multichecker ./pkg/... 
	@echo "Проверка ./pkg/...  завершена успешно" 