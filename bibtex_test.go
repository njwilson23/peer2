package bibtex

import (
	"fmt"
	"testing"
)

func readtoarray(fnm string) []Entry {
	entries := make(chan Entry)
	go Read(fnm, entries)
	entriesArray := make([]Entry, 0)
	for entry := range entries {
		entriesArray = append(entriesArray, entry)
	}
	return entriesArray
}

func TestReadEntries(t *testing.T) {
	entries := make(chan Entry)
	//go Read("test.bib", entries)
	go Read("glacier.bib", entries)
	i := 0
	for {
		_, ok := <-entries
		if !ok {
			break
		} else {
			i += 1
		}
	}
	if i == 0 {
		fmt.Println("No entries read successfully")
		t.Fail()
	}
}

func TestSearchAuthor(t *testing.T) {
	entries := readtoarray("test.bib")
	var results []Entry

	// Test 1
	results = SearchAuthor(entries, "Nye")
	if len(results) == 0 {
		fmt.Println("No entries found matching 'Nye'")
		t.Fail()
	}

	// Test 2 - make sure accents are accounted for
	results = SearchAuthor(entries, "Luthi")
	if len(results) == 0 {
		fmt.Println("No entries found matching 'Luthi'")
		t.Fail()
	}
}

func TestSearchTitle(t *testing.T) {
	entries := readtoarray("test.bib")
	var results []Entry

	// Test 1
	results = SearchTitle(entries, "thermal")
	if len(results) == 0 {
		fmt.Println("No entries found matching 'thermal'")
		t.Fail()
	}

	// Test 2 - title searches should be case-insensitive
	results = SearchTitle(entries, "jakobshavn")
	if len(results) == 0 {
		fmt.Println("No entries found matching 'jakobshaven'")
		t.Fail()
	}
}
