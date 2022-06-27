package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"strings"
)

type Options struct {
	Opts map[string]bool
}

type DBOptions struct {
	SupportParams map[string]map[string]string
	ParamPrefix   *string
	OptionsSet    *Options
}

func (d *DBOptions) Put() {
	fmt.Println("discardfile Info:", *d.ParamPrefix, *d.OptionsSet)
}

func (d *DBOptions) Init() {
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
	d.OptionsSet = &Options{
		Opts: map[string]bool{
			utils.SuppressionTrigger: false,
			utils.IgnoreReplicates:   false,
			utils.GetReplicates:      false,
			utils.IgnoreForeignkey:   false,
		},
	}
}

func (d *DBOptions) IsType(raw *string, dbType *string, processType *string) error {
	d.Init()
	_, ok := d.SupportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

func (d *DBOptions) Parse(raw *string) error {
	options := utils.TrimKeySpace(strings.Split(*raw, " "))
	optionsLength := len(options) - 1
	for i := 0; i < len(options); i++ {
		if strings.EqualFold(options[i], utils.DBOptionsType) {
			d.ParamPrefix = &options[i]
			if i+1 > optionsLength {
				return errors.Errorf("%s value must be specified", options[i])
			}
		}
		_, ok := d.OptionsSet.Opts[options[i]]
		if ok {
			d.OptionsSet.Opts[strings.ToUpper(options[i])] = true
		}
	}

	return nil
}

type DBOptionsSet struct{}

var DBOptionsBus DBOptionsSet

func (d *DBOptionsSet) Init() {

}

func (d *DBOptionsSet) Registry() map[string]Parameter {
	return map[string]Parameter{utils.DBOptionsType: &DBOptions{}}
}
