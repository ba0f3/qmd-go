# QMD Makefile
# Regular build (no GGUF): make build
# Install: make install [PREFIX=/usr/local]
# GGUF build: make deps-gguf && make build-gguf

GO ?= go
PREFIX ?= ~/.local
DEPS_DIR = .deps
LLAMA_CPP = $(DEPS_DIR)/go-llama.cpp

.PHONY: build build-gguf deps-gguf deps-purego build-purego install clean update-llama-cpp update-llama-cpp-ggml

build:
	$(GO) build -tags fts5 -o qmd ./cmd/qmd

# Clone go-llama.cpp with submodules and build libbinding.a.
# Required once before build-gguf (go-llama.cpp needs the llama.cpp submodule for C++ headers).
# Note: The vendored llama.cpp may be old; if models fail to load, try update-llama-cpp.
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

# Restore llama.cpp submodule to the commit expected by go-llama.cpp and rebuild.
# Use this to fix a broken build after the submodule was changed (e.g. by a failed update).
# Note: "Updating" to origin/master usually breaks the build (grammar-parser.h removed in newer llama.cpp).
# The submodule is reset to the commit recorded in go-llama.cpp's tree.
update-llama-cpp:
	@if [ ! -d "$(LLAMA_CPP)/.git" ]; then \
		echo "Error: go-llama.cpp not cloned. Run 'make deps-gguf' first."; \
		exit 1; \
	fi
	@COMMIT=$$(cd $(LLAMA_CPP) && git rev-parse HEAD:llama.cpp 2>/dev/null); \
	if [ -z "$$COMMIT" ]; then \
		echo "Error: could not get llama.cpp submodule commit from go-llama.cpp."; \
		exit 1; \
	fi; \
	echo "Restoring llama.cpp submodule to go-llama.cpp expected commit $$COMMIT..."; \
	cd $(LLAMA_CPP)/llama.cpp && git fetch origin && git checkout $$COMMIT
	@echo "Normalizing line endings so patches apply..."
	@(cd $(LLAMA_CPP) && find . -type f \( -name '*.c' -o -name '*.cpp' -o -name '*.h' -o -name '*.patch' \) -print0 | xargs -0 -r sed -i 's/\r$$//' 2>/dev/null || true)
	@cd $(LLAMA_CPP) && git checkout HEAD -- Makefile 2>/dev/null || true
	@cd $(LLAMA_CPP) && rm -f prepare && $(MAKE) clean 2>/dev/null || true
	@cd $(LLAMA_CPP) && $(MAKE) libbinding.a

# Switch llama.cpp submodule to ggml-org/llama.cpp (has nomic_bert and other embedding architectures).
# WARNING: This may fail due to API incompatibilities between go-llama.cpp and newer llama.cpp versions.
# If build fails with errors like "load_binding_model was not declared" or "llama_binding_state was not declared",
# go-llama.cpp's binding.cpp is incompatible with this llama.cpp version. Use the default submodule instead.
# Discards local changes in the submodule and rebuilds libbinding.a.
# Note: ggml-org/llama.cpp has ggml headers in ggml/include/, so we patch CXXFLAGS to add that path.
update-llama-cpp-ggml:
	@if [ ! -d "$(LLAMA_CPP)/.git" ]; then \
		echo "Error: go-llama.cpp not cloned. Run 'make deps-gguf' first."; \
		exit 1; \
	fi
	@cd $(LLAMA_CPP)/llama.cpp && \
		git remote add ggml https://github.com/ggml-org/llama.cpp 2>/dev/null || true && \
		git fetch ggml && \
		git reset --hard ggml/master
	@echo "Patching go-llama.cpp Makefile to add include paths for ggml-org structure..."
	@cd $(LLAMA_CPP) && \
		if ! grep -q "llama.cpp/include" Makefile; then \
			sed -i 's|CXXFLAGS = -I./llama.cpp|CXXFLAGS = -I./llama.cpp -I./llama.cpp/include -I./llama.cpp/ggml/include|' Makefile; \
			sed -i 's|CFLAGS   = -I./llama.cpp|CFLAGS   = -I./llama.cpp -I./llama.cpp/include -I./llama.cpp/ggml/include|' Makefile; \
		fi && \
		if ! grep -q "llama.cpp/ggml/include" Makefile || grep -q "llama.cpp/ggml/include.*llama.cpp/ggml/include" Makefile; then \
			sed -i 's|-I./llama.cpp/ggml/include -I./llama.cpp/ggml/include|-I./llama.cpp/ggml/include|g' Makefile; \
		fi
	@cd $(LLAMA_CPP) && $(MAKE) clean 2>/dev/null || true
	@cd $(LLAMA_CPP) && $(MAKE) libbinding.a

# Clone llama.cpp submodule for purego build (kelindar/search method).
deps-purego:
	@mkdir -p llama-go
	@if [ ! -d "llama-go/llama.cpp/.git" ]; then \
		cd llama-go && git clone --depth 1 https://github.com/ggerganov/llama.cpp.git; \
	fi

# Build llama_go shared library (purego method, no CGO).
# Output: llama-go/build/libllama_go.so (Linux) or llama-go/build/llama_go.dll (Windows)
build-purego: deps-purego
	@mkdir -p llama-go/build
	@cd llama-go/build && \
		cmake -DBUILD_SHARED_LIBS=ON -DCMAKE_BUILD_TYPE=Release .. && \
		cmake --build . --config Release
	@echo "Built shared library: llama-go/build/libllama_go.so (or .dll/.dylib)"

# Build QMD with GGUF support. Tries purego first, falls back to go-llama.cpp (CGO).
# Optional: LDFLAGS="-s -w" for release builds (strip debug symbols).
build-gguf: deps-gguf
	$(GO) mod edit -replace=github.com/go-skynet/go-llama.cpp=./$(LLAMA_CPP)
ifeq ($(LDFLAGS),)
	@$(GO) build -tags gguf,fts5 -o qmd ./cmd/qmd
else
	@$(GO) build -tags gguf,fts5 -ldflags "$(LDFLAGS)" -o qmd ./cmd/qmd
endif

# Install qmd to PREFIX/bin (default /usr/local). Run after make build or make build-gguf.
install: build-gguf
	@mkdir -p $(PREFIX)/bin
	@cp -f qmd $(PREFIX)/bin/qmd
	@echo "installed $(PREFIX)/bin/qmd"

clean:
	rm -f qmd
	$(GO) mod edit -dropreplace=github.com/go-skynet/go-llama.cpp 2>/dev/null || true
