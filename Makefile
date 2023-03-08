# TODO: Build using https://goreleaser.com/
build:
	CGO_ENABLED=0 go build -v

build-linux:
	CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -v

release-linux: build-linux
	tar cvzf redis-ha-check-linux-amd64.tar.gz redis-ha-check

# TODO: Lint using golangci-lint
