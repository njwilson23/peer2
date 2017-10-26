package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SearchResult struct {
	path  string
	score float64
}

func (r SearchResult) String() string {
	return fmt.Sprintf("%s", r.path)
}

func search(roots []string, searchTerms []string) []SearchResult {

	var results []SearchResult

	for _, root := range roots {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if strings.ToLower(filepath.Ext(path)) != ".pdf" {
				return nil
			}

			score := 0.0
			for _, term := range searchTerms {
				if strings.Contains(path, term) {
					score++
				}
			}
			if score != 0 {
				results = append(results, SearchResult{
					path:  path,
					score: score,
				})
			}
			return nil
		})
	}

	return results
}
