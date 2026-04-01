# reckon

> SessionStart hook that builds a compact symbol index of your codebase so Claude always knows where things live.

When Claude needs to find a function, class, or method, it normally greps across hundreds of files. reckon builds a flat `.codeindex` file at session start — Claude checks it first, gets an exact file and line number, then reads only that file.

## Install

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/AgusRdz/reckon/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/AgusRdz/reckon/main/install.ps1 | iex
```

Both scripts download the binary, add it to `PATH`, and run `reckon init` to register the Claude Code hook automatically.

To override the install directory:

```bash
RECKON_INSTALL_DIR=/usr/local/bin curl -fsSL https://raw.githubusercontent.com/AgusRdz/reckon/main/install.sh | sh
```

```powershell
$env:RECKON_INSTALL_DIR = "C:\tools\reckon"; irm https://raw.githubusercontent.com/AgusRdz/reckon/main/install.ps1 | iex
```

### Build from source

```bash
git clone https://github.com/AgusRdz/reckon.git
cd reckon
make install    # builds and copies to ~/.local/bin (Linux/macOS) or %LOCALAPPDATA%\Programs\reckon (Windows)
reckon init
```

## Hook registration

`reckon init` (run automatically by the install scripts) registers the SessionStart hook — no manual config needed. It writes the hook entry to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [{ "type": "command", "command": "/path/to/reckon" }]
      }
    ]
  }
}
```

If you use [ctx](https://github.com/AgusRdz/ctx) for session state, add reckon **before** ctx in the hook array so the index pointer lands in context first.

To remove the hook:

```bash
reckon uninstall
```

## How it works

On every session start, reckon:

1. Walks the project directory (respecting skip patterns)
2. Extracts symbols from every source file via regex
3. Writes `.codeindex` — a tab-separated file, one symbol per line
4. Emits a JSON hook response telling Claude the index exists

Claude then greps `.codeindex` before searching the codebase. If the symbol isn't there, it falls back to normal search.

### Index format

```
LoginAsync    src/controllers/UserController.cs    145    method
AuthService   src/services/AuthService.cs            1    class
useAuthStore  src/stores/auth.store.ts              12    function
```

Columns: `symbol · file (CWD-relative) · line · kind`

Kinds: `class`, `method`, `function`, `interface`, `type`, `enum`, `struct`, `const`

## Language support

| Language | Extensions |
|---|---|
| Go | `.go` |
| TypeScript / JavaScript | `.ts`, `.tsx`, `.js`, `.jsx` |
| C# | `.cs` |
| Python | `.py` |
| Java | `.java` |
| Rust | `.rs` |
| Ruby | `.rb` |
| PHP | `.php` |

Extraction is regex-only — no compilation, no AST, no language server.

## CLI

```bash
reckon              # SessionStart mode: rebuild index, emit hook response
reckon index        # rebuild .codeindex only (no hook output), print stats
reckon stats        # show symbol count, file count, language breakdown
reckon init         # install Claude Code SessionStart hook
reckon uninstall    # remove hook and config
reckon version      # show version
reckon help         # show usage
```

## Config

Drop a `.codeindex.yml` in your project root to override skip patterns:

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

The patterns in `.codeindex.yml` are **appended** to the built-in defaults — you don't need to repeat them. Most projects need zero config.

Add `.codeindex` to your `.gitignore` — it's rebuilt on every session start and should not be committed.

## License

MIT
