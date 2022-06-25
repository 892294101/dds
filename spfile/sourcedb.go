package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
)

type sourceDB struct {
	ParamPrefix string
	Address     string
	Port        string
	Database    string
	Type        string
	UserId      string
	PassWord    string
}

func (s *sourceDB) Put() {
	fmt.Println("sourceDB Info: ", s)
}

func (s *sourceDB) Parse(raw string) error {
	matched, _ := regexp.MatchString(SourceDBRegular, raw)
	if matched == true {

		return nil
	}

	return errors.Errorf("source db parameter parsing failed: %s", raw)
}

type sourceDBSet struct{}

func (sd *sourceDBSet) Init() () {

}

func (sd *sourceDBSet) Registry() map[string]Parameter {
	return map[string]Parameter{SourceDBType: &sourceDB{}}
}
