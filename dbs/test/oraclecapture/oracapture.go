package main

import (
	"flag"
	"fmt"
	"github.com/892294101/dds/dbs/extract/oracle"
	"os"
)

var processName = flag.String("processid", "", "Please enter the process id name")

func main() {
	flag.Parse()
	if processName == nil || len(*processName) == 0 {
		os.Exit(1)
	}

	ext := oracle.NewCapture()
	if ext != nil {
		err := ext.InitConfig(*processName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "%s", err)
			if err != nil {
				return
			}
			os.Exit(2)
		}
	}

}
