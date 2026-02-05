// C API for llama.cpp embeddings (purego-compatible)
#ifndef LLAMA_GO_H
#define LLAMA_GO_H

#ifdef __cplusplus
extern "C" {
#endif

// Handle for a loaded model
typedef void* LlamaGoModel;

// Load a GGUF model from file path. Returns NULL on error.
// Use llama_go_get_error() to get error message.
LlamaGoModel llama_go_load(const char* model_path, int n_ctx, int n_gpu_layers);

// Free a model handle
void llama_go_free(LlamaGoModel model);

// Generate embedding for text. Returns number of dimensions, or -1 on error.
// embedding must be pre-allocated with at least max_dims floats.
int llama_go_embed(LlamaGoModel model, const char* text, float* embedding, int max_dims);

// Get last error message (thread-local)
const char* llama_go_get_error(void);

#ifdef __cplusplus
}
#endif

#endif // LLAMA_GO_H
