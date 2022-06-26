package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

type Encoding uint

const (
	utf8Default Encoding = iota
	UTF8
	ISO88591
)

type spfileBaseInfo struct {
	enc         Encoding //文件字符集
	file        string   //文件
	log         *logrus.Logger
	dbType      string
	processType string
}

func LoadSpfile(filePath string, enc Encoding, log *logrus.Logger, dbType string, processType string) (*Spfile, error) {
	if len(filePath) == 0 {
		return nil, errors.New(fmt.Sprintf("Parameter file path must be specified"))
	}
	fh := new(spfileBaseInfo)
	fh.enc = enc
	fh.file = filePath
	fh.log = log
	fh.dbType = dbType
	fh.processType = processType
	return fh.LoadFile(fh)
}

func (f *spfileBaseInfo) LoadFile(fh *spfileBaseInfo) (*Spfile, error) {
	_, err := os.Stat(fh.file)
	if os.IsNotExist(err) {
		return nil, errors.Errorf("File not found: %s", fh.file)
	}
	return &Spfile{paramBaseInfo: fh, log: fh.log}, nil
}
