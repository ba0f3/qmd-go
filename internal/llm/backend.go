package llm

import (
	"os"
	"strings"
)

// isGGUFSpec returns true if model looks like a GGUF spec: local path ending in .gguf,
// "repo:file.gguf" (Hugging Face), or "org/repo/file.gguf".
func isGGUFSpec(model string) bool {
	model = strings.TrimSpace(model)
	if model == "" {
		return false
	}
	if strings.HasSuffix(model, ".gguf") {
		return true
	}
	if idx := strings.Index(model, ":"); idx >= 0 {
		file := model[idx+1:]
		return strings.TrimSpace(file) != "" && strings.HasSuffix(file, ".gguf")
	}
	return false
}

// NewEmbedClient returns an LLM implementation for embeddings.
// If QMD_EMBED_BACKEND=gguf or model is a GGUF spec (path to .gguf, or "repo:file.gguf"),
// a GGUF client is used (requires build with -tags gguf). Otherwise uses Ollama/OpenAI-compatible API.
func NewEmbedClient(model string) (LLM, error) {
	useGGUF := os.Getenv("QMD_EMBED_BACKEND") == "gguf" || isGGUFSpec(model)
	if useGGUF {
		return newGGUFClient(model)
	}
	baseURL := os.Getenv("OLLAMA_HOST")
	if baseURL == "" {
		baseURL = "http://localhost:11434/v1"
	}
	return NewOpenAIClient(baseURL, model), nil
}
