package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
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

// Interface for sorting
type ByYear []Entry

func (a ByYear) Len() int           { return len(a) }
func (a ByYear) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByYear) Less(i, j int) bool { return a[i].Year < a[j].Year }

func (entry Entry) String() string {
	return fmt.Sprintf("@%v\nTitle: \"%v\"\nAuthor: %v\nYear: %v\nJournal: %v\n",
		entry.BibTeXkey,
		entry.Title,
		entry.Author,
		entry.Year,
		entry.Journal)
}

func (e Entry) TestAuthor(auth string) bool {
	return auth == "" || strings.Contains(e.Author, auth)
}

func (e Entry) TestTitle(title string) bool {
	titleLower := strings.ToLower(title)
	return titleLower == "" || strings.Contains(strings.ToLower(e.Title), titleLower)
}

func (e Entry) TestYear(year int) bool {
	return year == -1000000 || e.Year == year
}

type ParseError struct {
	Message string
}

func (err ParseError) Error() string {
	return fmt.Sprintf(err.Message)
}

// Given a string that is a valid BibTeX value, return a unicode representation
func UnicodeBibValue(str string) string {
	var ustr string
	ustr = strings.Replace(str, "\\\"o", "ő", -1)
	ustr = strings.Replace(ustr, "\\\"u", "ű", -1)
	ustr = strings.Replace(ustr, "{\\ae}", "æ", -1)
	ustr = strings.Replace(ustr, "{", "", -1)
	ustr = strings.Replace(ustr, "}", "", -1)
	ustr = strings.Trim(ustr, " \t,\"")
	return ustr
}

// Given a line from a BibTeX file, attempt to return a key : value pair
func parseLine(s string) (string, string, error) {
	var key, value string
	var err error
	if strings.ContainsRune(s, '=') {
		pieces := strings.Split(s, "=")
		key = strings.ToLower(strings.Trim(pieces[0], " \t"))
		value = UnicodeBibValue(pieces[1])
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
func parseEntry(lines []string) (Entry, error) {
	var err error
	var title, author, journal, key string
	var year int
	for _, line := range lines {
		switch k, v, _ := parseLine(line); k {
		case "author":
			author = v
		case "title":
			title = v
		case "year":
			year, err = strconv.Atoi(strings.Trim(v, "\""))
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

func combinerunninglines(lines []string) []string {
	var depth, start int
	var joinedline string
	var linestojoin []string
	joinedlines := make([]string, 0)
	intag := false
	inquote := false

	re_quote := regexp.MustCompile(`[^\\]\"`)

	for i := range lines {
		depth += strings.Count(lines[i], "{")
		depth -= strings.Count(lines[i], "}")

		if len(re_quote.FindAllString(lines[i], -1)) == 1 {
			inquote = !inquote
			if inquote {
				depth += 1
			} else {
				depth -= 1
			}
		}

		if depth > 1 {
			// inside a tag
			if !intag {
				start = i
				intag = true
			}
		} else {
			if intag {
				linestojoin = lines[start : i+1]
				for j, line := range linestojoin {
					line = strings.TrimSpace(line)
					linestojoin[j] = line
				}
				joinedline = strings.Join(linestojoin, " ")
				intag = false
			} else {
				joinedline = lines[i]
			}
			joinedlines = append(joinedlines, joinedline)
		}
	}
	return joinedlines
}

func parseEntries(lines []string, entries chan Entry) error {
	var depth, start int
	var err error
	joinedlines := combinerunninglines(lines)
	inentry := false
	for i, line := range joinedlines {
		depth += strings.Count(line, "{")
		depth -= strings.Count(line, "}")
		if depth == 0 {
			if inentry {
				// end of entry
				inentry = false
				entry, err := parseEntry(joinedlines[start : i+1])
				if err == nil {
					entries <- entry
				} else {
					fmt.Println(err)
				}
			}
		} else if depth == 1 {
			// inside entry but not tag
			if !inentry {
				inentry = true
				start = i
			}
		} else {
			// should never happen
			if depth > 1 {
				err = ParseError{"running lines"}
			} else if depth < 1 {
				err = ParseError{"malformed BibTeX: negative bracket depth"}
			}
			break
		}
		i++
		if i == len(lines) {
			break
		}
	}
	return err
}

// Open and read a BibTeX database and return an array of BibTeX entries
// This prints any errors raised
func ReadBibTeX(fnm string, entries chan Entry) {
	defer close(entries)
	data, err := ioutil.ReadFile(fnm)
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(data), "\n")
	err = parseEntries(lines, entries)
	if err != nil {
		fmt.Println(err)
	}
}

// Removes LaTeX-y symbols from *s*.
func sanitize(s string) string {
	out := s
	if strings.Contains(s, "\\\"") {
		out = strings.Replace(out, "\\\"", "", -1)
	}
	if strings.ContainsAny(s, "{}\\'") {
		out = strings.Replace(out, "{", "", -1)
		out = strings.Replace(out, "}", "", -1)
		out = strings.Replace(out, "\\", "", -1)
		out = strings.Replace(out, "'", "", -1)
	}
	return out
}

// Search a slice of BibTeX entries for author text matching a substring
func SearchAuthor(entries []Entry, s string) []Entry {
	found := make([]Entry, 0)
	for _, entry := range entries {
		auth := sanitize(entry.Author)
		if strings.Contains(auth, s) {
			found = append(found, entry)
		}
	}
	return found
}

// Search a slice of BibTeX entries for title text matching a substring
func SearchTitle(entries []Entry, s string) []Entry {
	found := make([]Entry, 0)
	s = strings.ToLower(s)
	for _, entry := range entries {
		title := strings.ToLower(sanitize(entry.Title))
		if strings.Contains(title, s) {
			found = append(found, entry)
		}
	}
	return found
}

// Search a slice of BibTeX entries for year
func SearchYear(entries []Entry, ymin int, ymax int) []Entry {
	found := make([]Entry, 0)
	for _, entry := range entries {
		if (ymin <= entry.Year) && (entry.Year <= ymax) {
			found = append(found, entry)
		}
	}
	return found
}
