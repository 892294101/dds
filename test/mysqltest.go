package main

import (
	"fmt"
	"github.com/892294101/dds/ddslog"
	"github.com/892294101/go-mysql/canal"
	"github.com/892294101/go-mysql/mysql"
	"github.com/892294101/go-mysql/replication"
	"github.com/892294101/parser/ast"
	"math/rand"
	"os"
	"time"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnTableDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, tabDDL interface{}) error {

	/*
	*ast.CreateTableStmt
	*ast.AlterTableStmt:
	*ast.DropTableStmt:
	*ast.RenameTableStmt:
	*ast.TruncateTableStmt:
	 */
	switch td := tabDDL.(type) {
	case *ast.CreateTableStmt:
		fmt.Println("OnTableDDL: ", td.Text())
	case *ast.AlterTableStmt:
		fmt.Println("OnTableDDL: ", td.Text())
	case *ast.DropTableStmt:
		fmt.Println("OnTableDDL: ", td.Text())
	case *ast.RenameTableStmt:
		fmt.Println("OnTableDDL: ", td.Text())
	case *ast.TruncateTableStmt:
		fmt.Println("OnTableDDL: ", td.Text())

	}
	return nil
}

func (h *MyEventHandler) OnDataBaseDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, dbDDL interface{}) error {
	/*
	*ast.CreateDatabaseStmt
	*ast.DropDatabaseStmt
	*ast.AlterDatabaseStmt
	*
	 */

	switch dd := dbDDL.(type) {
	case *ast.CreateDatabaseStmt:
		fmt.Println("OnDataBaseDDL: ", dd.Text())
	case *ast.DropDatabaseStmt:
		fmt.Println("OnDataBaseDDL: ", dd.Text())
	case *ast.AlterDatabaseStmt:
		fmt.Println("OnDataBaseDDL: ", dd.Text())

	}

	return nil
}

func (h *MyEventHandler) OnIndexDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, iDDL interface{}) error {
	switch dd := iDDL.(type) {
	case *ast.CreateIndexStmt:
		fmt.Println("OnIndexDDL: ", dd.Text())
	case *ast.DropIndexStmt:
		fmt.Println("OnIndexDDL: ", dd.Text())

	}
	return nil
}

func (h *MyEventHandler) OnViewDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, vDDL interface{}) error {
	switch dd := vDDL.(type) {
	case *ast.CreateViewStmt:
		fmt.Println("OnViewDDL", dd.Text())
		fmt.Println("OnViewDDL", dd.Text())

	}
	return nil
}

func (h *MyEventHandler) OnSequenceDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, sDDL interface{}) error {
	switch sd := sDDL.(type) {
	case *ast.CreateSequenceStmt:
		fmt.Println("OnSequenceDDL", sd.Text())
	case *ast.DropSequenceStmt:
		fmt.Println("OnSequenceDDL", sd.Text())
	case *ast.AlterSequenceStmt:
		fmt.Println("OnSequenceDDL", sd.Text())

	}
	return nil
}

func (h *MyEventHandler) OnUserDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, uDDL interface{}) error {
	switch ud := uDDL.(type) {
	case *ast.CreateUserStmt:
		fmt.Println("OnUserDDL", ud.Text())
	case *ast.AlterUserStmt:
		fmt.Println("OnUserDDL", ud.Text())
	case *ast.DropUserStmt:
		fmt.Println("OnUserDDL", ud.Text())
	case *ast.RenameUserStmt:
		fmt.Println("OnUserDDL", ud.Text())
	}
	return nil
}

func (h *MyEventHandler) OnGrantDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, gDDL interface{}) error {
	switch gd := gDDL.(type) {
	case *ast.GrantProxyStmt:
		fmt.Println("OnGrantDDL", gd.Text())
	case *ast.GrantRoleStmt:
		fmt.Println("OnGrantDDL", gd.Text())
	case *ast.GrantStmt:
		fmt.Println("OnGrantDDL", gd.Text())
	}
	return nil
}

// 没用
/*func (h *MyEventHandler) OnGTID(m mysql.GTIDSet) error {
	fmt.Println("GTIDSet: ", m.String())

	return nil
}*/
// 没用
/*func (h *MyEventHandler) OnXID(m mysql.Position) error {
	fmt.Println("OnXID: ", m.String())
	return nil
}
*/
// 没用
/*func (h *MyEventHandler) OnRotate(re *replication.RotateEvent) error {
	fmt.Println("OnRotate: ", re.Position, re.NextLogName)
	return nil
}
*/
// 没用
/*func (h *MyEventHandler) OnTableChanged(schema string, table string) error {
	fmt.Println("OnTableChanged: ", schema, table)
	return nil
}*/
// 没用
/*func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}*/
// 此位置很重要
func (h *MyEventHandler) OnPosSynced(m mysql.Position, g mysql.GTIDSet, b bool) error {
	fmt.Println("=======================================")
	//每一次操作都会输出
	fmt.Printf("OnPosSynced: m.Name:%s m.Pos:%d \n", m.Name, m.Pos)

	if g != nil {
		fmt.Printf("OnPosSynced sync : %s %v\n", g, b)
	}

	return nil
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	fmt.Printf("Header: %d %v %v %v %v\n", e.Header.Timestamp, e.Header.EventSize, e.Header.LogPos, e.Header.Flags, e.Header.ServerID)
	fmt.Printf("OnRow: Table.Schema: %s Table.Name: %s\n", e.Table.Schema, e.Table.Name)
	fmt.Print("\n")

	for i, column := range e.Table.Columns {
		fmt.Println("OnRow: Table.Columns: ", i, " ", column)
	}

	fmt.Println("-----------------------------------------------------------------------")
	for _, row := range e.Rows {
		fmt.Printf("action: %s\n", e.Action)
		for i2, i3 := range row {
			fmt.Printf("column number: %d Data: %v\n", i2, i3)
		}

	}

	fmt.Println("-----------------------------------------------------------------------")

	return nil
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
	c.MaxReconnectAttempts = 3
	c.HeartbeatPeriod = time.Second
	c.UseDecimal = true

	ch, err := canal.NewCanal(c)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Register a handler to handle RowsEvent
	ch.SetEventHandler(&MyEventHandler{})

	fn := mysql.Position{Name: fmt.Sprintf("mysql-bin.%06d", 2), Pos: 660}

	ch.RunFrom(fn)

	/*pfile, err := spfile.LoadSpfile("D:\\workspace\\gowork\\src\\github.com/892294101\\dds\\build\\param\\httk_0001.desc",
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
