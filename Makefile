# QMD Makefile
# Regular build (no GGUF): make build
# Install: make install [PREFIX=/usr/local]
# GGUF build: make deps-gguf && make build-gguf

GO ?= go
PREFIX ?= ~/.local
DEPS_DIR = .deps
LLAMA_CPP = $(DEPS_DIR)/go-llama.cpp

.PHONY: build build-gguf deps-gguf install clean

build:
	$(GO) build -tags fts5 -o qmd ./cmd/qmd

# Clone go-llama.cpp with submodules and build libbinding.a.
# Required once before build-gguf (go-llama.cpp needs the llama.cpp submodule for C++ headers).
deps-gguf:
	@mkdir -p $(DEPS_DIR)
	@if [ ! -d "$(LLAMA_CPP)/.git" ]; then \
		git clone --recurse-submodules https://github.com/go-skynet/go-llama.cpp $(LLAMA_CPP); \
	else \
		(cd $(LLAMA_CPP) && git submodule update --init --recursive); \
	fi
	@echo "Normalizing line endings so go-llama.cpp patches apply..."
	@(cd $(LLAMA_CPP) && find . -type f \( -name '*.c' -o -name '*.cpp' -o -name '*.h' -o -name '*.patch' \) -print0 | xargs -0 -r sed -i 's/\r$$//' 2>/dev/null || true)
	@cd $(LLAMA_CPP) && $(MAKE) libbinding.a

# Build QMD with GGUF support. Runs deps-gguf if needed and points Go at the local go-llama.cpp.
build-gguf: deps-gguf
	$(GO) mod edit -replace=github.com/go-skynet/go-llama.cpp=./$(LLAMA_CPP)
	@$(GO) build -tags gguf,fts5 -o qmd ./cmd/qmd

# Install qmd to PREFIX/bin (default /usr/local). Run after make build or make build-gguf.
install: build-gguf
	@mkdir -p $(PREFIX)/bin
	@cp -f qmd $(PREFIX)/bin/qmd
	@echo "installed $(PREFIX)/bin/qmd"

clean:
	rm -f qmd
	$(GO) mod edit -dropreplace=github.com/go-skynet/go-llama.cpp 2>/dev/null || true
