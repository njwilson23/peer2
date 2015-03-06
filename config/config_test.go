package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	config := ParseConfig("test.yaml")
	if config.Reader != "evince" {
		t.Fail()
	}
	if config.Bibfiles[0] != "biblio.bib" {
		t.Fail()
	}
	if config.Bibfiles[1] != "biblio2.bib" {
		t.Fail()
	}

	if config.SearchRoots[0] != "~/Downloads" {
		t.Fail()
	}
	if config.SearchRoots[1] != "~/Documents/pdfs" {
		t.Fail()
	}
}
