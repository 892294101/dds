package main

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"myGithubLib/dds/extract/mysql/ddslog"
	"os"
	"time"
)

func main() {

	log, err := ddslog.InitDDSlog()
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
	}

}
