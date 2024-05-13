package main

import (
	"fmt"
	"github.com/892294101/dds-spfile"
	"github.com/892294101/dds/ddslog"
	"os"
)

func main() {

	log, err := ddslog.InitDDSlog("HTTK_0002")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pfile, err := ddsspfile.LoadSpfile(fmt.Sprintf("%s.desc", "HTTK_0001"), ddsspfile.UTF8, log, ddsspfile.GetMySQLName(), ddsspfile.GetExtractName())
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
