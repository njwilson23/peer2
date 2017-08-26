package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var cmd *exec.Cmd
	if len(os.Args) == 1 {
		fmt.Println("Usage: peer pdf|bib args...")
		os.Exit(0)
	} else if os.Args[1] == "bib" {
		cmd = exec.Command("peer-bib", os.Args[2:]...)
	} else if os.Args[1] == "pdf" {
		cmd = exec.Command("peer-pdf", os.Args[2:]...)
	} else {
		cmd = exec.Command("peer-pdf", os.Args[1:]...)
	}
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
