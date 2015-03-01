package bibtex

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

func (entry Entry) String() string {
	return fmt.Sprintf("@%v\nTitle: \"%v\"\nAuthor: %v\nYear: %v\nJournal: %v\n",
		entry.BibTeXkey,
		entry.Title,
		entry.Author,
		entry.Year,
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
		value = strings.Trim(pieces[1], " \t{},\"")
	} else if strings.HasPrefix(s, "@") {
		key = "BibTeXkey"
		value = strings.TrimSuffix(strings.Split(s, "{")[1], ",")
	} else {
		err = ParseError{"not a valid key-value"}
	}
	return key, value, err
}

// BibTeX can use either curly braces or quotation marks to enclose content.
// This indentifies which convention is being used.
/*func enclosingmark(line string) (rune, error) {
	var sym rune
	var err error
	idx := strings.Index(line, "=")
	if idx == -1 {
		err = ParseError{"no equals sign to split tag from content"}
		sym = ' '
	} else {
		jdx := strings.IndexAny(line[idx:], "{\"")
		r := []rune(line)
		sym = r[idx+jdx]
	}
	return sym, err
}*/

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

func ParseEntries(lines []string, entries chan Entry) error {
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
				entry, err := ParseEntry(joinedlines[start : i+1])
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
func Read(fnm string, entries chan Entry) {
	defer close(entries)
	data, err := ioutil.ReadFile(fnm)
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(data), "\n")
	err = ParseEntries(lines, entries)
	if err != nil {
		fmt.Println(err)
	}
}

// Removes LaTeX-y symbols from *s*. If *lc* is true, then result is converted
// to lower case
func sanitize(s string, lc bool) string {
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

// Search a slice of BibTeX entries for title text matching a substring
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
