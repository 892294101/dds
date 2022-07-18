package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/892294101/dds/utils"
	"regexp"
)

type ownerTable struct {
	ownerValue string
	tableValue string
}

type ETL struct {
	addColumn    string
	deleteColumn string
	updateColumn string
	mapColumn    string
}

type TableSets struct {
	supportParams  map[string]map[string]string
	paramPrefix    *string
	TableList      map[ownerTable]*ETL
	tableListIndex []ownerTable
}

func (t *TableSets) put() string {
	var msg string
	for i, index := range t.tableListIndex {
		_, ok := t.TableList[index]
		if ok {
			if i > 0 {
				msg += fmt.Sprintf("\n")
			}
			msg += fmt.Sprintf("%s %s.%s", *t.paramPrefix, index.ownerValue, index.tableValue)
		}

	}
	return msg
}

// 当传入参数时, 初始化特定参数的值
func (t *TableSets) init() {
	t.supportParams = map[string]map[string]string{
		utils.MySQL: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
		utils.MariaDB: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
	}
	t.TableList = make(map[ownerTable]*ETL)

}

// 当没有参数时, 初始化此参数默认值
func (t *TableSets) initDefault() error {
	t.init()
	t.paramPrefix = &utils.DBOptionsType
	return nil
}

func (t *TableSets) isType(raw *string, dbType *string, processType *string) error {
	t.init()
	_, ok := t.supportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

// 新参数进入后, 第一次需要进入解析动作
func (t *TableSets) parse(raw *string) error {
	reg, err := regexp.Compile(utils.TableRegular)
	if reg == nil || err != nil {
		return errors.Errorf("%s parameter Regular compilation error: %s", utils.TableType, *raw)
	}

	result := reg.FindStringSubmatch(*raw)
	if len(result) < 1 {
		return errors.Errorf("%s parameter Regular get substring error: %s", utils.TableType, *raw)
	}
	result = utils.TrimKeySpace(result)

	if t.paramPrefix == nil {
		t.paramPrefix = &result[1]
	}

	ownerTable := ownerTable{result[3], result[5]}
	_, ok := t.TableList[ownerTable]
	if !ok {
		t.TableList[ownerTable] = nil
		t.tableListIndex = append(t.tableListIndex, ownerTable)
	}

	return nil
	/*matched, _ := regexp.MatchString(utils.TableRegular, *raw)
	if matched == true {
		rawText := *raw
		rawText = rawText[:len(rawText)-1]

		tab := utils.TrimKeySpace(strings.Split(rawText, " "))
		for i := 0; i < len(tab); i++ {
			if strings.EqualFold(tab[i], utils.TableType) {
				t.paramPrefix = &tab[i]
			} else {
				tabVal := strings.Split(tab[i], ".")
				ownerTable := ownerTable{tabVal[0], tabVal[1]}
				_, ok := t.TableList[ownerTable]
				if !ok {
					t.TableList[ownerTable] = nil
					t.tableListIndex = append(t.tableListIndex, ownerTable)
				}

			}
		}
		return nil
	}

	if ok := strings.HasSuffix(*raw, ";"); !ok {
		return errors.Errorf("%s parameter must end with a semicolon: %s", utils.TableType, *raw)
	}

	return errors.Errorf("Incorrect %s parameter user(or db) and table name rules: %s", utils.TableType, *raw)*/
}

// 当出现第二次参数进入, 需要进入add动作
func (t *TableSets) add(raw *string) error {
	reg, err := regexp.Compile(utils.TableRegular)
	if reg == nil || err != nil {
		return errors.Errorf("%s parameter Regular compilation error: %s", utils.TableType, *raw)
	}

	result := reg.FindStringSubmatch(*raw)
	if len(result) < 1 {
		return errors.Errorf("%s parameter Regular get substring error: %s", utils.TableType, *raw)
	}
	result = utils.TrimKeySpace(result)

	ownerTable := ownerTable{result[3], result[5]}
	_, ok := t.TableList[ownerTable]
	if !ok {
		t.TableList[ownerTable] = nil
		t.tableListIndex = append(t.tableListIndex, ownerTable)
	}

	/*matched, _ := regexp.MatchString(utils.TableRegular, *raw)
	if matched == true {
		rawText := *raw
		rawText = rawText[:len(rawText)-1]

		tab := utils.TrimKeySpace(strings.Split(rawText, " "))
		for i := 0; i < len(tab); i++ {
			if strings.EqualFold(tab[i], utils.TableType) {
				t.paramPrefix = &tab[i]
			} else {
				tabVal := strings.Split(tab[i], ".")
				ownerTable := ownerTable{tabVal[0], tabVal[1]}
				_, ok := t.TableList[ownerTable]
				if !ok {
					t.TableList[ownerTable] = nil
					t.tableListIndex = append(t.tableListIndex, ownerTable)
				}

			}
		}
		return nil
	}
	if ok := strings.HasSuffix(*raw, ";"); !ok {
		return errors.Errorf("%s parameter must end with a semicolon: %s", utils.TableType, *raw)
	}
	return errors.Errorf("Incorrect %s parameter user(or db) and table name rules: %s", utils.TableType, *raw)
	*/
	return nil
}

type TableSet struct {
	table *TableSets
}

var TableSetBus TableSet

func (t *TableSet) Init() {
	t.table = new(TableSets)
}

func (t *TableSet) Add(raw *string) error {
	return t.table.add(raw)
}

func (t *TableSet) ListParamText() string {
	return t.table.put()
}

func (t *TableSet) GetParam() interface{} {
	return t.table
}

func (t *TableSet) Registry() map[string]Parameter {
	t.Init()
	return map[string]Parameter{utils.TableType: t.table}
}
