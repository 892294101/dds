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
	rawData         []string        // 文件原始数据
	paramBaseInfo   *spfileBaseInfo // 文件句柄
	log             *logrus.Logger  //日志系统
	instructionsSet map[string]Parameters
}

func (s *Spfile) Production() error {
	f, err := os.Open(s.paramBaseInfo.file)
	if err != nil {
		return errors.Errorf("Failed to open parameter file %s: %s", s.paramBaseInfo.file, err)
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
		case utils.HasPrefixIgnoreCase(params, utils.ProcessType):
			CallType = utils.ProcessType
			pro = &ProcessBus
		case utils.HasPrefixIgnoreCase(params, utils.SourceDBType):
			CallType = utils.SourceDBType
			pro = &sourceDBSetBus
		case utils.HasPrefixIgnoreCase(params, utils.TrailDirType):
			CallType = utils.TrailDirType
			pro = &trailDirBus
		default:
			return errors.Errorf("Unknown parameter: %s", params)
		}
		for _, rawData := range pro.Registry() {
			if err := rawData.IsType(&params, &s.paramBaseInfo.dbType, &s.paramBaseInfo.processType); err != nil {
				return err
			}
			if err := rawData.Parse(&params); err != nil {
				return err
			}
			s.instructionsSet[CallType] = pro
			rawData.Put()
		}

	}
	return nil
}
