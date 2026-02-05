# QMD - Query Markup Documents

An on-device search engine for everything you need to remember. Index your markdown notes, meeting transcripts, documentation, and knowledge bases. Search with keywords or natural language. Ideal for your agentic flows.

QMD combines BM25 full-text search, vector semantic search, and LLM re-ranking—all running locally via node-llama-cpp with GGUF models.

**This project is a fork of [tobi/qmd](https://github.com/tobi/qmd)** (Go port, optional GGUF embedding backend).

## Quick Start

```sh
# Install (requires Go 1.23+)
go install -tags fts5 github.com/ba0f3/qmd-go
# Or from repo:
# git clone https://github.com/ba0f3/qmd-go && cd qmd && go build -tags fts5 -o qmd-go ./cmd/qmd

# Create collections for your notes, docs, and meeting transcripts
qmd collection add ~/notes --name notes
qmd collection add ~/Documents/meetings --name meetings
qmd collection add ~/work/docs --name docs

# Add context to help with search results
qmd context add qmd://notes "Personal notes and ideas"
qmd context add qmd://meetings "Meeting transcripts and notes"
qmd context add qmd://docs "Work documentation"

# Generate embeddings for semantic search
qmd embed

# Search across everything
qmd search "project timeline"           # Fast keyword search
qmd vsearch "how to deploy"             # Semantic search
qmd query "quarterly planning process"  # Hybrid + reranking (best quality)

# Get a specific document
qmd get "meetings/2024-01-15.md"

# Get a document by docid (shown in search results)
qmd get "#abc123"

# Get multiple documents by glob pattern
qmd multi-get "journals/2025-05*.md"

# Search within a specific collection
qmd search "API" -c notes

# Export all matches for an agent
qmd search "API" --all --files --min-score 0.3
```

### Using with AI Agents

QMD's `--json` and `--files` output formats are designed for agentic workflows:

```sh
# Get structured results for an LLM
qmd search "authentication" --json -n 10

# List all relevant files above a threshold
qmd query "error handling" --all --files --min-score 0.4

# Retrieve full document content
qmd get "docs/api-reference.md" --full
```

### MCP Server

QMD exposes an MCP (Model Context Protocol) server for integration with Claude, Cursor, and other AI tools.

**Tools exposed:**
- `search` – Fast BM25 keyword search (supports collection filter)
- `vsearch` – Semantic vector search (supports collection filter)
- `query` – Hybrid search (BM25 + vector, RRF)
- `get` – Retrieve document by path or docid
- `multi_get` – Retrieve multiple documents by glob or list
- `status` – Index health and collection info

**Resources:** Documents are readable via `qmd://` URIs (e.g. `qmd://collection/path/to/file.md`).

**Claude Desktop** (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "qmd": {
      "command": "qmd",
      "args": ["mcp"]
    }
  }
}
```

**Claude Code** — Install the plugin (recommended):

```bash
claude marketplace add tobi/qmd
claude plugin add qmd@qmd
```

Or configure MCP manually in `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "qmd": {
      "command": "qmd",
      "args": ["mcp"]
    }
  }
}
```

## Architecture

- **SQLite FTS5** – Full-text search (BM25)
- **Embeddings** – Stored in SQLite; generated via Ollama or any OpenAI-compatible embedding API
- **Hybrid search** – `query` combines BM25 and vector results using Reciprocal Rank Fusion (RRF)
- **Chunking** – 800 tokens per chunk, 15% overlap (character-based in Go)
- **Index** – `~/.cache/qmd/index.sqlite` (or `INDEX_PATH`)

```
Query ──┬──► BM25 (FTS5) ──► Ranked list
        │
        └──► Embed query ──► Vector similarity ──► Ranked list
                                    │
                                    ▼
                            RRF fusion ──► Final results
```

## Requirements

- **Go** 1.23 or later (for building)
- **SQLite** with FTS5 (e.g. `github.com/mattn/go-sqlite3` – enable via build tag `fts5`)
- **Embeddings**: [Ollama](https://ollama.ai) (default) or any OpenAI-compatible API (set `OLLAMA_HOST` or use `OPENAI_API_BASE` + `OPENAI_API_KEY` for embed)

### Build tags

| Tag   | Purpose |
|-------|---------|
| `fts5` | **Use for normal builds.** Enables SQLite FTS5 full-text search. Without it, `qmd search` and related features will not work. |
| `gguf` | Optional. Enables local/Hugging Face GGUF embedding models (no Ollama required). Requires CGO and [go-llama.cpp](https://github.com/go-skynet/go-llama.cpp); use `make deps-gguf` then `make build-gguf`. |

**Examples:** `make build` (outputs `qmd-go`) or `go build -tags fts5 -o qmd-go ./cmd/qmd`. For GGUF embeddings: `make build-gguf`.

## Installation

**From source (recommended)**

```sh
git clone https://github.com/tobi/qmd
cd qmd
make build
make install
# Installs to /usr/local/bin/qmd; ensure /usr/local/bin is in your PATH
```

To install to a different prefix (e.g. your user directory):

```sh
make install PREFIX=$HOME/.local
# Then add $HOME/.local/bin to your PATH
```

**Using Go install**

```sh
go install -tags fts5 github.com/ba0f3/qmd-go@latest
# Binary is installed as qmd-go in $GOPATH/bin or $HOME/go/bin; add that to your PATH
```

### Development

```sh
git clone https://github.com/tobi/qmd
cd qmd
./qmd --help          # Runs via go run
make build && ./qmd-go status
go test ./...
```

## Usage

### Collection Management

```sh
# Create a collection from current directory
qmd collection add . --name myproject

# Create a collection with explicit path and custom glob mask
qmd collection add ~/Documents/notes --name notes --mask "**/*.md"

# List all collections
qmd collection list

# Remove a collection
qmd collection remove myproject

# Rename a collection
qmd collection rename myproject my-project

# List files in a collection
qmd ls notes
qmd ls notes/subfolder
```

### Generate Vector Embeddings

```sh
# Embed all indexed documents (Ollama or OpenAI-compatible API)
qmd embed

# Force re-embed everything
qmd embed -f
```

Set `QMD_EMBED_MODEL` (default: `nomic-embed-text` for Ollama) or use an OpenAI-compatible endpoint.

### Context Management

Context adds descriptive metadata to collections and paths, helping search understand your content.

```sh
# Add context to a collection (using qmd:// virtual paths)
qmd context add qmd://notes "Personal notes and ideas"
qmd context add qmd://docs/api "API documentation"

# Add context from within a collection directory
cd ~/notes && qmd context add "Personal notes and ideas"

# Add global context (applies to all collections)
qmd context add / "Knowledge base for my projects"

# List all contexts
qmd context list

# Remove context
qmd context rm qmd://notes/old
```

### Search Commands

| Command   | Description                                      |
|----------|--------------------------------------------------|
| `search` | BM25 full-text search only                       |
| `vsearch` | Vector semantic search only                    |
| `query`  | Hybrid: BM25 + vector, RRF (no LLM reranker)     |

```sh
# Full-text search (fast, keyword-based)
qmd search "authentication flow"

# Vector search (semantic similarity)
qmd vsearch "how to login"

# Hybrid search (best quality without external reranker)
qmd query "user authentication"
```

### Options

```sh
# Search options
-n <num>           # Number of results (default: 5, or 20 for --files/--json)
-c, --collection   # Restrict search to a specific collection
--all              # Return all matches (use with --min-score to filter)
--min-score <num>  # Minimum score threshold (default: 0)
--full             # Show full document content
--line-numbers     # Add line numbers to output
--index <name>     # Use named index (default: index)

# Output formats (for search and multi-get)
--files            # Output: docid,score,filepath,context
--json             # JSON output with snippets
--csv              # CSV output
--md               # Markdown output
--xml              # XML output

# Get options
qmd get <file>[:line]  # Get document, optionally starting at line
-l <num>               # Maximum lines to return
--from <num>           # Start from line number

# Multi-get options
-l <num>           # Maximum lines per file
--max-bytes <num>  # Skip files larger than N bytes (default: 10KB)
```

### Index Maintenance

```sh
# Show index status and collections
qmd status

# Re-index all collections
qmd update

# Re-index with git pull first (for remote repos)
qmd update --pull

# Get document by filepath (with docid fallback)
qmd get notes/meeting.md

# Get document by docid (from search results)
qmd get "#abc123"

# Get document starting at line 50, max 100 lines
qmd get notes/meeting.md:50 -l 100

# Get multiple documents by glob pattern
qmd multi-get "journals/2025-05*.md"

# Get multiple documents by comma-separated list (supports docids)
qmd multi-get "doc1.md, doc2.md, #abc123"

# Limit multi-get to files under 20KB
qmd multi-get "docs/*.md" --max-bytes 20480

# Output multi-get as JSON for agent processing
qmd multi-get "docs/*.md" --json

# Clean up cache and orphaned data
qmd cleanup
```

## Data Storage

Index stored in: `~/.cache/qmd/index.sqlite` (or `INDEX_PATH`; use `--index <name>` for `~/.cache/qmd/<name>.sqlite`).

### Schema (overview)

- **documents** – Paths, titles, content hash, collection, active flag
- **content** – Full document text (keyed by hash)
- **documents_fts** – FTS5 full-text index
- **content_vectors** / **embedding_blobs** – Chunk embeddings for vector search
- Config (collections, context) – YAML in `~/.config/qmd/index.yml` (or per `--index`)

## Embedding backends

**Default (no build tag):** Embeddings use Ollama or any OpenAI-compatible API. Set `OLLAMA_HOST` and `QMD_EMBED_MODEL` (e.g. `nomic-embed-text`).

**GGUF (optional):** You can use local or Hugging Face GGUF embedding models without an external server.

1. Set `QMD_EMBED_BACKEND=gguf` **or** set `QMD_EMBED_MODEL` to a GGUF spec (see below).
2. Build QMD with GGUF support. [go-llama.cpp](https://github.com/go-skynet/go-llama.cpp) uses a git submodule for C++ code, so use the Makefile (one-time deps, then build):

   ```sh
   make deps-gguf   # clone go-llama.cpp + submodules, build libbinding.a (once)
   make build-gguf  # add replace in go.mod and build qmd with -tags gguf,fts5
   ```

   This clones go-llama.cpp into `.deps/go-llama.cpp` and runs `make libbinding.a` there, then builds qmd so the CGO compile finds the C++ headers and library.

3. **Model spec** for `QMD_EMBED_MODEL` when using GGUF:
   - **Local path:** `/path/to/embedding-model-Q8_0.gguf`
   - **Hugging Face (repo:file):** `ggml-org/embeddinggemma-300M-GGUF:embeddinggemma-300M-Q8_0.gguf`
   - **Hugging Face (path-style):** `ggml-org/embeddinggemma-300M-GGUF/embeddinggemma-300M-Q8_0.gguf`

   The first time you use a Hugging Face spec, the file is downloaded to `~/.cache/qmd/models` (or `QMD_MODEL_CACHE`).

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XDG_CACHE_HOME` | `~/.cache` | Cache directory (index SQLite) |
| `INDEX_PATH` | (derived) | Override index DB path |
| `OLLAMA_HOST` | `http://localhost:11434/v1` | Ollama API base for embed (API backend only) |
| `QMD_EMBED_MODEL` | `nomic-embed-text` | Embedding model name or GGUF path/spec |
| `QMD_EMBED_BACKEND` | (auto) | `gguf` to force GGUF backend (requires build with `-tags gguf`) |
| `QMD_MODEL_CACHE` | `~/.cache/qmd/models` | Directory for downloaded GGUF models |

For OpenAI-compatible APIs, set your provider’s base URL and API key (e.g. `OPENAI_API_BASE`, `OPENAI_API_KEY`); the Go CLI uses the same env names as typical OpenAI clients where applicable.

## License

MIT
