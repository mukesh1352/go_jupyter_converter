package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type PathRewriter struct {
	patterns    []patternRule
	datasetRoot string
}

type patternRule struct {
	pattern     *regexp.Regexp
	replaceFunc func(datasetRoot string, match []string) string
}

// Entry point
func NewPathRewriter(datasetRootOptional ...string) *PathRewriter {
	var datasetRoot string
	if len(datasetRootOptional) > 0 && datasetRootOptional[0] != "" {
		datasetRoot = datasetRootOptional[0]
	} else {
		cfg, err := LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		datasetRoot = cfg.DatasetRoot
	}

	return &PathRewriter{
		datasetRoot: datasetRoot,
		patterns: []patternRule{
			newRule(`(?i)(pd\.read_csv\(\s*["'])([^"']+)(["']\s*\))`),
			newRule(`(?i)(open\(\s*["'])([^"']+)(["'])`),
			newRule(`(?i)(np\.load\(\s*["'])([^"']+)(["'])`),
			newRule(`(?i)(json\.load\(\s*open\(\s*["'])([^"']+)(["'])`),
		},
	}
}

// Generic rule generator with fallback
func newRule(regex string) patternRule {
	re := regexp.MustCompile(regex)

	return patternRule{
		pattern: re,
		replaceFunc: func(datasetRoot string, match []string) string {
			originalPath := match[2]

			var fullPath string
			if filepath.IsAbs(originalPath) {
				fullPath = originalPath
			} else {
				relPath := filepath.Clean(originalPath)
				fullPath = filepath.Join(datasetRoot, relPath)
			}

			// Escape for Python (forward slashes recommended)
			escapeForPython := func(path string) string {
				return filepath.ToSlash(path) // avoids unicode escape errors
			}

			// If file exists
			if _, err := os.Stat(fullPath); err == nil {
				log.Printf("📦 Resolving path: original=%s, full=%s", originalPath, fullPath)
				return fmt.Sprintf("%s%s%s", match[1], escapeForPython(fullPath), match[3])
			}

			// Fallback search by extension
			ext := filepath.Ext(originalPath)
			var fallback string

			_ = filepath.WalkDir(datasetRoot, func(path string, d os.DirEntry, err error) error {
				if filepath.Ext(path) == ext && !d.IsDir() {
					fallback = path
					return filepath.SkipDir
				}
				return nil
			})

			if fallback != "" {
				log.Printf("⚠️  File not found: %s, using fallback: %s", fullPath, fallback)
				return fmt.Sprintf("%s%s%s", match[1], escapeForPython(fallback), match[3])
			}

			log.Printf("❌ No fallback found for missing file: %s", fullPath)
			return match[0] // return original
		},
	}
}

// Rewrite code by applying all rules
func (rw *PathRewriter) Rewrite(code string) string {
	for _, rule := range rw.patterns {
		code = rule.pattern.ReplaceAllStringFunc(code, func(s string) string {
			match := rule.pattern.FindStringSubmatch(s)
			if len(match) == 4 {
				return rule.replaceFunc(rw.datasetRoot, match)
			}
			return s
		})
	}
	return code
}
