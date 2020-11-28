ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
WORKSPACE = $(ROOT)/workspace
GO_BIN_DATA = $(WORKSPACE)/bin/go-bindata

## No need to change anything below this line (normally, unless you know what you are doing)
EXECUTABLES = go install
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))
SPLIT = $(subst -, ,$@)
GAME = $(word 2, $(SPLIT))


$(GO_BIN_DATA):
	cd $(WORKSPACE) && go build github.com/go-bindata/go-bindata/go-bindata
	mv $(ROOT)/workspace/go-bindata $@

version: $(GO_BIN_DATA)
	$(GO_BIN_DATA) -version

data-%: $(GO_BIN_DATA)
	cd $(ROOT)/internal/$(GAME) && $(GO_BIN_DATA) -pkg $(GAME) data/

bin-data: data-bloodrage data-terraforming data-tinytown

build:
	go build -o $(ROOT)/bin/persianbgbot ./cmd/persianbgbot

install: build
	install $(ROOT)/bin/persianbgbot /usr/local/bin/persianbgbot
