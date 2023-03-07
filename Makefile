# TODO: Build using https://goreleaser.com/
build:
	CGO_ENABLED=0 go build -v

build-linux:
	CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -v

# TODO: Lint using golangci-lint
