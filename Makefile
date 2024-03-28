BIN:=kv
PACKAGE:=github.com/dcai/kv
GO?=go
GOOS?=$(word 1, $(subst /, " ", $(word 4, $(shell go version))))
MAKEFILE:=$(realpath $(lastword $(MAKEFILE_LIST)))
ROOT_DIR:= $(shell dirname $(MAKEFILE))
SOURCES:=$(wildcard *.go src/*.go src/*/*.go shell/*sh) $(MAKEFILE)
BUILD_FLAGS:=-a -ldflags "-s -w"

BINARY32       := kv-$(GOOS)_386
BINARY64       := kv-$(GOOS)_amd64
BINARYS390     := kv-$(GOOS)_s390x
BINARYARM5     := kv-$(GOOS)_arm5
BINARYARM6     := kv-$(GOOS)_arm6
BINARYARM7     := kv-$(GOOS)_arm7
BINARYARM8     := kv-$(GOOS)_arm8
BINARYPPC64LE  := kv-$(GOOS)_ppc64le
BINARYRISCV64  := kv-$(GOOS)_riscv64
BINARYLOONG64  := kv-$(GOOS)_loong64
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
	BINARY := $(BINARY64)
else ifeq ($(UNAME_M),amd64)
	BINARY := $(BINARY64)
else ifeq ($(UNAME_M),s390x)
	BINARY := $(BINARYS390)
else ifeq ($(UNAME_M),i686)
	BINARY := $(BINARY32)
else ifeq ($(UNAME_M),i386)
	BINARY := $(BINARY32)
else ifeq ($(UNAME_M),armv5l)
	BINARY := $(BINARYARM5)
else ifeq ($(UNAME_M),armv6l)
	BINARY := $(BINARYARM6)
else ifeq ($(UNAME_M),armv7l)
	BINARY := $(BINARYARM7)
else ifeq ($(UNAME_M),armv8l)
	# armv8l is always 32-bit and should implement the armv7 ISA, so
	# just use the same filename as for armv7.
	BINARY := $(BINARYARM7)
else ifeq ($(UNAME_M),arm64)
	BINARY := $(BINARYARM8)
else ifeq ($(UNAME_M),aarch64)
	BINARY := $(BINARYARM8)
else ifeq ($(UNAME_M),ppc64le)
	BINARY := $(BINARYPPC64LE)
else ifeq ($(UNAME_M),riscv64)
	BINARY := $(BINARYRISCV64)
else ifeq ($(UNAME_M),loongarch64)
	BINARY := $(BINARYLOONG64)
else
$(error Build on $(UNAME_M) is not supported, yet.)
endif

all: build
	@./kv get

build:
	@mkdir -p target
	@go build -o $(BIN) $(PACKAGE)/cmd/kv

vars:
	@echo $(ROOT_DIR)
	@echo $(GOOS)
	@echo $(GO)
	@echo $(SOURCES)

test-set: build
	@rm data.json
	./kv set NODE_ENV production
	./kv set testkey hello
	./kv set aws_key ssiikkuu_uuyy13
	./kv set jira_url https://fake_org.atlassian.net
	./kv set jira_user user@user.com
	./kv set jira_token aGVsbG8=
	@echo "=========================================="
	cat data.json | jq

test-getall: build
	@./kv get --all

test-getone: build
	@./kv get NODE_ENV

release:
	go build -ldflags "-w -s"

target/$(BINARYARM8): $(SOURCES)
	@echo "-----------" $@
	GOARCH=arm64 $(GO) build $(BUILD_FLAGS) -o $@

goreleaser-local:
	@goreleaser build --single-target --snapshot --clean
