package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"strings"
)

type DiscardFile struct {
	SupportParams map[string]map[string]string
	ParamPrefix   *string
	Dir           *string
}

func (d *DiscardFile) Put() {
	fmt.Println("discardfile Info:", *d.ParamPrefix, *d.Dir)
}

func (d *DiscardFile) Init() {
	d.SupportParams = map[string]map[string]string{
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

func (d *DiscardFile) IsType(raw *string, dbType *string, processType *string) error {
	d.Init()
	_, ok := d.SupportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

func (d *DiscardFile) Parse(raw *string) error {
	discards := utils.TrimKeySpace(strings.Split(*raw, " "))
	discardLength := len(discards) - 1
	for i := 0; i < len(discards); i++ {
		switch {
		case strings.EqualFold(discards[i], utils.DiscardFileType):
			d.ParamPrefix = &discards[i]
			if i+1 > discardLength {
				return errors.Errorf("%s value must be specified", discards[i])
			}
			NextVal := &discards[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if d.Dir != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			d.Dir = NextVal
			i += 1
		default:
			return errors.Errorf("unknown parameter: %s", discards[i])
		}
	}

	return nil
}

type DiscardFileSet struct{}

var DiscardFileBus DiscardFileSet

func (t *DiscardFileSet) Init() {

}

func (t *DiscardFileSet) Registry() map[string]Parameter {
	return map[string]Parameter{utils.DiscardFileType: &DiscardFile{}}
}
