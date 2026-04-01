package metrics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry is one index-build event written to the usage log.
type Entry struct {
	Timestamp time.Time `json:"ts"`
	Project   string    `json:"project"` // absolute path
	Symbols   int       `json:"symbols"`
	Files     int       `json:"files"`
	Languages []string  `json:"languages"`
}

// Record appends one build event to the usage log.
// Failures are silently ignored — metrics must never break the main flow.
func Record(project string, symbols, files int, languages []string) {
	path, err := logPath()
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Project:   project,
		Symbols:   symbols,
		Files:     files,
		Languages: languages,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	f.Write(append(data, '\n'))
}

// ReadAll returns all entries from the usage log, oldest first.
func ReadAll() ([]Entry, error) {
	path, err := logPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if json.Unmarshal(line, &e) == nil {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

// Clear deletes the usage log.
func Clear() error {
	path, err := logPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func logPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "reckon", "usage.log"), nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
