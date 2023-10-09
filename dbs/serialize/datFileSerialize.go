package serialize

import (
	"bytes"
	"encoding/binary"
	"github.com/892294101/dds/dbs/utils"
	"github.com/golang/snappy"
	"github.com/pkg/errors"
)

const (
	FileBegin uint8 = iota
	FileEnd
	EOF_FILE       byte = 0x05
	DATA_FILE_HEAD byte = 0x02
)

// 事务提交和结束
type FileMarkV1 struct {
	WriteMark *uint8
	WriteTime *uint64
	fm        *FileMark
}

// 事务提交和结束
type FileMark struct {
	WriteMark uint8
	WriteTime uint64
}

func (f *FileMarkV1) Init(v *FileMark) error {
	if v == nil {
		return errors.Errorf("Serialization event cannot be empty")
	}
	f.fm = v
	return nil
}

func (f *FileMarkV1) InitBuffer() error {
	return nil
}

func (f *FileMarkV1) EncodeFileMark(Buffer *bytes.Buffer) error {
	var err error

	if err = binary.Write(Buffer, binary.LittleEndian, DATA_FILE_HEAD); err != nil {
		return err
	}

	if err = binary.Write(Buffer, binary.LittleEndian, f.fm.WriteMark); err != nil {
		return err
	}

	if err = binary.Write(Buffer, binary.LittleEndian, f.fm.WriteTime); err != nil {
		return err
	}

	if err = binary.Write(Buffer, binary.LittleEndian, EOF_FILE); err != nil {
		return err
	}
	return nil
}
func (f *FileMarkV1) EncodeData() ([]byte, error) {
	Buffer := utils.DataRowsBufferGet()
	defer utils.DataRowsBufferPut(Buffer)
	if err := f.EncodeFileMark(Buffer); err != nil {
		return nil, err
	}
	return Buffer.Bytes(), nil
}

func (f *FileMarkV1) Compress(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.Errorf("FileMark Compression input cannot be empty")
	}
	return snappy.Encode(nil, data), nil
}

func (f *FileMarkV1) Encode(data interface{}, compress bool) ([]byte, error) {
	switch v := data.(type) {
	case *FileMark:
		if err := f.Init(v); err != nil {
			return nil, err
		}
		data, err := f.EncodeData()
		if err != nil {
			return nil, err
		}
		if compress {
			return f.Compress(data)
		}
		return data, nil
	default:
		return nil, errors.Errorf("Unknown serialization event")
	}
}
