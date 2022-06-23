package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"myGithubLib/spfile"
)

func main() {

	p, err := spfile.LoadSpfile("D:\\workspace\\gowork\\src\\myGithubLib\\dds\\build\\param\\httk_0001.desc",spfile.UTF8,logrus.New())
	if err != nil {
		fmt.Println(err)
	}

	p.Production()

}
