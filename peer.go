package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var cmd *exec.Cmd
	if os.Args[1] == "bib" {
		cmd = exec.Command("peer-bib", os.Args[1:]...)
	} else {
		cmd = exec.Command("peer-pdf", os.Args[1:]...)
	}
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
