# Version
VERSION=0.3.0

# Architecture
ARCH=amd64

# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
LDFLAGS="-X main.VERSION=$(VERSION)"
BINARY_NAME=sws

all: build

build-all: build build-linux build-darwin build-windows

build: 
	$(GOBUILD) -o $(BINARY_NAME) -ldflags $(LDFLAGS) -v

clean: 
	@rm -rf ./sws-binaries/linux_$(ARCH)/sws
	@rm -rf ./sws-binaries/darwin_$(ARCH)/sws
	@rm -rf ./sws-binaries/windows_$(ARCH)/sws.exe
	$(GOCLEAN)
      
# Cross compilation
build-linux:
	@mkdir -p ./sws-binaries/linux_$(ARCH)
	GOOS=linux GOARCH=$(ARCH) $(GOBUILD) -o ./sws-binaries/linux_$(ARCH)/$(BINARY_NAME) -ldflags $(LDFLAGS) -v
	@cp -r ./node_config ./sws-binaries/linux_$(ARCH)/
	@zip -r ./sws-binaries/linux_$(ARCH)-v$(VERSION).zip ./sws-binaries/linux_$(ARCH)
build-darwin:
	@mkdir -p ./sws-binaries/darwin_$(ARCH)
	GOOS=darwin GOARCH=$(ARCH) $(GOBUILD) -o ./sws-binaries/darwin_$(ARCH)/$(BINARY_NAME) -ldflags $(LDFLAGS) -v
	@cp -r ./node_config ./sws-binaries/darwin_$(ARCH)/
	@zip -r ./sws-binaries/darwin_$(ARCH)-v$(VERSION).zip ./sws-binaries/darwin_$(ARCH)
build-windows:
	@mkdir -p ./sws-binaries/windows_$(ARCH)
	GOOS=windows GOARCH=$(ARCH) $(GOBUILD) -o ./sws-binaries/windows_$(ARCH)/$(BINARY_NAME).exe -ldflags $(LDFLAGS) -v
	@cp -r ./node_config ./sws-binaries/windows_$(ARCH)/
	@zip -r ./sws-binaries/windows_$(ARCH)-v$(VERSION).zip ./sws-binaries/windows_$(ARCH)

version:
	@echo $(VERSION)