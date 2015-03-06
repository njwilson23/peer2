package main

import (
	"flag"
	"fmt"
	"github.com/njwilson23/peer2/config"
	"gopkg.in/yaml.v2"
)

// USAGE:
// peer-pdf [query] [-r] [-p] [-o N]

func main() {

	config := config.ParseConfig(".peer.yaml")

}
