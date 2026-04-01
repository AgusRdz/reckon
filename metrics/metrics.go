package metrics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	TypeBuild  = "build"
	TypeSearch = "search"
)

// Result is one symbol match from a search.
type Result struct {
	Symbol string `json:"symbol"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	Kind   string `json:"kind"`
}

// Entry is one event written to the usage log.
// Type is either "build" or "search"; unused fields are omitted.
type Entry struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"ts"`
	Project   string    `json:"project"`
	// build fields
	Symbols   int      `json:"symbols,omitempty"`
	Files     int      `json:"files,omitempty"`
	Languages []string `json:"languages,omitempty"`
	// search fields
	Query   string   `json:"query,omitempty"`
	Hits    int      `json:"hits,omitempty"`
	Results []Result `json:"results,omitempty"`
}

// RecordBuild appends a build event to the usage log.
func RecordBuild(project string, symbols, files int, languages []string) {
	write(Entry{
		Type:      TypeBuild,
		Timestamp: time.Now().UTC(),
		Project:   project,
		Symbols:   symbols,
		Files:     files,
		Languages: languages,
	})
}

// RecordSearch appends a search event to the usage log.
func RecordSearch(project, query string, results []Result) {
	write(Entry{
		Type:      TypeSearch,
		Timestamp: time.Now().UTC(),
		Project:   project,
		Query:     query,
		Hits:      len(results),
		Results:   results,
	})
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
			// back-compat: entries written before type field default to build
			if e.Type == "" {
				e.Type = TypeBuild
			}
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

func write(e Entry) {
	// Failures are silently ignored — metrics must never break the main flow.
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
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	f.Write(append(data, '\n'))
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
