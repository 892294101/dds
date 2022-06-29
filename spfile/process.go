package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"regexp"
	"strings"
)

type Process struct {
	SupportParams map[string]map[string]string
	ParamPrefix   string
	Name          string
}

func (p *Process) Put() string {
	return fmt.Sprintf("%s %s", p.ParamPrefix, p.Name)
}

// 初始化参数可以支持的数据库和进程
func (p *Process) Init() {
	p.SupportParams = map[string]map[string]string{
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

func (p *Process) InitDefault() error {
	return nil
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
	matched, _ := regexp.MatchString(utils.ProcessRegular, *raw)
	if matched == true {
		rd := strings.Split(*raw, " ")
		p.ParamPrefix = rd[0]
		p.Name = rd[1]
		return nil
	}

	return errors.Errorf("%s parameter parsing failed: %s", utils.ProcessType, *raw)
}

func (p *Process) Add(raw *string) error {
	return nil
}
type processSet struct {
	process *Process
}

var ProcessBus processSet

func (p *processSet) Init() {
	p.process = new(Process)
}

func (p *processSet) Add(raw *string) error {
	return nil
}

func (p *processSet) ListParamText() string {
	return p.process.Put()
}

func (p *processSet) Registry() map[string]Parameter {
	p.Init()
	return map[string]Parameter{utils.ProcessType: p.process}
}
