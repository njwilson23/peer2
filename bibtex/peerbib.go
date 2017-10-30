package main

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/urfave/cli.v1"
)

func main() {

	app := cli.NewApp()
	app.Name = "peerbib"
	app.Version = "0.3.0dev"
	app.Usage = "peer [--bibtex BIBFILE] [--author AUTHOR] [--year YEAR] [--title TITLE] [search_terms...]"

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
		cli.BoolFlag{
			Name:  "key-only, k",
			Usage: "Only print BibTeX key",
		},
	}

	app.Action = func(c *cli.Context) error {

		var bibtexResults []Entry
		searchAuthor := c.String("author")
		searchTitle := c.String("title")
		searchYear := c.Int("year")

		bibfile := c.String("bibtex")
		entries := make(chan Entry)
		go ReadBibTeX(bibfile, entries)

		for entry := range entries {

			if entry.TestAuthor(searchAuthor) &&
				entry.TestTitle(searchTitle) &&
				entry.TestYear(searchYear) {
				bibtexResults = append(bibtexResults, entry)
			}

		}

		sort.Sort(ByYear(bibtexResults))

		if c.Bool("key-only") {
			for _, entry := range bibtexResults {
				fmt.Println(entry.BibTeXkey)
			}

			return nil
		}

		for _, entry := range bibtexResults {
			fmt.Println(fmt.Sprintf("@%v\n%v (%v)\n\"%v\"\n",
				entry.BibTeXkey,
				entry.Author,
				entry.Year,
				entry.Title))
		}

		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
