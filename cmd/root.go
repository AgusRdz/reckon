package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Root handles the SessionStart hook invocation.
func Root(version string) {
	// SessionStart provides no meaningful stdin; consume and discard.
	io.ReadAll(os.Stdin) //nolint:errcheck

	cwd, err := os.Getwd()
	if err != nil {
		emitEmpty()
		return
	}

	symbols, stats, err := BuildIndex(cwd)
	if err != nil || len(symbols) == 0 {
		emitEmpty()
		return
	}

	msg := fmt.Sprintf(
		"Symbol index rebuilt: .codeindex — %s symbols across %s files.\nGrep .codeindex before searching the codebase. If not found there, search normally.",
		formatNum(stats.Symbols),
		formatNum(stats.Files),
	)
	respond(msg)
}

func emitEmpty() {
	respond("")
}

func respond(output string) {
	resp := map[string]string{"action": "continue", "output": output}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}
