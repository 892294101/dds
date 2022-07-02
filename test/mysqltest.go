package main

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"math/rand"
	"myGithubLib/dds/extract/mysql/ddslog"
	"os"
	"time"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {

	fmt.Printf("Header: %d %v %v %v %v\n", e.Header.Timestamp, e.Header.EventSize, e.Header.LogPos, e.Header.Flags, e.Header.ServerID)
	fmt.Printf("OnRow: Table.Schema: %s Table.Name:%s  Table.PKColumns:%v Table.UnsignedColumns:%v\n", e.Table.Schema, e.Table.Name, e.Table.PKColumns, e.Table.UnsignedColumns)

	for i, column := range e.Table.Columns {
		fmt.Println("OnRow: Table.Columns: ", i, " ", column )
	}

	for i, index := range e.Table.Indexes {
		fmt.Println("OnRow: Table.Indexes: ", i, " ", index)
	}

	fmt.Println("OnRow: ", e.Action, " ", e.Rows)

	return nil
}
func (h *MyEventHandler) OnGTID(m mysql.GTIDSet) error {
	fmt.Println("GTIDSet: ", m.String())

	return nil
}
func (h *MyEventHandler) OnXID(m mysql.Position) error {
	fmt.Println("OnXID: ", m.String())
	return nil
}
func (h *MyEventHandler) OnRotate(re *replication.RotateEvent) error {
	fmt.Println("OnRotate: ", re.Position, re.NextLogName)
	return nil
}
func (h *MyEventHandler) OnTableChanged(schema string, table string) error {
	fmt.Println("OnTableChanged: ", schema, table)
	return nil
}

func (h *MyEventHandler) OnDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent) error {

	fmt.Println("=======================================")
	fmt.Println("OnDDL: ", queryEvent.ExecutionTime)
	fmt.Println("OnDDL: ", string(queryEvent.Schema))
	fmt.Println("OnDDL: ", string(queryEvent.Query))
	fmt.Println("OnDDL: ", nextPos.Name, "  ", nextPos.Pos)
	queryEvent.Dump(os.Stdout)


	//queryEvent.Dump(os.Stdout)

	return nil
}

func (h *MyEventHandler) OnPosSynced(m mysql.Position, g mysql.GTIDSet, b bool) error {
	fmt.Println("=======================================")
	//每一次操作都会输出
	fmt.Printf("OnPosSynced: m.Name:%s m.Pos:%d \n", m.Name, m.Pos)

	if g != nil {
		fmt.Printf("OnPosSynced: %s\n", g)
	}
	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func main() {
	log, err := ddslog.InitDDSlog()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init failed: %s", err)
		os.Exit(1)
	}

	c := new(canal.Config)

	c.Addr = "10.130.41.234:3306"
	c.User = "admin"
	c.Password = "admin"

	c.Charset = mysql.DEFAULT_CHARSET
	c.ServerID = uint32(rand.New(rand.NewSource(time.Now().Unix())).Intn(1000)) + 1001

	c.Flavor = "mysql"

	c.Dump.DiscardErr = true
	c.Dump.SkipMasterData = false

	c.Logger = log

	ch, err := canal.NewCanal(c)
	if err != nil {
		log.Fatalf("%s", err)
	}

	c.DisableRetrySync = true
	c.UseDecimal = true

	// Register a handler to handle RowsEvent
	ch.SetEventHandler(&MyEventHandler{})

	fn := mysql.Position{Name: fmt.Sprintf("mysql-bin.%06d", 1), Pos: 4492}

	ch.RunFrom(fn)
	// Start canal
	ch.Run()

/*
		 pfile, err := spfile.LoadSpfile("D:\\workspace\\gowork\\src\\myGithubLib\\dds\\build\\param\\httk_0001.desc",
			spfile.UTF8,
			log,
			spfile.GetMySQLName(),
			spfile.GetExtractName())
		if err != nil {
			log.Fatalf("%s", err)
		}

		if err := pfile.Production(); err != nil {
			log.Fatalf("%s", err)
		}

		ext := oramysql.NewMySQLSync()
		err = ext.InitSyncerConfig(log, pfile)
		if err != nil {
			log.Fatalf("%s", err)
		}

		err = ext.StartSyncToStream(1, 150)
		if err != nil {
			log.Fatalf("StartSyncToStream failed: %s", err)
		}*/

}
