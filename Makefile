NAME=hostctl
VERSION=0.1.0

build:
	go build -ldflags "-X main.Version $(VERSION)"
