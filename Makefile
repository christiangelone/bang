SHELL=/bin/sh
include Makefile.*

.PHONY: fmt
fmt:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Formatting code..."
	@echo "----------------------------------------------------------------"
	$(GO) fmt ./...
	$(GOMOD) tidy

.PHONY: lint
lint:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Linting code..."
	@echo "----------------------------------------------------------------"
	$(GOLINT) run

.PHONY: test
test:
	@echo "----------------------------------------------------------------"
	@echo " ✅  Testing code..."
	@echo "----------------------------------------------------------------"
	$(GO) test ./... -v -coverprofile=coverage.out

.PHONY: coverage
coverage:
	@echo "----------------------------------------------------------------"
	@echo " 📊  Checking coverage..."
	@echo "----------------------------------------------------------------"
	$(GOTOOL) cover -html=coverage.out -o coverage.html
	$(GOTOOL) cover -func=coverage.out

.PHONY: compile
compile:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Compiling code..."
	@echo "----------------------------------------------------------------"
	$(GOBUILD) ./...
	$(PROTOTOOL) compile

.PHONY: deps
deps:
	@echo "----------------------------------------------------------------"
	@echo " ⬇️  Downloading dependencies..."
	@echo "----------------------------------------------------------------"
	$(GOGET) ./...

.PHONY: build
build: deps fmt
	@echo "----------------------------------------------------------------"
	@echo " 📦 Building binary..."
	@echo "----------------------------------------------------------------"
	$(GOBUILD) -a -ldflags "-w" -tags 'netgo osusergo' -o bang main.go

.PHONY: release
release:
	if [ -z "$(TAG)" ]; then echo 'You need to pass a TAG ❌. ex: make release TAG=v0.0.1' && exit 1; fi
	@echo "----------------------------------------------------------------"
	@echo " 📦 Tagging & Building binary..."
	@echo "----------------------------------------------------------------"
	@echo "Tag: \033[1;33m$(TAG)\033[0m"
	git tag $(TAG)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -ldflags "-w -X 'main.Version=$(TAG)'" -tags 'netgo osusergo' -o bang_darwin_x86_64 main.go
	GOOS=linux GOARCH=amd64 $(GOBUILD) -a -ldflags "-w -X 'main.Version=$(TAG)'" -tags 'netgo osusergo' -o bang_linux_x86_64 main.go


.PHONY: all
all: lint test build

################################################################################