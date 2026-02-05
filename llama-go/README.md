# llama-go: Purego wrapper for llama.cpp

This directory contains a C++ wrapper around llama.cpp that exposes a simple C API for embeddings, compatible with Go's purego (no CGO required).

## Building

```bash
# Clone llama.cpp submodule
make deps-purego

# Build shared library
make build-purego
```

This creates `build/libllama_go.so` (Linux), `build/llama_go.dll` (Windows), or `build/libllama_go.dylib` (macOS).

## Usage

The Go code automatically finds the library in:
1. `LLAMA_GO_LIB` environment variable
2. `llama-go/build/` directory (development)
3. Current directory
4. System library paths (`/usr/lib`, `/usr/local/lib`)
5. Directory containing the executable

## API

See `llama_go.h` for the C API:
- `llama_go_load()` - Load a GGUF model
- `llama_go_embed()` - Generate embedding for text
- `llama_go_free()` - Free model handle
- `llama_go_get_error()` - Get last error message
