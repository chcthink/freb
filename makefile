git :
	git add . && git cz && git push origin

BINARY=freb
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -trimpath -ldflags="-s -w" -o $(BINARY)-$(GOOS)-$(GOARCH)

build-macos:
	$(MAKE) build GOOS=darwin GOARCH=arm64
build-linux:
	$(MAKE) build GOOS=linux GOARCH=amd64
build-windows:
	$(MAKE) build GOOS=windows GOARCH=amd64

build-all:
	$(MAKE) build-macos
	$(MAKE) build-linux
	$(MAKE) build-windows

clean:
	rm -f $(BINARY)-*