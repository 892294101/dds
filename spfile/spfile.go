package spfile

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"myGithubLib/dds/utils"
	"os"
	"strings"
)

const (
	AnnotationPrefix = "--"
)

type Spfile struct {
	rawData       []string              // 文件原始数据
	paramBaseInfo *spfileBaseInfo       // 文件句柄
	log           *logrus.Logger        //日志系统
	ParamSet      map[string]Parameters // 参数集
	paramSetIndex []string              // 参数集的索引, 因为map不排序
	mustParams    []string              // 必须存在的参数
}

// 初始化数据库和进程必须存在的参数
func (s *Spfile) init() error {
	s.ParamSet = make(map[string]Parameters)

	switch {
	// MySQL extract进程必须存在的参数
	case s.paramBaseInfo.dbType == GetMySQLName() && s.paramBaseInfo.processType == GetExtractName():
		s.mustParams = append(s.mustParams, utils.ProcessType)
		s.mustParams = append(s.mustParams, utils.SourceDBType)
		s.mustParams = append(s.mustParams, utils.TrailDirType)
		s.mustParams = append(s.mustParams, utils.DiscardFileType)
		s.mustParams = append(s.mustParams, utils.DBOptionsType)
		s.mustParams = append(s.mustParams, utils.TableType)
	}

	return nil
}

func (s *Spfile) Production() error {
	if err := s.init(); err != nil {
		return err
	}
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
	return s.scanParams()
}

func (s *Spfile) scanParams() error {
	for _, params := range s.rawData {
		var pro Parameters
		switch {
		case utils.HasPrefixIgnoreCase(params, utils.ProcessType):
			if s.ParamSet[utils.ProcessType] == nil {
				pro = &ProcessBus
				if err := s.firstParams(pro, &params); err != nil {
					return err
				}
			} else {
				return errors.Errorf("%s configuration cannot be set repeatedly", utils.ProcessType)
			}
		case utils.HasPrefixIgnoreCase(params, utils.SourceDBType):
			if s.ParamSet[utils.SourceDBType] == nil {
				pro = &sourceDBSetBus
				if err := s.firstParams(pro, &params); err != nil {
					return err
				}
			} else {
				return errors.Errorf("%s configuration cannot be set repeatedly", utils.SourceDBType)
			}

		case utils.HasPrefixIgnoreCase(params, utils.TrailDirType):
			if s.ParamSet[utils.TrailDirType] == nil {
				pro = &trailDirBus
				if err := s.firstParams(pro, &params); err != nil {
					return err
				}
			} else {
				return errors.Errorf("%s configuration cannot be set repeatedly", utils.TrailDirType)
			}

		case utils.HasPrefixIgnoreCase(params, utils.DiscardFileType):
			if s.ParamSet[utils.DiscardFileType] == nil {
				pro = &DiscardFileBus
				if err := s.firstParams(pro, &params); err != nil {
					return err
				}
			} else {
				return errors.Errorf("%s configuration cannot be set repeatedly", utils.DiscardFileType)
			}

		case utils.HasPrefixIgnoreCase(params, utils.DBOptionsType):
			if s.ParamSet[utils.DBOptionsType] == nil {
				pro = &DBOptionsBus
				if err := s.firstParams(pro, &params); err != nil {
					return err
				}
			} else {
				return errors.Errorf("%s configuration cannot be set repeatedly", utils.DBOptionsType)
			}

		case utils.HasPrefixIgnoreCase(params, utils.TableType+" "):
			if s.ParamSet[utils.TableType] == nil {
				pro = &TableSetBus
				if err := s.firstParams(pro, &params); err != nil {
					return err
				}
			} else {
				pro = s.ParamSet[utils.TableType]
				if err := s.addParams(pro, &params); err != nil {
					return err
				}
			}
		case utils.HasPrefixIgnoreCase(params, utils.TableExcludeType):
			if s.ParamSet[utils.TableExcludeType] == nil {
				pro = &ExcludeTableSetBus
				if err := s.firstParams(pro, &params); err != nil {
					return err
				}
			} else {
				pro = s.ParamSet[utils.TableExcludeType]
				if err := s.addParams(pro, &params); err != nil {
					return err
				}
			}
		default:
			return errors.Errorf("Unknown parameter: %s", params)
		}

	}
	return s.registerMustParams()
}

func (s *Spfile) firstParams(pro Parameters, params *string) error {

	for Type, rawData := range pro.Registry() {
		if err := rawData.isType(params, &s.paramBaseInfo.dbType, &s.paramBaseInfo.processType); err != nil {
			return err
		}
		if err := rawData.parse(params); err != nil {
			return err
		}
		s.ParamSet[Type] = pro
		s.paramSetIndex = append(s.paramSetIndex, Type)

	}
	return nil
}

func (s *Spfile) addParams(pro Parameters, params *string) error {
	return pro.Add(params)
}

func (s *Spfile) registerMustParams() error {
	for _, paramType := range s.mustParams {
		switch paramType {
		case utils.DBOptionsType: // 对缺失的参数补充默认值
			_, ok := s.ParamSet[utils.DBOptionsType]
			if !ok {
				s.ParamSet[utils.DBOptionsType] = &DBOptionsBus
				for _, parameter := range s.ParamSet[utils.DBOptionsType].Registry() {
					if err := parameter.initDefault(); err != nil {
						return err
					}
				}
				s.paramSetIndex = append(s.paramSetIndex, paramType)
			}
		default:
			_, ok := s.ParamSet[paramType]
			if !ok {
				return errors.Errorf("The %s parameter must be configured", paramType)
			}
		}

	}
	return nil
}

func (s *Spfile) PutParamsText() {
	for _, index := range s.paramSetIndex {
		res := s.ParamSet[index].ListParamText()
		fmt.Println(res)
	}
}

