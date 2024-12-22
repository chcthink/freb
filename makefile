BINARY=freb
VERSION ?= $(shell git describe --tags --abbrev=0)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

git:
	git add . && git cz && git push origin

tag:
	$(MAKE) git && git tag $(VERSION) && git push origin $(VERSION)

tag-build:
	$(MAKE) tag && $(MAKE) build-all

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -trimpath -ldflags="-s -w -X '$(BINARY)/cmd.version=$(VERSION)'" -o $(BINARY)-$(GOOS)-$(GOARCH)

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