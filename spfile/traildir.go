package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

var (
	TrailDirType    = "TRAILDIR"
	TrailDirRegular = "(^)(?i:(" + TrailDirType + "))(\\s+)(*)($)"
)

type TrailDir struct {
	SupportParams map[string]map[string]string
	ParamPrefix   string
	Dir           string
}

func (t *TrailDir) Put() {
	fmt.Println("traildir Info:", t.ParamPrefix, t.Dir)
}

// 初始化参数可以支持的数据库和进程
func (t *TrailDir) Init() {
	t.SupportParams = map[string]map[string]string{
		MySQL: {
			Extract:  Extract,
			Replicat: Replicat,
		},
		MariaDB: {
			Extract:  Extract,
			Replicat: Replicat,
		},
	}
}

func (t *TrailDir) IsType(raw *string, dbType *string, processType *string) error {
	t.Init()
	_, ok := t.SupportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

func (t *TrailDir) Parse(raw *string) error {
	matched, _ := regexp.MatchString(TrailDirRegular, *raw)
	if matched == true {
		rd := strings.Split(*raw, " ")
		t.ParamPrefix = rd[0]
		t.Dir = rd[1]
		return nil
	}

	return errors.Errorf("%s parameter parsing failed: %s", TrailDirType, *raw)
}

type trailDirSet struct{}

var trailDirBus trailDirSet

func (t *trailDirSet) Init() {

}

func (t *trailDirSet) Registry() map[string]Parameter {
	return map[string]Parameter{TrailDirType: &Process{}}
}
