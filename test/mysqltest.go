package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"myGithubLib/dds/extract/mysql/spfile"
	"os"
)

func main() {

	//m := treebidimap.NewWith(utils.IntComparator, utils.StringComparator)

	p, err := spfile.LoadSpfile("D:\\workspace\\gowork\\src\\myGithubLib\\dds\\build\\param\\httk_0001.desc",
		spfile.UTF8,
		logrus.New(),
		spfile.GetMySQLName(),
		spfile.GetExtractName())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := p.Production(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p.PutParamsText()

	/*log, err := ddslog.InitDDSlog()
	if err != nil {
		fmt.Fprint(os.Stderr, "DDS log error: %s", err)
		os.Exit(1)
	}

	cfg := replication.BinlogSyncerConfig{
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "",
		Logger:   log,
	}

	syncer := replication.NewBinlogSyncer(cfg)

	streamer, err := syncer.StartSync(mysql.Position{Name: "1", Pos: 1})
	if err != nil {
		log.Fatalf("%s", err)
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			continue
		}
		ev.Dump(os.Stdout)
	}*/

}
