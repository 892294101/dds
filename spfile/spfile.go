package spfile

import (
	"bufio"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"myGithubLib/dds/extract/mysql/utils"
	"os"
	"strings"
)

const (
	AnnotationPrefix = "--"
)

type Spfile struct {
	rawData         []string       // 文件原始数据
	fileInfo        *fileHandle    // 文件句柄
	log             *logrus.Logger //日志系统
	instructionsSet map[string]Parameters
}

func (s *Spfile) Production() error {
	f, err := os.Open(s.fileInfo.file)
	if err != nil {
		return errors.Errorf("Failed to open parameter file %s: %s", s.fileInfo.file, err)
	}
	reader := bufio.NewScanner(f)
	for reader.Scan() {
		val := strings.TrimSpace(reader.Text())
		if !strings.HasPrefix(val, AnnotationPrefix) && val != "" {
			s.rawData = append(s.rawData, val)
		}
	}
	s.instructionsSet = make(map[string]Parameters)

	for _, params := range s.rawData {
		var pro Parameters
		var CallType string
		switch {
		case utils.HasPrefixIgnoreCase(params, ProcessType):
			CallType = ProcessType
			pro = &processSet{}
		case utils.HasPrefixIgnoreCase(params, SourceDBType):
			CallType = SourceDBType
			pro = &sourceDBSet{}
		default:
			return errors.Errorf("Unknown parameter: %s", params)
		}

		for _, parameter := range pro.Registry() {
			if err := parameter.Parse(params); err != nil {
				return err
			}
			s.instructionsSet[CallType] = pro
			parameter.Put()
		}

	}
	return nil
}
