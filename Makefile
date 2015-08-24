NAME=hostctl
ARCH=$(shell uname -m)
VERSION=0.1.0dev

.PHONY: build release docs

build:
	glu build darwin,linux

deps:
	go get github.com/gliderlabs/glu
	go get -d .

release:
	glu release v$(VERSION)

docs:
	boot2docker ssh "sync; sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'" || true
	docker run --rm -it -p 8000:8000 -v $(PWD):/work gliderlabs/pagebuilder mkdocs serve
