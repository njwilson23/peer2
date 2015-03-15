package bibtex

import (
	"fmt"
	"io/ioutil"
	"sort"
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
	go Read("test.bib", entries)
	i := 0
	for {
		_, ok := <-entries
		if !ok {
			break
		} else {
			i += 1
		}
	}
	if i != 3 {
		fmt.Println(i, "entries read (should be 3)")
		t.Fail()
	}
}

func TestReadEntriesMacsyma(t *testing.T) {
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
	if i != 446 {
		fmt.Println(i, "entries read (should be 446)")
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

	// Test 1
	entries := readtoarray("macsyma.bib")
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
	entries = readtoarray("test.bib")
	results = SearchTitle(entries, "jakobshavn")
	if len(results) == 0 {
		fmt.Println("No entries found matching 'jakobshavn'")
		t.Fail()
	}
}

func TestSearchYear(t *testing.T) {
	entries := readtoarray("test.bib")
	results := SearchYear(entries, 2005, 2013)
	if len(results) != 2 {
		fmt.Sprintln("%v entries found for 2005-2013 (should be 2)", len(results))
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

func TestUnicodeBibValue(t *testing.T) {
	var testStr, outStr string

	testStr = "	Quantifying melt rates under {G}reenland's {N}ioghalvfjerdsbr{\\ae}"
	outStr = UnicodeBibValue(testStr)
	if outStr != "Quantifying melt rates under Greenland's Nioghalvfjerdsbræ" {
		fmt.Println(outStr)
		t.Fail()
	}

	testStr = "\t{M\\\"unchow}"
	outStr = UnicodeBibValue(testStr)
	if outStr != "Műnchow" {
		fmt.Println(outStr)
		t.Fail()
	}
}

func TestCombineRunningLinesBraces(t *testing.T) {
	entry_str := `@Article{Roache:1985:NAG,
  author =       {Patrick J. Roache and Stanly Steinberg},
  title =        {New Approach to Grid Generation Using a Variational
                 Formulation},
  journal =      {AIAA Paper},
  pages =        {360--370},
  year =         {1985},
  CODEN =        {AAPRAQ},
  ISSN =         {0146-3705},
  bibdate =      {Wed Jan 15 15:35:13 MST 1997},
  bibsource =    {Compendex database;
                 http://www.math.utah.edu/pub/tex/bib/macsyma.bib},
  acknowledgement = ack-nhfb,
  affiliationaddress = {Ecodynamics Research Associates, Albuquerque,
                 NM, USA},
  classification = {631; 921; 931},
  conference =   {Collection of Technical Papers --- AIAA 7th
                 Computational Fluid Dynamics Conference.},
  journalabr =   {AIAA Paper},
  keywords =     {behavioral errors; computational fluid dynamics; fluid
                 dynamics; mathematical techniques; symbolic
                 manipulation; Thompson-Thames-Mastin method (TTM
                 method); VAX 780; Vaxima},
  meetingaddress = {Cincinnati, OH, Engl},
  sponsor =      {AIAA, New York, NY, USA},
}`
	entry_strs := strings.Split(entry_str, "\n")
	joinedlines := combinerunninglines(entry_strs)
	if len(joinedlines) != 19 {
		fmt.Println("wrong number of line in joined entry:", len(joinedlines))
		t.Fail()
	}
}

func TestCombineRunningLinesQuotes(t *testing.T) {
	entry_str := `@Article{Roache:1985:NAG,
  author =       "Patrick J. Roache and Stanly Steinberg",
  title =        "New Approach to Grid Generation Using a Variational
                 Formulation",
  journal =      "AIAA Paper",
  pages =        "360--370",
  year =         "1985",
  CODEN =        "AAPRAQ",
  ISSN =         "0146-3705",
  bibdate =      "Wed Jan 15 15:35:13 MST 1997",
  bibsource =    "Compendex database;
                 http://www.math.utah.edu/pub/tex/bib/macsyma.bib",
  acknowledgement = ack-nhfb,
  affiliationaddress = "Ecodynamics Research Associates, Albuquerque,
                 NM, USA",
  classification = "631; 921; 931",
  conference =   "Collection of Technical Papers --- AIAA 7th
                 Computational Fluid Dynamics Conference.",
  journalabr =   "AIAA Paper",
  keywords =     "behavioral errors; computational fluid dynamics; fluid
                 dynamics; mathematical techniques; symbolic
                 manipulation; Thompson-Thames-Mastin method (TTM
                 method); VAX 780; Vaxima",
  meetingaddress = "Cincinnati, OH, Engl",
  sponsor =      "AIAA, New York, NY, USA",
}`
	entry_strs := strings.Split(entry_str, "\n")
	joinedlines := combinerunninglines(entry_strs)
	if len(joinedlines) != 19 {
		fmt.Println("wrong number of line in joined entry:", len(joinedlines))
		t.Fail()
	}
}

func TestCombineRunningLinesInFile(t *testing.T) {

	data, err := ioutil.ReadFile("macsyma.bib")
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(data), "\n")
	joinedlines := combinerunninglines(lines)
	if len(joinedlines) != 6384 {
		fmt.Println(fmt.Sprintf("unexpected number of lines (%v)", len(joinedlines)))
	}
}

func TestSortEntries(t *testing.T) {
	entries := []Entry{Entry{"FirstTitle", "A. Hodges", 1973, "Tests and Units", "@Hodges1973First"},
		Entry{"SecondTitle", "Dana Sukoi", 1985, "Reproducibility Mechanics", "@Sukoi1985Second"},
		Entry{"ThirdTitle", "Carl McIntyre", 1968, "Journal of Validation", "@McIntyre1968Third"}}
	sort.Sort(ByYear(entries))
	if entries[0].Title != "ThirdTitle" {
		t.Fail()
	}
	if entries[1].Title != "FirstTitle" {
		t.Fail()
	}
	if entries[2].Title != "SecondTitle" {
		t.Fail()
	}
}
