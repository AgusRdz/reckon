package hooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func claudeSettingsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "settings.json")
}

func loadSettings() map[string]interface{} {
	data, err := os.ReadFile(claudeSettingsPath())
	if err != nil {
		return map[string]interface{}{}
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return map[string]interface{}{}
	}
	return m
}

func saveSettings(m map[string]interface{}) error {
	path := claudeSettingsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func exePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

// Install registers reckon as a Claude Code SessionStart hook.
func Install(version string) {
	exe, err := exePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: failed to get executable path: %v\n", err)
		os.Exit(1)
	}

	settings := loadSettings()
	hooksMap := getOrCreateMap(settings, "hooks")
	sessionStart := getOrCreateSlice(hooksMap, "SessionStart")
	sessionStart = removeOurEntries(sessionStart)

	// Use forward slashes and quoted path to match Claude Code's expected format.
	// Prepend so reckon runs before ctx restore — index must be injected first.
	exeFwd := strings.ReplaceAll(exe, "\\", "/")
	entry := map[string]interface{}{
		"hooks": []interface{}{
			map[string]interface{}{
				"type":    "command",
				"command": fmt.Sprintf("%q hook", exeFwd),
			},
		},
	}
	sessionStart = append([]interface{}{entry}, sessionStart...)

	hooksMap["SessionStart"] = sessionStart
	settings["hooks"] = hooksMap

	if err := saveSettings(settings); err != nil {
		fmt.Fprintf(os.Stderr, "reckon: failed to write settings: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("reckon %s hook installed\n", version)
	fmt.Printf("  binary: %s\n", exe)
	fmt.Printf("  config: %s\n", claudeSettingsPath())
}

// Uninstall removes the reckon SessionStart hook.
func Uninstall() {
	settings := loadSettings()

	hooksMap, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		fmt.Println("no hook found")
		return
	}
	sessionStart, ok := hooksMap["SessionStart"].([]interface{})
	if !ok {
		fmt.Println("no hook found")
		return
	}

	before := len(sessionStart)
	sessionStart = removeOurEntries(sessionStart)
	hooksMap["SessionStart"] = sessionStart
	settings["hooks"] = hooksMap

	if len(sessionStart) == before {
		fmt.Println("no hook found")
		return
	}

	if err := saveSettings(settings); err != nil {
		fmt.Fprintf(os.Stderr, "reckon: failed to write settings: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("hook removed from ~/.claude/settings.json")
}

// IsInstalled checks if the hook is installed.
func IsInstalled() (bool, string) {
	cmd := GetHookCommand()
	return cmd != "", cmd
}

// GetHookCommand returns the currently registered hook command, or "".
func GetHookCommand() string {
	settings := loadSettings()
	hooksMap, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		return ""
	}
	sessionStart, ok := hooksMap["SessionStart"].([]interface{})
	if !ok {
		return ""
	}
	for _, entry := range sessionStart {
		m, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		hooksList, ok := m["hooks"].([]interface{})
		if !ok {
			continue
		}
		for _, h := range hooksList {
			hm, ok := h.(map[string]interface{})
			if !ok {
				continue
			}
			cmd, _ := hm["command"].(string)
			if strings.Contains(cmd, "reckon") && strings.HasSuffix(cmd, " hook") {
				return cmd
			}
		}
	}
	return ""
}

func getOrCreateMap(m map[string]interface{}, key string) map[string]interface{} {
	v, ok := m[key].(map[string]interface{})
	if !ok {
		v = map[string]interface{}{}
	}
	return v
}

func getOrCreateSlice(m map[string]interface{}, key string) []interface{} {
	v, ok := m[key].([]interface{})
	if !ok {
		v = []interface{}{}
	}
	return v
}

func removeOurEntries(entries []interface{}) []interface{} {
	var result []interface{}
	for _, entry := range entries {
		m, ok := entry.(map[string]interface{})
		if !ok {
			result = append(result, entry)
			continue
		}
		hooksList, ok := m["hooks"].([]interface{})
		if !ok {
			result = append(result, entry)
			continue
		}
		isOurs := false
		for _, h := range hooksList {
			hm, ok := h.(map[string]interface{})
			if !ok {
				continue
			}
			cmd, _ := hm["command"].(string)
			if strings.Contains(cmd, "reckon") {
				// match both old format ("reckon") and new ("reckon hook")
				isOurs = true
				break
			}
		}
		if !isOurs {
			result = append(result, entry)
		}
	}
	return result
}
