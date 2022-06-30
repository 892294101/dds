package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"regexp"
)

type ExcludeOwnerTable struct {
	OwnerValue string
	TableValue string
}

type ExcludeTableSets struct {
	SupportParams  map[string]map[string]string
	ParamPrefix    *string
	TableList      map[ExcludeOwnerTable]*string
	TableListIndex []ExcludeOwnerTable
}

func (e *ExcludeTableSets) Put() string {
	var msg string
	for _, index := range e.TableListIndex {
		_, ok := e.TableList[index]
		if ok {
			msg += fmt.Sprintf("%s %s.%s\n", *e.ParamPrefix, index.OwnerValue, index.TableValue)
		}

	}
	return msg
}

// 当传入参数时, 初始化特定参数的值
func (e *ExcludeTableSets) Init() {
	e.SupportParams = map[string]map[string]string{
		utils.MySQL: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
		utils.MariaDB: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
	}
	e.TableList = make(map[ExcludeOwnerTable]*string)

}

// 当没有参数时, 初始化此参数默认值
func (e *ExcludeTableSets) InitDefault() error {
	e.Init()
	e.ParamPrefix = &utils.DBOptionsType
	return nil
}

func (e *ExcludeTableSets) IsType(raw *string, dbType *string, processType *string) error {
	e.Init()
	_, ok := e.SupportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

// 新参数进入后, 第一次需要进入解析动作
func (e *ExcludeTableSets) Parse(raw *string) error {
	reg, err := regexp.Compile(utils.TableExcludeRegular)
	if reg == nil || err != nil {
		return errors.Errorf("%s parameter Regular compilation error: %s", utils.TableExcludeType, *raw)
	}

	result := reg.FindStringSubmatch(*raw)
	if len(result) < 1 {
		return errors.Errorf("%s parameter Regular get substring error: %s", utils.TableExcludeType, *raw)
	}
	result = utils.TrimKeySpace(result)

	if e.ParamPrefix == nil {
		e.ParamPrefix = &result[1]
	}

	ownerTable := ExcludeOwnerTable{result[3], result[5]}
	_, ok := e.TableList[ownerTable]
	if !ok {
		e.TableList[ownerTable] = nil
		e.TableListIndex = append(e.TableListIndex, ownerTable)
	}
	return nil
}

// 当出现第二次参数进入, 需要进入add动作
func (e *ExcludeTableSets) Add(raw *string) error {
	reg, err := regexp.Compile(utils.TableExcludeRegular)
	if reg == nil || err != nil {
		return errors.Errorf("%s parameter Regular compilation error: %s", utils.TableExcludeType, *raw)
	}

	result := reg.FindStringSubmatch(*raw)
	if len(result) < 1 {
		return errors.Errorf("%s parameter Regular get substring error: %s", utils.TableExcludeType, *raw)
	}
	result = utils.TrimKeySpace(result)

	ownerTable := ExcludeOwnerTable{result[3], result[5]}
	_, ok := e.TableList[ownerTable]
	if !ok {
		e.TableList[ownerTable] = nil
		e.TableListIndex = append(e.TableListIndex, ownerTable)
	}

	return nil
}

type ExcludeTableSet struct {
	table *ExcludeTableSets
}

var ExcludeTableSetBus ExcludeTableSet

func (e *ExcludeTableSet) Init() {
	e.table = new(ExcludeTableSets)
}

func (e *ExcludeTableSet) Add(raw *string) error {
	return e.table.Add(raw)
}

func (e *ExcludeTableSet) ListParamText() string {
	return e.table.Put()
}

func (e *ExcludeTableSet) Registry() map[string]Parameter {
	e.Init()
	return map[string]Parameter{utils.TableExcludeType: e.table}
}
