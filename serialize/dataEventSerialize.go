package serialize

import (
	"encoding/binary"
	"github.com/892294101/dds/utils"
	"github.com/892294101/go-mysql/canal"
	"github.com/golang/snappy"
	"github.com/pkg/errors"
)

var (
	DATA_ROW_HEAD = []byte{0x01, 0x00} // 第一个字节表示这是一个数据行事件，第二个字节是用于补全数据字节头（OK_HEAD），因为它传递是被replication给截掉了。
)

// 数据行事件==================================================================
type DataRowEvent struct {
	RowEvent *canal.RowsEvent
	RawData  []byte
}

type DataEventV1 struct {
	RowEvent *DataRowEvent
}

func (d *DataEventV1) Init(v *DataRowEvent) error {
	if v == nil {
		return errors.Errorf("Serialization event cannot be empty")
	}
	d.RowEvent = v
	return nil
}

func (d *DataEventV1) InitBuffer() error {
	return nil
}

func (d *DataEventV1) EncodeData() ([]byte, error) {
	Buffer := utils.DataRowsBufferGet()
	defer utils.DataRowsBufferPut(Buffer)

	if err := binary.Write(Buffer, binary.LittleEndian, DATA_ROW_HEAD); err != nil {
		return nil, errors.Errorf("DATA_ROW_HEAD binary write error: %v", err)
	}

	if err := binary.Write(Buffer, binary.LittleEndian, d.RowEvent.RawData); err != nil {
		return nil, errors.Errorf("RawData binary write error: %v", err)
	}

	if err := binary.Write(Buffer, binary.LittleEndian, EOF_DATA); err != nil {
		return nil, errors.Errorf("EOF_DATA binary write error: %v", err)
	}

	return Buffer.Bytes(), nil
}

func (d *DataEventV1) Clear() error {
	d.RowEvent = nil
	return nil
}

func (d *DataEventV1) Compress(data *[]byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("DataEvent Compression input cannot be empty")
	}
	return snappy.Encode(nil, *data), nil
}

func (d *DataEventV1) Encode(data interface{}, compress bool) ([]byte, error) {
	switch v := data.(type) {
	case *DataRowEvent:
		if err := d.Init(v); err != nil {
			return nil, err
		}
		data, err := d.EncodeData()
		if err != nil {
			return nil, err
		}
		_ = d.Clear()

		if compress {
			return d.Compress(&data)
		}
		//fmt.Println(len(data))
		return data, nil
	default:
		return nil, errors.Errorf("Unknown serialization event")
	}
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}
