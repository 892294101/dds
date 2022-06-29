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

func (o *Options) SetOption(s *string) error {
	ops := strings.ToUpper(*s)
	switch ops {
	case utils.GetReplicates:
		o.Opts[utils.GetReplicates] = true
	case utils.IgnoreReplicates:
		o.Opts[utils.GetReplicates] = false
	default:
		ops := strings.ToUpper(*s)
		_, ok := o.Opts[ops]
		if ok {
			o.Opts[ops] = true
		} else {
			return errors.Errorf("unknown parameter: %s", *s)
		}
	}
	return nil
}

func (o *Options) GetOption(s *string) {

}

type DBOptions struct {
	SupportParams map[string]map[string]string
	ParamPrefix   *string
	OptionsSet    *Options
}

func (d *DBOptions) Put() string {
	var msg string
	msg += fmt.Sprintf("%s", *d.ParamPrefix)
	for s, b := range d.OptionsSet.Opts {
		msg += fmt.Sprintf(" %s %v", s, b)
	}
	return msg
}

// 当传入参数时, 初始化特定参数的值
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
			utils.IgnoreForeignkey:   false,
			utils.GetReplicates:      false,
		},
	}
}

// 当没有参数时, 初始化此参数默认值
func (d *DBOptions) InitDefault() error {
	d.Init()
	d.ParamPrefix = &utils.DBOptionsType
	return nil
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
		} else {
			err := d.OptionsSet.SetOption(&options[i])
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (d *DBOptions) Add(raw *string) error {
	return nil
}

type DBOptionsSet struct {
	dbOps *DBOptions
}

var DBOptionsBus DBOptionsSet

func (d *DBOptionsSet) Init() {
	d.dbOps = new(DBOptions)
}

func (d *DBOptionsSet) Add(raw *string) error {
	return nil
}

func (d *DBOptionsSet) ListParamText() string {
	return d.dbOps.Put()
}

func (d *DBOptionsSet) Registry() map[string]Parameter {
	d.Init()
	return map[string]Parameter{utils.DBOptionsType: d.dbOps}
}
