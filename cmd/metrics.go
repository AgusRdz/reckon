package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/agusrdz/reckon/metrics"
)

// Metrics prints usage statistics from the local build log.
func Metrics(args []string) {
	if len(args) > 0 && args[0] == "--clear" {
		if err := metrics.Clear(); err != nil {
			fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("usage log cleared")
		return
	}

	entries, err := metrics.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}
	if len(entries) == 0 {
		fmt.Println("no usage data yet — run 'reckon index' or wait for a SessionStart")
		return
	}

	type projectStats struct {
		builds    int
		symbols   int
		lastBuild time.Time
	}

	byProject := map[string]*projectStats{}
	var total int
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var todayCount int

	for _, e := range entries {
		total++
		if !e.Timestamp.Before(today) {
			todayCount++
		}
		ps := byProject[e.Project]
		if ps == nil {
			ps = &projectStats{}
			byProject[e.Project] = ps
		}
		ps.builds++
		ps.symbols += e.Symbols
		if e.Timestamp.After(ps.lastBuild) {
			ps.lastBuild = e.Timestamp
		}
	}

	// Sort projects by most builds
	projects := make([]string, 0, len(byProject))
	for p := range byProject {
		projects = append(projects, p)
	}
	sort.Slice(projects, func(i, j int) bool {
		return byProject[projects[i]].builds > byProject[projects[j]].builds
	})

	fmt.Printf("total builds:  %s  (%s today)  across %s project(s)\n\n",
		formatNum(total),
		formatNum(todayCount),
		formatNum(len(byProject)),
	)

	const (
		colProject = 32
		colBuilds  = 8
		colAvg     = 14
		colLast    = 0
	)
	fmt.Printf("  %-*s  %-*s  %-*s  %s\n",
		colProject, "project",
		colBuilds, "builds",
		colAvg, "avg symbols",
		"last build",
	)
	fmt.Printf("  %s  %s  %s  %s\n",
		strings.Repeat("-", colProject),
		strings.Repeat("-", colBuilds),
		strings.Repeat("-", colAvg),
		strings.Repeat("-", 12),
	)

	for _, p := range projects {
		ps := byProject[p]
		avg := ps.symbols / ps.builds
		fmt.Printf("  %-*s  %-*s  %-*s  %s\n",
			colProject, truncate(shortPath(p), colProject),
			colBuilds, formatNum(ps.builds),
			colAvg, formatNum(avg),
			timeAgo(ps.lastBuild),
		)
	}
}

// shortPath returns ~ for home dir prefix and shortens long paths.
func shortPath(p string) string {
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(p, home) {
		p = "~" + p[len(home):]
	}
	return filepath.ToSlash(p)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return "…" + s[len(s)-(max-1):]
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%d minute%s ago", m, plural(m))
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%d hour%s ago", h, plural(h))
	case d < 48*time.Hour:
		return "yesterday"
	default:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
