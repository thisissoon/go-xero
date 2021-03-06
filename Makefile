test:
	go test -coverprofile=coverage.out -cover -v

coverage: test
	go tool cover -html=coverage.out

godoc:
	godoc -http=:8000
