package main

import (
	"fmt"
	"github.com/892294101/dds/ddslog"
	"github.com/892294101/dds/spfile"
	"os"
)

func main() {

	log, err := ddslog.InitDDSlog("HTTK_0002")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pfile, err := spfile.LoadSpfile(fmt.Sprintf("%s.desc", "HTTK_0002"), spfile.UTF8, log, spfile.GetOracleName(), spfile.GetExtractName())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := pfile.Production(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pfile.LoadToDatabase()
	pfile.PutParamsText()
}
