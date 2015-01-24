package bibtex

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Entry struct {
	Title     string
	Author    string
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
	for _, line := range lines {
		switch k, v, _ := ParseLine(line); k {
		case "author":
			author = v
		case "title":
			title = v
		case "journal":
			journal = v
		case "BibTeXkey":
			key = v
		}

	}
	entry := Entry{title, author, journal, key}
	return entry, err
}

// Open and read a BibTeX database and return an array of BibTeX entries
func Read(fnm string, entriesPtr *[]Entry) error {

	data, err := ioutil.ReadFile(fnm)
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(data), "\n")
	entries := *entriesPtr

	// Split into entries, and parse them one by one
	// I suppose this could be done asynchronously...
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
				entries = append(entries, entry)
			} else {
				fmt.Println(err)
			}
			start = i
		} else if count < 0 {
			err = ParseError{"malformed bibtex"}
			break
		}

	}
	*entriesPtr = entries
	return err
}
