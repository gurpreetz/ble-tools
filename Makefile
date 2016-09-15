#
# Confidential!!!
# Source code property of Blue Clover Design LLC.
#
# Demonstration, distribution, replication, or other use of the
# source codes is NOT permitted without prior written consent
# from Blue Clover Design.
#

GIT_DESC := $(shell git describe --tags --always --dirty --match "v[0-9]*")
VERSION_TAG := $(patsubst v%,%,$(GIT_DESC))

GO=go
GO_BUILD_OPTS=-ldflags "-X main.versionTag=$(VERSION_TAG)"
GO_GET=$(GO) get
GO_BUILD_OSX=GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_OPTS)

BINS :=
BINS += $(GOPATH)/bin/darwin_amd64/ble-tools


COMMON_DEPS := 
COMMON_DEPS += bleTools.go 
COMMON_DEPS += cmdLine.go 
COMMON_DEPS += xmlParser.go
COMMON_DEPS += csvParser.go

default: build

.PHONY: install
install: $(BINS)
	$(GO_GET) -d . && $(GO) install $(GO_BUILD_OPTS)

.PHONY: build
build: $(BINS)

$(GOPATH)/bin/darwin_amd64/ble-tools:
	@install -d $(GOPATH)
	$(GO_GET) -d . && $(GO_BUILD_OSX) -o $@ .

.PHONY: versions
versions:
	@echo "VERSION_TAG: $(VERSION_TAG)"

.PHONY: clean
clean:
	@rm -f $(BINS) 

