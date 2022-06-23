package spfile

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	AnnotationPrefix = "--"
)

type Spfile struct {
	rawData  []string       // 文件原始数据
	fileInfo *fileHandle    // 文件句柄
	log      *logrus.Logger //日志系统
}

func (s *Spfile) Production() (*Spfile, error) {
	f, err := os.Open(s.fileInfo.file)
	if err != nil {
		return nil, errors.Errorf("Failed to open parameter file %s: %s", s.fileInfo.file, err)
	}
	reader := bufio.NewScanner(f)
	for reader.Scan() {
		val := strings.TrimSpace(reader.Text())
		if !strings.HasPrefix(val, AnnotationPrefix) && val != "" {
			s.rawData = append(s.rawData, val)
		}
	}
	for _, datum := range s.rawData {

		fmt.Println(strings.TrimSpace(datum))
	}
	return s, nil
}
