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

	var builds, searches []metrics.Entry
	for _, e := range entries {
		if e.Type == metrics.TypeSearch {
			searches = append(searches, e)
		} else {
			builds = append(builds, e)
		}
	}

	printBuildStats(builds)
	if len(searches) > 0 {
		fmt.Println()
		printSearchStats(searches)
	}
}

func printBuildStats(entries []metrics.Entry) {
	if len(entries) == 0 {
		return
	}

	type projectStats struct {
		builds    int
		symbols   int
		lastBuild time.Time
	}

	byProject := map[string]*projectStats{}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var todayCount int

	for _, e := range entries {
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

	projects := make([]string, 0, len(byProject))
	for p := range byProject {
		projects = append(projects, p)
	}
	sort.Slice(projects, func(i, j int) bool {
		return byProject[projects[i]].builds > byProject[projects[j]].builds
	})

	fmt.Printf("builds:  %s total  (%s today)  across %s project(s)\n\n",
		formatNum(len(entries)),
		formatNum(todayCount),
		formatNum(len(byProject)),
	)
	const (colP, colB, colA = 32, 8, 14)
	fmt.Printf("  %-*s  %-*s  %-*s  %s\n", colP, "project", colB, "builds", colA, "avg symbols", "last build")
	fmt.Printf("  %s  %s  %s  %s\n", strings.Repeat("-", colP), strings.Repeat("-", colB), strings.Repeat("-", colA), strings.Repeat("-", 12))
	for _, p := range projects {
		ps := byProject[p]
		avg := ps.symbols / ps.builds
		fmt.Printf("  %-*s  %-*s  %-*s  %s\n",
			colP, truncate(shortPath(p), colP),
			colB, formatNum(ps.builds),
			colA, formatNum(avg),
			timeAgo(ps.lastBuild),
		)
	}
}

func printSearchStats(entries []metrics.Entry) {
	type queryStats struct {
		count int
		hits  int
	}
	type projectStats struct {
		searches int
		hits     int
		misses   int
	}

	byQuery := map[string]*queryStats{}
	byProject := map[string]*projectStats{}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var todayCount, totalHits, totalMisses int

	for _, e := range entries {
		hit := e.Hits > 0
		if !e.Timestamp.Before(today) {
			todayCount++
		}
		if hit {
			totalHits++
		} else {
			totalMisses++
		}

		qs := byQuery[e.Query]
		if qs == nil {
			qs = &queryStats{}
			byQuery[e.Query] = qs
		}
		qs.count++
		if hit {
			qs.hits++
		}

		ps := byProject[e.Project]
		if ps == nil {
			ps = &projectStats{}
			byProject[e.Project] = ps
		}
		ps.searches++
		if hit {
			ps.hits++
		} else {
			ps.misses++
		}
	}

	total := len(entries)
	hitRate := 0
	if total > 0 {
		hitRate = totalHits * 100 / total
	}

	fmt.Printf("searches:  %s total  (%s today)  hit rate: %d%%  (%s found / %s not found)\n\n",
		formatNum(total),
		formatNum(todayCount),
		hitRate,
		formatNum(totalHits),
		formatNum(totalMisses),
	)

	// Top queries (by frequency, show up to 10)
	type queryEntry struct {
		query string
		qs    *queryStats
	}
	qlist := make([]queryEntry, 0, len(byQuery))
	for q, qs := range byQuery {
		qlist = append(qlist, queryEntry{q, qs})
	}
	sort.Slice(qlist, func(i, j int) bool {
		return qlist[i].qs.count > qlist[j].qs.count
	})
	if len(qlist) > 10 {
		qlist = qlist[:10]
	}

	const (colQ, colC, colH = 30, 8, 8)
	fmt.Printf("  top queries:\n")
	fmt.Printf("  %-*s  %-*s  %s\n", colQ, "query", colC, "count", "found")
	fmt.Printf("  %s  %s  %s\n", strings.Repeat("-", colQ), strings.Repeat("-", colC), strings.Repeat("-", 5))
	for _, qe := range qlist {
		found := "no"
		if qe.qs.hits > 0 {
			found = "yes"
		}
		fmt.Printf("  %-*s  %-*s  %s\n", colQ, truncate(qe.query, colQ), colC, formatNum(qe.qs.count), found)
	}

	// Per-project search stats
	if len(byProject) > 1 {
		fmt.Printf("\n  per project:\n")
		projects := make([]string, 0, len(byProject))
		for p := range byProject {
			projects = append(projects, p)
		}
		sort.Slice(projects, func(i, j int) bool {
			return byProject[projects[i]].searches > byProject[projects[j]].searches
		})
		const colP2 = 32
		fmt.Printf("  %-*s  %-*s  %-*s  %s\n", colP2, "project", colC, "searches", colH, "hits", "misses")
		fmt.Printf("  %s  %s  %s  %s\n", strings.Repeat("-", colP2), strings.Repeat("-", colC), strings.Repeat("-", colH), strings.Repeat("-", 6))
		for _, p := range projects {
			ps := byProject[p]
			fmt.Printf("  %-*s  %-*s  %-*s  %s\n",
				colP2, truncate(shortPath(p), colP2),
				colC, formatNum(ps.searches),
				colH, formatNum(ps.hits),
				formatNum(ps.misses),
			)
		}
	}
}

// shortPath returns ~ for home dir prefix.
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
