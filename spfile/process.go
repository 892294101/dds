package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/892294101/dds/utils"
	"regexp"
	"strings"
)

type ProcessInfo struct {
	name *string
}

func (p *ProcessInfo) GetName() *string { return p.name }

type Process struct {
	supportParams map[string]map[string]string
	paramPrefix   *string
	ProInfo       *ProcessInfo
}

func (p *Process) put() string {
	return fmt.Sprintf("%s %s", *p.paramPrefix, *p.ProInfo.name)
}

// 初始化参数可以支持的数据库和进程
func (p *Process) init() {
	p.supportParams = map[string]map[string]string{
		utils.MySQL: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
		utils.MariaDB: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
	}
}

func (p *Process) initDefault() error {
	return nil
}

func (p *Process) isType(raw *string, dbType *string, processType *string) error {
	p.init()
	_, ok := p.supportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

func (p *Process) parse(raw *string) error {
	matched, _ := regexp.MatchString(utils.ProcessRegular, *raw)
	if matched == true {
		rd := strings.Split(*raw, " ")
		p.paramPrefix = &rd[0]
		p.ProInfo.name = &rd[1]
		return nil
	}

	return errors.Errorf("%s parameter parsing failed: %s", utils.ProcessType, *raw)
}

func (p *Process) add(raw *string) error {
	return nil
}

type processSet struct {
	process *Process
}

var ProcessBus processSet

func (p *processSet) Init() {
	p.process = new(Process)
	p.process.ProInfo = new(ProcessInfo)
}

func (p *processSet) Add(raw *string) error {
	return nil
}

func (p *processSet) ListParamText() string {
	return p.process.put()
}

func (p *processSet) GetParam() interface{} {
	return p.process
}

func (p *processSet) Registry() map[string]Parameter {
	p.Init()
	return map[string]Parameter{utils.ProcessType: p.process}
}
