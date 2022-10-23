package serialize

import (
	"bytes"
	"encoding/binary"
	"github.com/892294101/dds/utils"
	"github.com/892294101/go-mysql/mysql"
	"github.com/892294101/go-mysql/replication"
	"github.com/golang/snappy"
	"github.com/pkg/errors"
)

var DefaultXid uint64 = 0

const (
	TransBegin = iota
	TransCommit

	TRANSACTION_HEAD byte = 0x03
	EOF_TRANSACTION  byte = 0x00
)

// 事务提交和结束
type TransactionEvent struct {
	Xid       *replication.XIDEvent // 事务提交是的POS和XID
	Pos       *mysql.Position       // 当时开始时对应事件的pos，MySQL开始事务没有xid
	TransType int                   // 事务启动或结束标识
}

// 事务提交和结束
type TransactionV1 struct {
	TransType uint8   // 事务启动或结束标识
	Xid       *uint64 // 事务提交时XID
	Fn        *uint64 // 文件号
	Pos       *uint64 // 当时开始时对应事件的pos，MySQL开始事务没有xid

}

func (t *TransactionV1) Init(v *TransactionEvent) error {
	if v == nil {
		return errors.Errorf("Transaction Serialization event cannot be empty")
	}
	switch v.TransType {
	case TransCommit:
		t.Xid = &v.Xid.XID
	case TransBegin:
		t.Xid = &DefaultXid
	}
	t.TransType = uint8(v.TransType)
	f, p, err := utils.ConvertPositionToNumber(v.Pos)
	if err != nil {
		return err
	}
	t.Fn = f
	t.Pos = p
	return nil
}

func (t *TransactionV1) InitBuffer() error {
	return nil
}

func (t *TransactionV1) EncodeTransaction(Buffer *bytes.Buffer) error {
	var err error
	if err = binary.Write(Buffer, binary.LittleEndian, TRANSACTION_HEAD); err != nil {
		return err
	}

	if err = binary.Write(Buffer, binary.LittleEndian, t.TransType); err != nil {
		return err
	}

	if err = binary.Write(Buffer, binary.LittleEndian, t.Fn); err != nil {
		return err
	}

	if err = binary.Write(Buffer, binary.LittleEndian, t.Pos); err != nil {
		return err
	}

	if err = binary.Write(Buffer, binary.LittleEndian, EOF_TRANSACTION); err != nil {
		return err
	}
	return nil
}

func (t *TransactionV1) EncodeData() ([]byte, error) {
	Buffer := utils.DataRowsBufferGet()
	defer utils.DataRowsBufferPut(Buffer)

	if err := t.EncodeTransaction(Buffer); err != nil {
		return nil, err
	}

	return Buffer.Bytes(), nil
}

func (t *TransactionV1) Compress(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.Errorf("Transaction Compression input cannot be empty")
	}
	return snappy.Encode(nil, data), nil
}

func (t *TransactionV1) Encode(data interface{}, compress bool) ([]byte, error) {
	switch v := data.(type) {
	case *TransactionEvent:
		if err := t.Init(v); err != nil {
			return nil, err
		}
		data, err := t.EncodeData()
		if err != nil {
			return nil, err
		}
		if compress {
			return t.Compress(data)
		}
		return data, nil
	default:
		return nil, errors.Errorf("Unknown serialization event")
	}
}
