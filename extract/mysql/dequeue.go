package oramysql

import (
	"github.com/892294101/dds-utils"
	"github.com/892294101/dds/serialize"
	"github.com/pkg/errors"
	"time"
	"unsafe"
)

const (
	DataEvent = iota
	TransactionEvent
)

func (e *ExtractEvent) WriteToFile(v []byte) error {
	defer dds_utils.ErrorCheckOfRecover(e.WriteToFile, e.log)
	if err := e.fileWrite.Write(v); err != nil {
		return err
	}

	if err := e.fileWrite.Sync(); err != nil {
		return err
	}

	return nil
}

func (e *ExtractEvent) WriteCache(v interface{}) error {
	defer dds_utils.ErrorCheckOfRecover(e.WriteCache, e.log)
	switch d := v.(type) {
	case *serialize.DataRowEvent:
		sch := *(*string)(unsafe.Pointer(&d.RowEvent.RowsEvent.Table.Schema))
		tab := *(*string)(unsafe.Pointer(&d.RowEvent.RowsEvent.Table.Table))
		// 过滤表是否存储
		ok, err := e.pFile.DMLFilter(&sch, &tab)
		if err != nil {
			return err
		}

		if ok {
			e.log.Debugf("write %v.%v event row data to file", sch, tab)
			r, err := e.Encode(d)
			if err != nil {
				e.log.Errorf("%v", err)
				e.ClearProcessInfoFile()
				e.CloseAll()
			}
			if len(e.TranHead.transactionHead) > 0 {
				e.log.Debugf("transaction head write to file from transaction cache")
				if err := e.WriteToFile(e.TranHead.transactionHead); err != nil {
					e.log.Errorf("%v", err)
					e.ClearProcessInfoFile()
					e.CloseAll()
				}
				e.log.Debugf("transaction head is write success and cache transaction header is clear")
				e.TranHead.transactionHead = e.TranHead.transactionHead[0:0]
			}
			e.TranHead.eventCount += 1
			e.log.Debugf("current event number: %d (data rows: %v)", e.TranHead.eventCount, len(d.RowEvent.RowsEvent.Rows))
			if err := e.WriteToFile(r); err != nil {
				return err
			} else {
				e.log.Debugf("row event write successfully")
			}
		} else {
			e.log.Debugf("%v.%v data discard", sch, tab)
		}
	case *serialize.TransactionEvent:
		if e.TranHead == nil {
			e.TranHead = new(TranHeadFilter)
		}
		switch d.TransType {
		case serialize.TransBegin:
			e.log.Debugf("transaction head received: %v:%v", d.Pos.Name, d.Pos.Pos)
			// 如果是开始事务则暂存在内存中
			data, err := e.Encode(d)
			if err != nil {
				e.log.Errorf("%v", err)
				e.ClearProcessInfoFile()
				e.CloseAll()
			}
			e.TranHead.transactionHead = data
			// 当事件时begin事务时，把当前时间写入到检查点元数据文件
			if err := e.md.SetTransactionBeginTime(uint64(time.Now().Unix())); err != nil {
				e.log.Warnf("unable to set transaction begin timestamp: %s", err)
			}
		case serialize.TransCommit:
			e.log.Debugf("transaction tail received: %v:%v. transaction xid: %v", d.Pos.Name, d.Pos.Pos, d.Xid.XID)
			// 更新log file number和position
			s, p, err := dds_utils.ConvertPositionToNumber(d.Pos)
			if err != nil {
				e.log.Errorf("%v", err)
				e.ClearProcessInfoFile()
				e.CloseAll()
			}
			err = e.md.SetPosition(*s, *p)
			if err != nil {
				e.log.Errorf("%v", err)
				e.ClearProcessInfoFile()
				e.CloseAll()
			}

			// 当检查点信息完成后设置begin事务时间戳为0
			if err := e.md.SetTransactionBeginTime(0); err != nil {
				e.log.Warnf("unable to set transaction begin timestamp: %s", err)
			}

			if e.TranHead.eventCount > 0 {
				e.log.Debugf("transaction tail received. write transaction tail to file")
				data, err := e.Encode(d)
				if err != nil {
					e.log.Errorf("%v", err)
					e.ClearProcessInfoFile()
					e.CloseAll()
				}
				err = e.WriteToFile(data)
				if err != nil {
					e.log.Errorf("%v", err)
					e.ClearProcessInfoFile()
					e.CloseAll()
				}
				// 当存在所需要的事件是，事务头写入完成后，会清除事务头，所以事务结束后，只需要清空时间数据即可
				e.TranHead.eventCount = 0
				e.log.Debugf("transaction tail write completed. clear the number of event row number")
			} else {
				// 当事件数据不存在是，说明没有row数据是需要的，那么当接收到事务结束后，只清空事务头即可
				e.TranHead.transactionHead = e.TranHead.transactionHead[0:0]
				e.log.Debugf("transaction tail received. clear transaction header because the transaction is empty and commit")
			}

		}
	}
	return nil
}

func (e *ExtractEvent) Encode(v interface{}) ([]byte, error) {
	defer dds_utils.ErrorCheckOfRecover(e.Encode, e.log)
	switch d := v.(type) {
	case *serialize.DataRowEvent:
		return e.serialize[DataEvent].Encode(d, false)
	case *serialize.TransactionEvent:
		return e.serialize[TransactionEvent].Encode(d, false)
	default:
		return nil, errors.Errorf("encode event row unknown")
	}
}

func (e *ExtractEvent) ProcessCacheData() {
	defer dds_utils.ErrorCheckOfRecover(e.ProcessCacheData, e.log)
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for range ticker.C {
			e.log.Infof("data row cache queue size: %v", len(e.replicationStream.event.cache.queue))
		}
	}()

	//FristRows := make(chan int, 1)
	for {
		select {
		case v := <-e.stopSrv.power:
			// 当接收到终止信号后，执行停止操作，但是当数据不完整的情况下正常模式是无法停止的，除非强制停止。
			switch v {
			case STOP:
				if e.fileWrite.Dirty == false {
					e.ClearProcessInfoFile()
					e.CloseAll()
				} else {
					e.log.Infof("%v", "A termination signal was received, but could not be stopped because the transaction was not completed")
				}
			case FORCE:
				e.log.Warnf("%v", "Forced stop signal received")
				e.ClearProcessInfoFile()
				e.CloseAll()
			}

		case v := <-e.replicationStream.event.cache.queue:
			switch data := v.(type) {
			case *serialize.DataRowEvent:
				if err := e.WriteCache(data); err != nil {
					e.log.Errorf("%v", err)
					e.ClearProcessInfoFile()
					e.CloseAll()
				}

			case *serialize.TransactionEvent:
				if err := e.WriteCache(data); err != nil {
					e.log.Errorf("%v", err)
					e.ClearProcessInfoFile()
					e.CloseAll()
				}
			}
		}
	}

}
