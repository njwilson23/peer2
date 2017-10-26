package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/njwilson23/peer2/bibtex"
	"gopkg.in/urfave/cli.v1"
)

func main() {

	app := cli.NewApp()
	app.Name = "peer"
	app.Version = "0.3.0dev"
	app.Usage = "peer [--bibtex bibfile] [--path filepath] [--open N] [--reference N] search_terms..."

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Value: wd,
			Usage: "Search root",
		},
		cli.IntFlag{
			Name:  "open, o",
			Value: -1,
			Usage: "Open one of the search results",
		},
		cli.StringFlag{
			Name:  "bibtex, b",
			Value: "",
			Usage: "Search a configured BibTeX file",
		},
		cli.StringFlag{
			Name:  "author",
			Value: "",
			Usage: "Author filter for BibTeX searches",
		},
		cli.StringFlag{
			Name:  "title",
			Value: "",
			Usage: "Title filter for BibTeX searches",
		},
		cli.IntFlag{
			Name:  "year",
			Value: -1000000,
			Usage: "Published year filter for BibTeX searches",
		},
	}

	app.Action = func(c *cli.Context) error {
		// search bibtex
		if c.String("bibtex") != "" {

			var bibtexResults []bibtex.Entry
			searchAuthor := c.String("author")
			searchTitle := c.String("title")
			searchYear := c.Int("year")

			bibfile := c.String("bibtex")
			entries := make(chan bibtex.Entry)
			go bibtex.Read(bibfile, entries)

			for entry := range entries {

				if entry.TestAuthor(searchAuthor) &&
					entry.TestTitle(searchTitle) &&
					entry.TestYear(searchYear) {
					bibtexResults = append(bibtexResults, entry)
				}

			}

			sort.Sort(bibtex.ByYear(bibtexResults))
			for _, entry := range bibtexResults {
				fmt.Println(fmt.Sprintf("@%v\n%v (%v)\n\t\"%v\"\n",
					entry.BibTeXkey,
					entry.Author,
					entry.Year,
					entry.Title))
			}

			return nil
		}

		// search PDFs
		searchTerms := c.Args()
		if len(searchTerms) == 0 {
			return errors.New("at least one search term must be provided")
		}
		roots := []string{c.String("path")}
		results := search(roots, searchTerms)

		if c.Int("open") != -1 {
			idx := c.Int("open")
			if idx > len(results) {
				return errors.New("invalid index to open")
			}
			cmd := exec.Command("evince", results[idx-1].path)
			cmd.Start()
		}

		for _, result := range results {
			fmt.Printf("%70s\t%.2f\n", result, result.score)
		}
		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
