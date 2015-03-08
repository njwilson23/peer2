# peer-next

This is a rewrite of [peer](https://github.com/njwilson23/peer) in Go. It also
includes a simple BibTeX parser.

## Things it does:

- Search for PDFs

    peer pdf search_times...

- Open a matching PDFs

    peer pdf -o N search_terms...

- Scan BibTeX for references

    peer bib -author Jenkins -year 1999

## Things it might someday do:

- return formatted references
- add papers to bibtex file
- manage a reference database, sorting and updating as new papers are downloaded

