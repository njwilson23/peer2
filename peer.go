package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

func main() {

	app := cli.NewApp()
	app.Name = "peer"
	app.Version = "0.3.0dev"
	app.Usage = "peer [--path FILEPATH] [--open N] [--reference N] SEARCH_TERMS..."

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
		cli.BoolFlag{
			Name:  "print0",
			Usage: "Print results seperated by a space",
		},
	}

	app.Action = func(c *cli.Context) error {

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

		if c.Bool("print0") {
			names := make([]string, len(results))
			for i, r := range results {
				names[i] = r.path
			}
			fmt.Println(strings.Join(names, " "))
			return nil
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
