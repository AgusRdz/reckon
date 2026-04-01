package cmd

import (
	"testing"
)

func TestFormatNum(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{999, "999"},
		{1000, "1,000"},
		{1243, "1,243"},
		{9999, "9,999"},
		{10000, "10,000"},
		{100000, "100,000"},
		{1000000, "1,000,000"},
		{1234567, "1,234,567"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatNum(tt.input)
			if got != tt.want {
				t.Errorf("formatNum(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtToLang(t *testing.T) {
	tests := []struct {
		ext  string
		want string
	}{
		{".go", "Go"},
		{".ts", "TypeScript"},
		{".tsx", "TypeScript"},
		{".js", "JavaScript"},
		{".jsx", "JavaScript"},
		{".cs", "C#"},
		{".py", "Python"},
		{".java", "Java"},
		{".rs", "Rust"},
		{".rb", "Ruby"},
		{".php", "PHP"},
		{".unknown", ".unknown"},
		{"", ""},
		{".txt", ".txt"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := extToLang(tt.ext)
			if got != tt.want {
				t.Errorf("extToLang(%q) = %q, want %q", tt.ext, got, tt.want)
			}
		})
	}
}
