ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
WORKSPACE = $(ROOT)/workspace

## No need to change anything below this line (normally, unless you know what you are doing)
EXECUTABLES = go install
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))
SPLIT = $(subst -, ,$@)
GAME = $(word 2, $(SPLIT))

build:
	go build -o $(ROOT)/bin/persianbgbot ./cmd/persianbgbot

embed:
	go build -tags embed -o $(ROOT)/bin/persianbgbot ./cmd/persianbgbot

install: embed
	install $(ROOT)/bin/persianbgbot /usr/local/bin/persianbgbot
