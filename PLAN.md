# reckon

> SessionStart hook that builds a compact symbol index of the current codebase and injects a navigation hint into context, so Claude always checks the index before searching files.

---

## Problem

When Claude needs to find where something lives in a codebase — the login flow, the email sender, the auth middleware — it typically greps or globs across the entire project. This is:

- **Slow** — searches through hundreds of files
- **Token-expensive** — intermediate results consume context
- **Imprecise** — wrong file names or conventions cause misses

This tool builds a flat symbol index (`.codeindex` or similar) at session start. Claude checks it first, gets an exact file + line, then reads only that file. Falls back to normal search only if the symbol isn't in the index.

No CLAUDE.md edits required — the SessionStart hook output itself instructs Claude how to use the index.

---

## How it works

```
SessionStart
      │
      ▼
  reckon
      │
      ├── detect languages in CWD (by file extension)
      │
      ├── extract symbols from every source file
      │       └── for each file: class, method, function, interface, enum, type
      │
      ├── write .codeindex (tab-separated, one symbol per line)
      │
      └── stdout: hook response telling Claude the index exists + how to use it
```

### Hook output (stdout)
```json
{
  "action": "continue",
  "output": "Symbol index rebuilt: .codeindex — 1,243 symbols across 87 files.\nGrep .codeindex before searching the codebase. If not found there, search normally."
}
```

No type bodies, no content injection — just a pointer to the index file.

---

## Index Format

Tab-separated, one symbol per line:

```
LoginAsync          src/controllers/UserController.cs    145    method
ValidateToken       src/services/AuthService.cs           89    method
AuthService         src/services/AuthService.cs            1    class
LoginDto            src/models/LoginDto.cs                 1    class
useAuthStore        src/stores/auth.store.ts              12    function
AuthGuard           src/guards/auth.guard.ts               1    class
```

**Columns:** `symbol | relative/path/to/file | line | kind`

- Paths are CWD-relative
- `kind`: `class`, `method`, `function`, `interface`, `type`, `enum`, `struct`, `const`
- One file: `.codeindex` at project root
- Gitignore it

### Claude's workflow
1. Grep `.codeindex` for the concept (e.g. `login`, case-insensitive)
2. Get exact file + line in one shot
3. Read only that file
4. If nothing found → fall back to normal codebase search

---

## Language Support

Extraction is purely regex-based. No compilation, no AST, no language server.

Each language is a small extractor implementing:

```go
type Extractor interface {
    Extensions() []string
    Extract(file string, content []byte) []Symbol
}
```

**Ships with:**
| Language | Extensions |
|---|---|
| TypeScript / JavaScript | `.ts`, `.tsx`, `.js`, `.jsx` |
| C# | `.cs` |
| Go | `.go` |
| Python | `.py` |
| Java | `.java` |
| Rust | `.rs` |
| Ruby | `.rb` |
| PHP | `.php` |

Unknown extension → skip. New language → add one file, register extractor.

---

## CLI Surface

```bash
reckon                  # SessionStart mode: rebuild index, emit hook response
reckon index            # rebuild .codeindex only (no hook output), print stats
reckon stats            # show symbol count, file count, language breakdown
reckon version
```

`reckon index` is the on-demand rebuild command users run manually or wire to a git hook.

---

## Config

Optional `.codeindex.yml` at project root:

```yaml
skip_patterns:
  - "**/*.test.ts"
  - "**/*.spec.ts"
  - "**/__mocks__/**"
  - "*.generated.ts"
  - "**/node_modules/**"
  - "**/bin/**"
  - "**/obj/**"
```

Sensible defaults — most projects need zero config.

---

## Phases

### Phase 1 — Core
- [ ] Symbol struct + Extractor interface
- [ ] Index writer (tab-separated, `.codeindex`)
- [ ] File walker with default skip patterns
- [ ] SessionStart hook response
- [ ] `reckon index` CLI command

### Phase 2 — Extractors
- [ ] TypeScript/JavaScript extractor
- [ ] C# extractor
- [ ] Go extractor
- [ ] Python extractor

### Phase 3 — More extractors + config
- [ ] Java, Rust, Ruby, PHP extractors
- [ ] `.codeindex.yml` skip_patterns
- [ ] `reckon stats` command

### Phase 4 — Polish
- [ ] GoReleaser + Homebrew tap
- [ ] README with install + hook registration instructions

---

## Hook Registration

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [{ "type": "command", "command": "reckon" }]
      }
    ]
  }
}
```

Run before `ctx restore` so the index pointer lands in context before session state is restored.

---

## Non-goals

- No type body injection (that was the old design)
- No cross-file type resolution
- No semantic aliases (may revisit later)
- No caching between runs — always rebuild at SessionStart (fast enough with regex)
- No compilation or language server integration
- No MCP server
