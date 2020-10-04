fmt:
	go fmt ./...

test: fmt
	go test ./... -cover -coverprofile=cover.out

cov: test
	go tool cover -html=cover.out

web:
	cd api/cmd && go run main.go
