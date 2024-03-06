PKG_PATH = "github.com/containers/prometheus-podman-exporter"
BIN := ./bin
GO := go
TARGET := prometheus-podman-exporter
COVERAGE_PATH ?= .coverage
GINKO_CLI_VERSION = $(shell grep 'ginkgo/v2' go.mod | grep -o ' v.*' | sed 's/ //g')
GOBIN := $(shell $(GO) env GOBIN)
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PRE_COMMIT = $(shell command -v bin/venv/bin/pre-commit ~/.local/bin/pre-commit pre-commit | head -n1)
PKG_MANAGER ?= $(shell command -v dnf yum|head -n1)
SELINUXOPT ?= $(shell test -x /usr/sbin/selinuxenabled && selinuxenabled && echo -Z)
BUILDFLAGS := -mod=vendor
BUILDTAGS ?= \
	$(shell hack/systemd_tag.sh) \
	$(shell hack/btrfs_installed_tag.sh) \
	$(shell hack/btrfs_tag.sh)

VERSION = $(shell cat VERSION  | grep VERSION | cut -d'=' -f2)
REVISION = $(shell cat VERSION  | grep REVISION | cut -d'=' -f2)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)

#=================================================
# Build binary, clean, install and uninstall
#=================================================

all: binary

.PHONY: clean
clean:
	@rm -rf $(BIN)

.PHONY: binary
binary: $(TARGET)  ## Build prometheus-podman-exporter binary
	@true

.PHONY: binary-remote
binary-remote: ## Build prometheus-podman-exporter for remote connection
	@mkdir -p $(BIN)/remote
	@export CGO_ENABLED=0 && $(GO) build $(BUILDFLAGS) --tags "remote containers_image_openpgp" -ldflags="-X '$(PKG_PATH)/cmd.buildVersion=$(VERSION)' -X '$(PKG_PATH)/cmd.buildRevision=$(REVISION)' -X '$(PKG_PATH)/cmd.buildBranch=$(BRANCH)'" -o $(BIN)/remote/$(TARGET)

.PHONY: binary-remote-amd64
binary-remote-amd64: ## Build amd64 prometheus-podman-exporter for remote connection
	@mkdir -p $(BIN)/remote
	@echo "building amd64"
	@export CGO_ENABLED=0 && GOARCH=amd64 $(GO) build $(BUILDFLAGS) --tags "remote containers_image_openpgp" -ldflags="-X '$(PKG_PATH)/cmd.buildVersion=$(VERSION)' -X '$(PKG_PATH)/cmd.buildRevision=$(REVISION)' -X '$(PKG_PATH)/cmd.buildBranch=$(BRANCH)'" -o $(BIN)/remote/$(TARGET)-amd64

.PHONY: binary-remote-s390x
binary-remote-s390x: ## Build s390x prometheus-podman-exporter for remote connection
	@mkdir -p $(BIN)/remote
	@echo "building s390x"
	@export CGO_ENABLED=0 && GOARCH=s390x $(GO) build $(BUILDFLAGS) --tags "remote containers_image_openpgp" -ldflags="-X '$(PKG_PATH)/cmd.buildVersion=$(VERSION)' -X '$(PKG_PATH)/cmd.buildRevision=$(REVISION)' -X '$(PKG_PATH)/cmd.buildBranch=$(BRANCH)'" -o $(BIN)/remote/$(TARGET)-s390x

.PHONY: binary-remote-ppc64le
binary-remote-ppc64le: ## Build ppc64le prometheus-podman-exporter for remote connection
	@mkdir -p $(BIN)/remote
	@echo "building ppc64le"
	@export CGO_ENABLED=0 && GOARCH=ppc64le $(GO) build $(BUILDFLAGS) --tags "remote containers_image_openpgp" -ldflags="-X '$(PKG_PATH)/cmd.buildVersion=$(VERSION)' -X '$(PKG_PATH)/cmd.buildRevision=$(REVISION)' -X '$(PKG_PATH)/cmd.buildBranch=$(BRANCH)'" -o $(BIN)/remote/$(TARGET)-ppc64le

.PHONY: binary-remote-arm64
binary-remote-arm64: ## Build arm64 prometheus-podman-exporter for remote connection
	@mkdir -p $(BIN)/remote
	@echo "building arm64"
	@export CGO_ENABLED=0 && GOARCH=arm64 $(GO) build $(BUILDFLAGS) --tags "remote containers_image_openpgp" -ldflags="-X '$(PKG_PATH)/cmd.buildVersion=$(VERSION)' -X '$(PKG_PATH)/cmd.buildRevision=$(REVISION)' -X '$(PKG_PATH)/cmd.buildBranch=$(BRANCH)'" -o $(BIN)/remote/$(TARGET)-arm64

.PHONY: $(TARGET)
$(TARGET): $(SRC)
	@echo "running go build"
	@mkdir -p $(BIN)
	$(GO) build $(BUILDFLAGS) -tags "$(BUILDTAGS)" -ldflags="-X '$(PKG_PATH)/cmd.buildVersion=$(VERSION)' -X '$(PKG_PATH)/cmd.buildRevision=$(REVISION)' -X '$(PKG_PATH)/cmd.buildBranch=$(BRANCH)'" -o $(BIN)/$(TARGET)

.PHONY: install
install:    ## Install prometheus-podman-exporter binary
	@install ${SELINUXOPT} -D -m0755 $(BIN)/$(TARGET) $(DESTDIR)/$(TARGET)

.PHONY: uninstall
uninstall:  ## Uninstall prometheus-podman-exporter binary
	@rm -f $(DESTDIR)/$(TARGET)

#=================================================
# Required tools installation tartgets
#=================================================

.PHONY: install.tools
install.tools: .install.pre-commit .install.codespell .install.golangci-lint .install.ginkgo ## Install needed tools

.PHONY: .install.codespell
.install.codespell:
	sudo ${PKG_MANAGER} -y install codespell

.PHONY: .install.pre-commit
.install.pre-commit:
	if [ -z "$(PRE_COMMIT)" ]; then \
		python3 -m pip install --user pre-commit; \
	fi

.PHONY: .install.ginkgo
.install.ginkgo:
	if [ ! -x "$(GOBIN)/ginkgo" ]; then \
		$(GO) install -mod=mod github.com/onsi/ginkgo/v2/ginkgo@$(GINKO_CLI_VERSION) ; \
	fi

.PHONY: .install.golangci-lint
.install.golangci-lint:
	VERSION=1.56.2 ./hack/install_golangci.sh

#=================================================
# Linting/Formatting/Code Validation targets
#=================================================

.PHONY: validate
validate: gofmt lint govet pre-commit codespell vendor ## Validate prometheus-podman-exporter code (fmt, lint, ...)

.PHONY: vendor
vendor: ## Check vendor
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
	@bash ./hack/tree_status.sh

.PHONY: lint
lint: ## Run golangci-lint
	@echo "running golangci-lint"
	$(BIN)/golangci-lint version
	$(BIN)/golangci-lint run

.PHONY: pre-commit
pre-commit:   ## Run pre-commit
ifeq ($(PRE_COMMIT),)
	@echo "FATAL: pre-commit was not found, make .install.pre-commit to installing it." >&2
	@exit 2
endif
	$(PRE_COMMIT) run -a

.PHONY: gofmt
gofmt:   ## Run gofmt
	@echo -e "gofmt check and fix"
	@gofmt -w $(SRC)

.PHONY: govet
govet:   ## Run govet
	@echo "running go vet"
	@go vet ../$(TARGET)

.PHONY: codespell
codespell: ## Run codespell
	@echo "running codespell"
	@codespell -S ./vendor,go.mod,go.sum,./.git


#=================================================
# Testing (units, functionality, ...) targets
#=================================================

.PHONY: test
test: ## Run tests
	rm -rf ${COVERAGE_PATH} && mkdir -p ${COVERAGE_PATH}
	$(GOBIN)/ginkgo \
		-r \
		--cover \
		--tags "$(BUILDTAGS) containers_image_openpgp" \
		--covermode atomic \
		--coverprofile coverprofile \
		--output-dir ${COVERAGE_PATH} \
		--succinct
	$(GO) tool cover -html=${COVERAGE_PATH}/coverprofile -o ${COVERAGE_PATH}/coverage.html
	$(GO) tool cover -func=${COVERAGE_PATH}/coverprofile > ${COVERAGE_PATH}/functions
	cat ${COVERAGE_PATH}/functions | sed -n 's/\(total:\).*\([0-9][0-9].[0-9]\)/\1 \2/p'

#=================================================
# Help menu
#=================================================

_HLP_TGTS_RX = '^[[:print:]]+:.*?\#\# .*$$'
_HLP_TGTS_CMD = grep -E $(_HLP_TGTS_RX) $(MAKEFILE_LIST)
_HLP_TGTS_LEN = $(shell $(_HLP_TGTS_CMD) | cut -d : -f 1 | wc -L)
_HLPFMT = "%-$(_HLP_TGTS_LEN)s %s\n"
.PHONY: help
help: ## Print listing of key targets with their descriptions
	@printf $(_HLPFMT) "Target:" "Description:"
	@printf $(_HLPFMT) "--------------" "--------------------"
	@$(_HLP_TGTS_CMD) | sort | \
		awk 'BEGIN {FS = ":(.*)?## "}; \
			{printf $(_HLPFMT), $$1, $$2}'
