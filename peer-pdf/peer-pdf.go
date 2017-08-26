package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/njwilson23/peer2/config"
	"launchpad.net/gnuflag"
)

// USAGE:
// peer-pdf [query] [-r] [-p] [-o N]

// MatchedFile represents a PDF file found matching the search
type MatchedFile struct {
	Filename string
	Priority int
}

// ForConsole arranges PDF files in order for console printing
type ForConsole []MatchedFile

func (a ForConsole) Len() int      { return len(a) }
func (a ForConsole) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Sort by priority, then by year
func (a ForConsole) Less(i, j int) bool {

	if a[i].Priority != a[j].Priority {
		return a[i].Priority < a[j].Priority
	}
	reYear := regexp.MustCompile(`[0-9]{4}`)
	iyears := reYear.FindString(a[i].Filename)
	jyears := reYear.FindString(a[j].Filename)
	iyear, jyear := 0, 1
	if iyears != "" && jyears != "" {
		iyear, _ = strconv.Atoi(iyears)
		jyear, _ = strconv.Atoi(jyears)
	}
	return iyear < jyear
}

// TestFilename checks whether a filename is matched by any of the search strings
func TestFilename(fnm string, searchstrs []string) bool {

	match := false

	if strings.ToLower(filepath.Ext(fnm)) == ".pdf" {

		match = true

		for _, searchstr := range searchstrs {
			if !strings.Contains(strings.ToLower(fnm), strings.ToLower(searchstr)) {
				match = false
				break
			}
		}
	}

	return match
}

// SearchRecursive walks a root, sending file matches to a channel
func SearchRecursive(root string, searchstrs []string, out chan<- MatchedFile, done chan<- bool) {

	filepath.Walk(root, func(fnm string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println("ERROR:", fnm, err)
			return filepath.SkipDir
		}

		if TestFilename(fnm, searchstrs) {
			out <- MatchedFile{Filename: fnm, Priority: 2}
		}

		return nil
	})

	done <- true
}

// Search checks a directory without recursing, sending matches to a channel
func Search(dirname string, searchstrs []string, out chan<- MatchedFile, done chan<- bool) {

	infos, err := ioutil.ReadDir(dirname)

	if err != nil {
		fmt.Println("ERROR:", dirname, err)
	} else {

		for _, info := range infos {
			if TestFilename(info.Name(), searchstrs) {
				out <- MatchedFile{Filename: info.Name(), Priority: 0}
			}
		}
	}

	done <- true
}

func main() {

	var searchstrs []string
	open := gnuflag.Int("o", 0, "Open option N using the configured reader")
	printpath := gnuflag.Bool("p", false, "Print full paths")
	rawpath := gnuflag.Bool("r", false, "Return raw paths")
	gnuflag.Parse(true)

	configPath, err := config.FindConfig()
	if err != nil {
		panic(err.Error())
	}
	config := config.ParseConfig(configPath)

	if len(gnuflag.Args()) != 0 {
		searchstrs = gnuflag.Args()
	} else {
		fmt.Println("ERROR: must provide at least one search query")
		os.Exit(1)
	}
	if len(config.SearchRoots) == 0 {
		fmt.Println("WARNING: list of search roots is empty")
	}

	results := make(chan MatchedFile)
	done := make(chan bool)
	nrunning := 0

	// Search current directory
	curdir, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR: cannot detect current directory")
		fmt.Println(err)
		os.Exit(1)
	}
	go Search(curdir, searchstrs, results, done)
	nrunning++

	// Search each directory root
	for _, root := range config.SearchRoots {
		go SearchRecursive(root, searchstrs, results, done)
		nrunning++
	}

	ncompleted := 0
	var foundfiles []MatchedFile
	for {
		select {
		case match := <-results:
			foundfiles = append(foundfiles, match)
		case <-done:
			ncompleted++
		case <-time.After(10 * time.Second):
			ncompleted = nrunning
		}
		if ncompleted == nrunning {
			//close(done)
			//close(results)
			break
		}
	}

	sort.Sort(ForConsole(foundfiles))

	if *open != 0 {
		// Open selection
		if (*open > len(results)) || (*open < 0) {
			fmt.Println("Index outside range of results found")
		} else {
			cmd := exec.Command(config.Reader, foundfiles[*open-1].Filename)
			cmd.Start()
		}
	} else {
		// Print matches
		if *rawpath {
			fmt.Print(strings.Trim(fmt.Sprint(results), "[]"))
		} else {

			for i, match := range foundfiles {
				if *printpath {
					fmt.Println(i+1, match)
				} else {
					fmt.Println(i+1, filepath.Base(match.Filename))
				}
			}
		}
	}
}
