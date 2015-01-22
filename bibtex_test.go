package bibtex

import (
	//"fmt"
	"testing"
)

func TestReadEntry(t *testing.T) {

	entries := make([]Entry, 0)
	//err := Read("test.bib", entries)
	Read("test.bib", entries)
	/*for _, entry := range entries {
		fmt.Println(entry)
	}*/
}
