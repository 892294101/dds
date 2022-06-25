package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type Process struct {
	ParamPrefix string
	Name        string
}

func (p *Process) Put() {
	fmt.Println("process Info: ", p.ParamPrefix, p.Name)
}

func (p *Process) Parse(raw string) error {
	matched, _ := regexp.MatchString(ProcessRegular, raw)
	if matched == true {
		rd := strings.Split(raw, " ")
		p.ParamPrefix = rd[0]
		p.Name = rd[1]
		return nil
	}

	return errors.Errorf("process parameter parsing failed: %s", raw)
}

type processSet struct{}

func (p *processSet) Init() () {

}

func (p *processSet) Registry() map[string]Parameter {
	return map[string]Parameter{ProcessType: &Process{}}
}
