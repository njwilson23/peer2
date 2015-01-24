package bibtex

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Entry struct {
	Title     string
	Author    string
	Year      int
	Journal   string
	BibTeXkey string
}

func (entry Entry) String() string {
	return fmt.Sprintf("@%v\nTitle: \"%v\"\nAuthor: %v\nJournal: %v",
		entry.BibTeXkey,
		entry.Title,
		entry.Author,
		entry.Journal)
}

type ParseError struct {
	Message string
}

func (err ParseError) Error() string {
	return fmt.Sprintf(err.Message)
}

// Given a line from a BibTeX file, attempt to return a key : value pair
func ParseLine(s string) (string, string, error) {
	var key, value string
	var err error
	if strings.ContainsRune(s, '=') {
		pieces := strings.Split(s, "=")
		key = strings.ToLower(strings.Trim(pieces[0], " \t"))
		value = strings.Trim(pieces[1], " \t{},")
	} else if strings.HasPrefix(s, "@") {
		key = "BibTeXkey"
		value = strings.TrimSuffix(strings.Split(s, "{")[1], ",")
	} else {
		err = ParseError{"not a valid key-value"}
	}
	return key, value, err
}

// Given an array of lines representing a complete BibTeX entry, return an Entry
// type
func ParseEntry(lines []string) (Entry, error) {
	var err error
	var title, author, journal, key string
	var year int
	for _, line := range lines {
		switch k, v, _ := ParseLine(line); k {
		case "author":
			author = v
		case "title":
			title = v
		case "year":
			year, err = strconv.Atoi(v)
			if err != nil {
				break
			}
		case "journal":
			journal = v
		case "BibTeXkey":
			key = v
		}

	}
	entry := Entry{title, author, year, journal, key}
	return entry, err
}

// Open and read a BibTeX database and return an array of BibTeX entries
// This prints any errors raised
func Read(fnm string, entries chan Entry) {
	defer close(entries)
	data, err := ioutil.ReadFile(fnm)
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(data), "\n")

	var count, start int
	inside_entry := false
	for i, line := range lines {
		count += strings.Count(line, "{")
		count -= strings.Count(line, "}")
		if !inside_entry && count != 0 {
			inside_entry = true
		}
		if count == 0 && inside_entry {
			entry, err := ParseEntry(lines[start : i+1])
			inside_entry = false
			if err == nil {
				entries <- entry
			} else {
				fmt.Println(err)
			}
			start = i
		} else if count < 0 {
			fmt.Println("malformed bibtex has more closing than opening braces")
			break
		}
	}
}

// Removes LaTeX-y symbols from *s*. If *lc* is true, then result is converted
// to lower case
func sanitize(s string, lc bool) string {
	out := s
	if strings.Contains(s, "\\\"") {
		out = strings.Replace(s, "\\\"", "", -1)
	}
	if strings.ContainsAny(s, "{}") {
		out = strings.Replace(s, "{", "", -1)
		out = strings.Replace(s, "}", "", -1)
	}
	if lc {
		out = strings.ToLower(out)
	}
	return out
}

// Search a slice of BibTeX entries for author text matching a substring
func SearchAuthor(entries []Entry, s string) []Entry {
	found := make([]Entry, 0)
	for _, entry := range entries {
		auth := sanitize(entry.Author, false)
		if strings.Contains(auth, s) {
			found = append(found, entry)
		}
	}
	return found
}

// Search a slice of BibTeX entries for author text matching a substring
func SearchTitle(entries []Entry, s string) []Entry {
	found := make([]Entry, 0)
	s = strings.ToLower(s)
	for _, entry := range entries {
		title := sanitize(entry.Title, true)
		if strings.Contains(title, s) {
			found = append(found, entry)
		}
	}
	return found
}
