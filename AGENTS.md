# AGENTS.md – Guidance for AI Agents

This file orients AI agents (Cursor, Claude Code, etc.) working on the QMD codebase. For full command reference and constraints, see **CLAUDE.md**.

## Project overview

- **QMD** = on-device search over markdown/knowledge: BM25 (FTS5), optional vector search, optional LLM reranking.
- This repo is the **Go** implementation of QMD (no Bun/TypeScript).
- **Index**: SQLite at `~/.cache/qmd/index.sqlite` (or `INDEX_PATH`). Do not modify it directly.

## Where to change what

| Area | Location | Notes |
|------|----------|--------|
| Go CLI | `cmd/qmd/*.go` | Cobra commands |
| Go store/indexer | `internal/store/`, `internal/indexer/` | SQLite, FTS5, indexing |
| Eval harness | `test/eval_harness_test.go` | Search quality eval over `test/eval-docs/` |
| Finetune / training | `finetune/` | Python (uv); see `finetune/CLAUDE.md` |

## Constraints (must follow)

- **Do not run automatically**: `qmd collection add`, `qmd embed`, `qmd update`. Suggest commands for the user to run.
- **Do not** modify the SQLite index directly; use the app’s commands/APIs.
## Conventions

- Code and comments in **English**.
- If a binary is not on PATH, use `/home/linuxbrew/.linuxbrew/bin/` as prefix (e.g. for `go`).
- Task management: use the CLI ticket system (`tk help`) when relevant.

## Testing

- **Go**: `go test ./...` (from repo root). Store/indexer tests use a temp DB; FTS5 may be unavailable in some builds (tests skip when appropriate).
- **Eval harness**: `go test -v ./test/ -run TestEvalHarnessSearch` (uses store directly; skips if FTS5 not in sqlite3 build).

## References

- **CLAUDE.md** – Commands, options, architecture, and “do not” rules.
- **README.md** – User-facing quick start, MCP, and architecture.
- **finetune/CLAUDE.md** – Finetune/data pipeline specifics.
- **Upstream**: [github.com/tobi/qmd](https://github.com/tobi/qmd).
