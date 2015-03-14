# peer (2)

This is a for-fun rewrite of [peer](https://github.com/njwilson23/peer) from
Python to Go. It also includes a simple BibTeX parser.

## Things it does:

- Search for PDFs

    `peer pdf search_terms...`

- Open a matching PDF

    `peer pdf search_terms... -o N`

- Scan BibTeX for references

    `peer bib --author Jenkins --year 1999`

## Things it might someday do:

- return formatted references

    `peer ref --style agu08.bst --author Jenkins --year 1999`

- add papers to bibtex file

	`peer import fnm`

- manage a reference database, sorting and updating as new papers are downloaded

