package main

import (
	"fmt"
	"github.com/892294101/dds/dbs/sci/terminal"

	"os"
)

func main() {
	err := terminal.OpenShell()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
	}

}
