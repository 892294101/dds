package oramysql

import (
	"fmt"
	"github.com/892294101/dds-utils"
	"github.com/892294101/dds/serialize"
	"github.com/892294101/go-mysql/canal"
	"github.com/892294101/go-mysql/mysql"
	"github.com/892294101/go-mysql/replication"
	"github.com/892294101/parser/ast"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	MaxCache = 2000
)

type queue struct {
	queue       chan interface{} // 缓存队列
	currentSize uint32
	maxSize     uint32
	log         *logrus.Logger // 日志记录器
}

func (q *queue) initQueue(log *logrus.Logger) error {
	if log == nil {
		return errors.Errorf("Initialization cache logger cannot be empty")
	}
	q.log = log
	q.queue = make(chan interface{}, MaxCache)
	q.maxSize = MaxCache
	q.currentSize = 0
	return nil
}

func (q *queue) Enqueue(data interface{}, size uint32) {

	defer ddsutils.ErrorCheckOfRecover(q.Enqueue, q.log)
	if q.currentSize+size < q.maxSize {
		q.queue <- data
		q.currentSize += size
	} else {
		for {
			t := time.NewTicker(time.Second)
			select {
			case <-t.C:
				if q.currentSize+size < q.maxSize {
					if data == nil {
						return
					}
					q.queue <- data
					q.currentSize += size
					t.Stop()
					return
				} else {
					q.log.Warnf("Data cache is full")
				}
			}

		}
	}

}

// 重新实现数据流接口
type EventExtract struct {
	canal.DummyEventHandler
	cache *queue
}

func (h *EventExtract) OnTableDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, d interface{}, rawData []byte) error {
	/*
	*ast.CreateTableStmt
	*ast.AlterTableStmt:
	*ast.DropTableStmt:
	*ast.RenameTableStmt:
	*ast.TruncateTableStmt:
	 */

	switch td := d.(type) {
	case *ast.CreateTableStmt:
		fmt.Println("OnTableDDL: ", td.Text(), nextPos.String())
	case *ast.AlterTableStmt:
		fmt.Println("OnTableDDL: ", td.Text(), nextPos.String())
	case *ast.DropTableStmt:
		fmt.Println("OnTableDDL: ", td.Text(), nextPos.String())
	case *ast.RenameTableStmt:
		fmt.Println("OnTableDDL: ", td.Text(), nextPos.String())
	case *ast.TruncateTableStmt:
		fmt.Println("OnTableDDL: ", td.Text(), nextPos.String())

	}
	return nil
}

func (h *EventExtract) OnDataBaseDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, d interface{}, rawData []byte) error {
	/*
	*ast.CreateDatabaseStmt
	*ast.DropDatabaseStmt
	*ast.AlterDatabaseStmt
	*
	 */

	switch dd := d.(type) {
	case *ast.CreateDatabaseStmt:
		fmt.Println("OnDataBaseDDL: ", dd.Text(), nextPos.String())
	case *ast.DropDatabaseStmt:
		fmt.Println("OnDataBaseDDL: ", dd.Text(), nextPos.String())
	case *ast.AlterDatabaseStmt:
		fmt.Println("OnDataBaseDDL: ", dd.Text(), nextPos.String())

	}

	return nil
}

func (h *EventExtract) OnIndexDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, d interface{}, rawData []byte) error {
	switch dd := d.(type) {
	case *ast.CreateIndexStmt:
		fmt.Println("OnIndexDDL: ", dd.Text(), nextPos.String())
	case *ast.DropIndexStmt:
		fmt.Println("OnIndexDDL: ", dd.Text(), nextPos.String())

	}
	return nil
}

func (h *EventExtract) OnViewDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, d interface{}, rawData []byte) error {
	switch dd := d.(type) {
	case *ast.CreateViewStmt:
		fmt.Println("OnViewDDL", dd.Text(), nextPos.String())
		fmt.Println("OnViewDDL", dd.Text(), nextPos.String())

	}
	return nil
}

func (h *EventExtract) OnSequenceDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, d interface{}, rawData []byte) error {
	switch sd := d.(type) {
	case *ast.CreateSequenceStmt:
		fmt.Println("OnSequenceDDL", sd.Text(), nextPos.String())
	case *ast.DropSequenceStmt:
		fmt.Println("OnSequenceDDL", sd.Text(), nextPos.String())
	case *ast.AlterSequenceStmt:
		fmt.Println("OnSequenceDDL", sd.Text(), nextPos.String())

	}
	return nil
}

func (h *EventExtract) OnUserDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, d interface{}, rawData []byte) error {
	switch ud := d.(type) {
	case *ast.CreateUserStmt:
		fmt.Println("OnUserDDL", ud.Text(), nextPos.String())
	case *ast.AlterUserStmt:
		fmt.Println("OnUserDDL", ud.Text(), nextPos.String())
	case *ast.DropUserStmt:
		fmt.Println("OnUserDDL", ud.Text(), nextPos.String())
	case *ast.RenameUserStmt:
		fmt.Println("OnUserDDL", ud.Text(), nextPos.String())
	}
	return nil
}

func (h *EventExtract) OnGrantDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent, d interface{}, rawData []byte) error {
	/*switch gd := gDDL.(type) {
	case *ast.GrantProxyStmt:
		fmt.Println("OnGrantDDL", gd.Text())
	case *ast.GrantRoleStmt:
		fmt.Println("OnGrantDDL", gd.Text())
	case *ast.GrantStmt:
		fmt.Println("OnGrantDDL", gd.Text())
	}*/
	return nil
}

// 事务开始
func (h *EventExtract) OnTransaction(m mysql.Position, queryEvent *replication.QueryEvent, d interface{}, rawData []byte) error {
	h.cache.queue <- &serialize.TransactionEvent{Xid: nil, Pos: &m, TransType: serialize.TransBegin}
	return nil
}

func (h *EventExtract) OnGTID(m mysql.GTIDSet, rawData []byte) error {
	return nil
}

// 事务结束
func (h *EventExtract) OnXID(m mysql.Position, e *replication.XIDEvent, rawData []byte) error {
	h.cache.queue <- &serialize.TransactionEvent{Xid: e, Pos: &m, TransType: serialize.TransCommit}
	return nil
}

// 没用

/*
当MySQL切换至新的binlog文件的时候，MySQL会在旧的binlog文件中写入一个ROTATE_EVENT，
表示新的binlog文件的文件名，以及第一个偏移地址。当在数据库中执行FLUSH LOGS语句或者binlog文件的大小超过max_binlog_size就会切换新的binlog文件
*/
func (h *EventExtract) OnRotate(re *replication.RotateEvent, rawData []byte) error {
	// fmt.Println("OnRotate: ", re.Position, string(re.NextLogName))
	return nil
}

func (h *EventExtract) String() string {
	return ""
}

// 此位置很重要
func (h *EventExtract) OnPosSynced(m mysql.Position, g mysql.GTIDSet, b bool) error {
	return nil
}

/*
如果update httk03 set id=4,name='wang4' 产生变化4行，那么e.RowsEvent将是4行。
如果begin update httk03 set id=4,name='wang4' 每次产生一行，那么e.RowsEvent将每次输出一行
*/
func (h *EventExtract) OnRow(e *canal.RowsEvent, rawData []byte) error {
	h.cache.queue <- &serialize.DataRowEvent{RowEvent: e, RawData: rawData}
	return nil
}
