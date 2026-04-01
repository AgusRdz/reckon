# CLAUDE.md — reckon

## What this is

A Go CLI tool that hooks into `SessionStart` to build a compact symbol index (`.codeindex`) of the current project and inject a navigation hint into context before Claude's first turn. Claude greps the index to locate any function/class/method instantly, then reads only that file. Falls back to normal codebase search if not found.

## Critical constraints

- **SessionStart runs before the first user turn** — must be fast. Target: < 2 seconds total.
- **Index only, no content injection** — never inject type bodies or file contents. Only a pointer to `.codeindex`.
- **All languages via regex** — no compilation, no AST, no language server. Extraction is regex + string parsing only.
- **CWD-relative paths** — all paths in `.codeindex` are relative to CWD (project root). Never absolute.
- **Always rebuild** — no caching between sessions. Files may have changed.
- **Silent on empty** — if no symbols found, emit `{ "action": "continue", "output": "" }`. Never inject a "nothing found" message.

## Repository structure

```
reckon/
├── main.go
├── cmd/
│   ├── root.go           # SessionStart handler (no args): rebuild index, emit hook response
│   ├── index.go          # `reckon index` — rebuild on demand, print stats
│   └── stats.go          # `reckon stats` command
├── extract/
│   ├── typescript.go     # TS/JS class, function, interface, type, enum extraction
│   ├── csharp.go         # C# class, method, interface, enum extraction
│   ├── golang.go         # Go func, type, struct, interface extraction
│   └── python.go         # Python class, def extraction
├── index/
│   └── index.go          # write/read .codeindex
├── walk/
│   └── walk.go           # file walker with skip patterns
├── detect/
│   └── project.go        # detect languages by file extension
├── config/
│   └── config.go         # .codeindex.yml loader
├── .codeindex.yml
├── PLAN.md
└── CLAUDE.md
```

## Protocol

**stdin:** empty (SessionStart provides no tool_input)

**stdout (index built):**
```json
{
  "action": "continue",
  "output": "Symbol index rebuilt: .codeindex — 1,243 symbols across 87 files.\nGrep .codeindex before searching the codebase. If not found there, search normally."
}
```

**stdout (no symbols found):**
```json
{ "action": "continue", "output": "" }
```

**stderr:** debug logs, extraction errors — never reaches context.

## Index format

Tab-separated, one symbol per line:

```
LoginAsync    src/controllers/UserController.cs    145    method
AuthService   src/services/AuthService.cs            1    class
useAuthStore  src/stores/auth.store.ts              12    function
```

Columns: `symbol \t file \t line \t kind`

kinds: `class`, `method`, `function`, `interface`, `type`, `enum`, `struct`, `const`

## SessionStart hook ordering

Run `reckon` before `ctx restore`. Configure via hook array order in `settings.json`.

## What NOT to build

- No type body injection
- No cross-file type resolution
- No semantic aliases
- No MCP server
- No caching between sessions
- No import graph traversal
- No compilation or language server integration
