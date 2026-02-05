// C++ wrapper for llama.cpp embeddings (purego-compatible)
#include "llama_go.h"
#include "llama.h"
#include <string>
#include <vector>
#include <mutex>
#include <cstring>

// Thread-local error storage
thread_local std::string g_last_error;

static void set_error(const char* msg) {
    g_last_error = msg ? msg : "unknown error";
}

extern "C" {

LlamaGoModel llama_go_load(const char* model_path, int n_ctx, int n_gpu_layers) {
    if (!model_path) {
        set_error("model_path is NULL");
        return nullptr;
    }

    try {
        llama_model_params model_params = llama_model_default_params();
        model_params.n_gpu_layers = n_gpu_layers;
        
        llama_model* model = llama_model_load_from_file(model_path, model_params);
        if (!model) {
            set_error("failed to load model");
            return nullptr;
        }

        // Store model pointer (we'll free it later)
        return reinterpret_cast<LlamaGoModel>(model);
    } catch (const std::exception& e) {
        set_error(e.what());
        return nullptr;
    } catch (...) {
        set_error("unknown exception during model load");
        return nullptr;
    }
}

void llama_go_free(LlamaGoModel model) {
    if (model) {
        llama_model* m = reinterpret_cast<llama_model*>(model);
        llama_model_free(m);
    }
}

int llama_go_embed(LlamaGoModel model, const char* text, float* embedding, int max_dims) {
    if (!model || !text || !embedding || max_dims <= 0) {
        set_error("invalid parameters");
        return -1;
    }

    try {
        llama_model* m = reinterpret_cast<llama_model*>(model);

        const llama_vocab* vocab = llama_model_get_vocab(m);
        if (!vocab) {
            set_error("failed to get vocab");
            return -1;
        }

        // Create context for embeddings
        llama_context_params ctx_params = llama_context_default_params();
        ctx_params.embeddings = true;  // Enable embeddings

        llama_context* ctx = llama_init_from_model(m, ctx_params);
        if (!ctx) {
            set_error("failed to create context");
            return -1;
        }

        // Tokenize input
        std::vector<llama_token> tokens;
        int n_tokens = llama_tokenize(vocab, text, (int32_t)strlen(text), nullptr, 0, true, false);
        if (n_tokens < 0) {
            llama_free(ctx);
            set_error("tokenization failed");
            return -1;
        }

        tokens.resize(n_tokens);
        n_tokens = llama_tokenize(vocab, text, (int32_t)strlen(text), tokens.data(), (int32_t)tokens.size(), true, false);
        if (n_tokens < 0) {
            llama_free(ctx);
            set_error("tokenization failed");
            return -1;
        }

        // Create batch: n_tokens, embd=0 (use token ids), n_seq_max=1
        llama_batch batch = llama_batch_init(n_tokens, 0, 1);
        batch.n_tokens = n_tokens;
        for (int i = 0; i < n_tokens; i++) {
            batch.token[i] = tokens[i];
            batch.pos[i] = i;
            batch.n_seq_id[i] = 1;
            batch.seq_id[i][0] = 0;
            batch.logits[i] = (i == n_tokens - 1) ? 1 : 0;  // only last token outputs embedding
        }

        // Evaluate tokens (batch by value)
        if (llama_decode(ctx, batch) < 0) {
            llama_batch_free(batch);
            llama_free(ctx);
            set_error("decode failed");
            return -1;
        }

        llama_batch_free(batch);

        // Get embedding (use model for n_embd; embeddings from context)
        int n_embd = llama_model_n_embd(m);
        const float* emb = llama_get_embeddings(ctx);
        if (!emb) {
            llama_free(ctx);
            set_error("failed to get embeddings");
            return -1;
        }

        // Copy embedding
        int copy_size = (n_embd < max_dims) ? n_embd : max_dims;
        std::memcpy(embedding, emb, copy_size * sizeof(float));

        llama_free(ctx);
        return n_embd;
    } catch (const std::exception& e) {
        set_error(e.what());
        return -1;
    } catch (...) {
        set_error("unknown exception during embedding");
        return -1;
    }
}

const char* llama_go_get_error(void) {
    return g_last_error.c_str();
}

} // extern "C"
