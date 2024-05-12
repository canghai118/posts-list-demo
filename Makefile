build:
	go build

dep:
	go mod tidy
	go mod vendor

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build