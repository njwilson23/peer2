package bibtex

import (
	"fmt"
	"io/ioutil"
	"strings"
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
	go Read("macsyma.bib", entries)
	i := 0
	for {
		_, ok := <-entries
		if !ok {
			break
		} else {
			i += 1
		}
	}
	if i != 527 {
		fmt.Println(i, "entries read (should be 527)")
		t.Fail()
	}
}

func TestSearchAuthor(t *testing.T) {
	entries := readtoarray("macsyma.bib")

	// Test 1
	results := SearchAuthor(entries, "Padget")
	if len(results) != 1 {
		fmt.Println(len(results), "entries found matching 'Padget' (should be 1)")
		t.Fail()
	}

	// Test 2 - make sure accents are accounted for
	results = SearchAuthor(entries, "Maartensson")
	if len(results) != 1 {
		fmt.Println(len(results), "entries found matching 'Maartensson' (should be 1)")
		t.Fail()
	}
}

func TestSearchTitle(t *testing.T) {
	entries := readtoarray("macsyma.bib")

	// Test 1
	results := SearchTitle(entries, "variational formulation")
	if len(results) == 0 {
		fmt.Println("No entries found matching 'variational formulation'")
		t.Fail()
	} else {
		for entry := range results {
			fmt.Println(entry)
		}
	}

	// Test 2 - title searches should be case-insensitive
	/*results = SearchTitle(entries, "jakobshavn")
	if len(results) == 0 {
		fmt.Println("No entries found matching 'jakobshaven'")
		t.Fail()
	}*/
}

func TestSearchYear(t *testing.T) {
	entries := readtoarray("test.bib")
	results := SearchYear(entries, 2005, 2013)
	if len(results) != 2 {
		fmt.Println(len(results), "entries found for 2005-2013 (should be 2)")
		t.Fail()
	}
}

func TestSanitize(t *testing.T) {
	difficult_name := "Bengt M{\\aa}rtensson"
	if sanitize(difficult_name, false) != "Bengt Maartensson" {
		fmt.Println("expected 'Bengt Maartensson' but got",
			sanitize(difficult_name, false))
		t.Fail()
	}

	difficult_word := "M\\'elange"
	if sanitize(difficult_word, true) != "melange" {
		fmt.Println("expected 'melange' but got",
			sanitize(difficult_word, true))
		t.Fail()
	}
}

func TestCombineRunningLines(t *testing.T) {

	data, err := ioutil.ReadFile("macsyma.bib")
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(data), "\n")
	joinedlines := combinerunninglines(lines)
	if len(joinedlines) != 11655 {
		fmt.Println(fmt.Sprintf("unexpected number of lines (%v)", len(joinedlines)))
	}
}
