package extract

// Symbol is a single symbol extracted from a source file.
type Symbol struct {
	Name string
	File string // CWD-relative path
	Line int
	Kind string // class, method, function, interface, type, enum, struct, const
}

// Extractor extracts symbols from source files.
type Extractor interface {
	Extensions() []string
	Extract(file string, content []byte) []Symbol
}

var extractors []Extractor

// All returns all registered extractors.
func All() []Extractor {
	return extractors
}

// Register registers an extractor.
func Register(e Extractor) {
	extractors = append(extractors, e)
}
