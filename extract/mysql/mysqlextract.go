package log

import (
	"context"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/sirupsen/logrus"

	"os"
	"time"
)

func main() {

	hlog := InitHlog()
	hlog.Infof("Initialize log file: %s", "log")

	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "",
	}
	cfg.Logger = hlog

	syncer := replication.NewBinlogSyncer(cfg)

	streamer, err := syncer.StartSync(mysql.Position{"1", 1})
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			// meet timeout
			continue
		}

		ev.Dump(os.Stdout)
	}
}
