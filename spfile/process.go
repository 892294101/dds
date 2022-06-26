package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

var (
	ProcessType    = "PROCESS"
	ProcessRegular = "(^)(?i:(" + ProcessType + "))(\\s+)((?:[A-Za-z0-9_]){4,12})($)"
)

type Process struct {
	SupportParams map[string]map[string]string
	ParamPrefix   string
	Name          string
}

func (p *Process) Put() {
	fmt.Println("process Info: ", p.ParamPrefix, p.Name)
}

// 初始化参数可以支持的数据库和进程
func (p *Process) Init() {
	p.SupportParams = map[string]map[string]string{
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

func (p *Process) IsType(raw *string, dbType *string, processType *string) error {
	p.Init()
	_, ok := p.SupportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

func (p *Process) Parse(raw *string) error {
	matched, _ := regexp.MatchString(ProcessRegular, *raw)
	if matched == true {
		rd := strings.Split(*raw, " ")
		p.ParamPrefix = rd[0]
		p.Name = rd[1]
		return nil
	}

	return errors.Errorf("%s parameter parsing failed: %s", ProcessType, *raw)
}

type processSet struct{}

var ProcessBus processSet

func (p *processSet) Init() {

}

func (p *processSet) Registry() map[string]Parameter {
	return map[string]Parameter{ProcessType: &Process{}}
}
