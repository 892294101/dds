package main

import (
	"context"
	"fmt"
	"github.com/892294101/dds/ddslog"
	"github.com/892294101/go-mysql/mysql"
	"github.com/892294101/go-mysql/replication"
	"os"
	"time"
)

func main() {

	log, err := ddslog.InitDDSlog("HTTK_0002")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "10.130.41.246",
		Port:     3306,
		User:     "root",
		Password: "Admin!@123",
		Logger:   log,
	}
	syncer := replication.NewBinlogSyncer(cfg)

	streamer, _ := syncer.StartSync(mysql.Position{"mysql-bin.000274", 2163725})

	d, err := replication.NewRawDataDecode(log)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	OK := []byte{0x00}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			// meet timeout
			continue
		}
		c, e := d.Decode(append(OK, ev.RawData...))
		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}
		fmt.Println(c.Header)

	}
}
