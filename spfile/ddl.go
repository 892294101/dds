package spfile

import (
	"myGithubLib/dds/extract/mysql/utils"
)

type DDL struct {
	SupportParams map[string]map[string]string
	ParamPrefix   *string
}

func (d *DDL) Init() {}

func (d *DDL) InitDefault() error {
	d.Init()
	d.ParamPrefix = &utils.DBOptionsType
	return nil
}

func (d *DDL) IsType(raw *string, dbType *string, processType *string) error {
	d.Init()

	return nil
}

func (d *DDL) Parse(raw *string) error {

	return nil
}
