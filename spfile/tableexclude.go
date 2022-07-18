package spfile

import (
	"fmt"
	"github.com/892294101/dds/utils"
	"github.com/pkg/errors"
	"regexp"
)

type excludeOwnerTable struct {
	ownerValue string
	tableValue string
}

type ExcludeTableSets struct {
	supportParams  map[string]map[string]string
	paramPrefix    *string
	TableList      map[excludeOwnerTable]*string
	tableListIndex []excludeOwnerTable
}

func (e *ExcludeTableSets) put() string {
	var msg string
	for _, index := range e.tableListIndex {
		_, ok := e.TableList[index]
		if ok {
			msg += fmt.Sprintf("%s %s.%s\n", *e.paramPrefix, index.ownerValue, index.tableValue)
		}
	}
	return msg
}

// 当传入参数时, 初始化特定参数的值
func (e *ExcludeTableSets) init() {
	e.supportParams = map[string]map[string]string{
		utils.MySQL: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
		utils.MariaDB: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
	}
	e.TableList = make(map[excludeOwnerTable]*string)

}

// 当没有参数时, 初始化此参数默认值
func (e *ExcludeTableSets) initDefault() error {
	e.init()
	e.paramPrefix = &utils.DBOptionsType
	return nil
}

func (e *ExcludeTableSets) isType(raw *string, dbType *string, processType *string) error {
	e.init()
	_, ok := e.supportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

// 新参数进入后, 第一次需要进入解析动作
func (e *ExcludeTableSets) parse(raw *string) error {
	reg, err := regexp.Compile(utils.TableExcludeRegular)
	if reg == nil || err != nil {
		return errors.Errorf("%s parameter Regular compilation error: %s", utils.TableExcludeType, *raw)
	}

	result := reg.FindStringSubmatch(*raw)
	if len(result) < 1 {
		return errors.Errorf("%s parameter Regular get substring error: %s", utils.TableExcludeType, *raw)
	}
	result = utils.TrimKeySpace(result)

	if e.paramPrefix == nil {
		e.paramPrefix = &result[1]
	}

	ownerTable := excludeOwnerTable{result[3], result[5]}
	_, ok := e.TableList[ownerTable]
	if !ok {
		e.TableList[ownerTable] = nil
		e.tableListIndex = append(e.tableListIndex, ownerTable)
	}
	return nil
}

// 当出现第二次参数进入, 需要进入add动作
func (e *ExcludeTableSets) add(raw *string) error {
	reg, err := regexp.Compile(utils.TableExcludeRegular)
	if reg == nil || err != nil {
		return errors.Errorf("%s parameter Regular compilation error: %s", utils.TableExcludeType, *raw)
	}

	result := reg.FindStringSubmatch(*raw)
	if len(result) < 1 {
		return errors.Errorf("%s parameter Regular get substring error: %s", utils.TableExcludeType, *raw)
	}
	result = utils.TrimKeySpace(result)

	ownerTable := excludeOwnerTable{result[3], result[5]}
	_, ok := e.TableList[ownerTable]
	if !ok {
		e.TableList[ownerTable] = nil
		e.tableListIndex = append(e.tableListIndex, ownerTable)
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
	return e.table.add(raw)
}

func (e *ExcludeTableSet) ListParamText() string {
	return e.table.put()
}

func (e *ExcludeTableSet) GetParam() interface{} {
	return e.table
}

func (e *ExcludeTableSet) Registry() map[string]Parameter {
	e.Init()
	return map[string]Parameter{utils.TableExcludeType: e.table}
}
