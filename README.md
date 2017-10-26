# peer

This is a simple PDF and BibTeX search tool. It started as a bash script,
then a Python program, and now it's written in go.

## Things it does:

- Search for PDFs by name

    `peer search_term1 [search_term2...]`

- Open a matching PDF

    `peer -o N search_terms...`

- Scan BibTeX for references

    `peer --bibtex bibfile.bib --author Jenkins --year 1999`

## Things it might someday do:

- return formatted references

    `peer ref --style agu08.bst --author Jenkins --year 1999`

- add papers to bibtex file

    `peer import fnm`

- manage a reference database, sorting and updating as new papers are downloaded

