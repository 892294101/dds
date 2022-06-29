package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"regexp"
	"strings"
)

type OwnerTable struct {
	OwnerValue string
	TableValue string
}

type ETL struct {
	AddColumn    string
	DeleteColumn string
	UpdateColumn string
	MapColumn    string
}

type TableSets struct {
	SupportParams  map[string]map[string]string
	ParamPrefix    *string
	TableList      map[OwnerTable]*ETL
	TableListIndex []OwnerTable
}

func (t *TableSets) Put() string {
	var msg string
	for _, index := range t.TableListIndex {

		_, ok := t.TableList[index]
		if ok {
			msg += fmt.Sprintf("%s %s.%s\n", *t.ParamPrefix, index.OwnerValue, index.TableValue)
		}

	}
	return msg
}

// 当传入参数时, 初始化特定参数的值
func (t *TableSets) Init() {
	t.SupportParams = map[string]map[string]string{
		utils.MySQL: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
		utils.MariaDB: {
			utils.Extract:  utils.Extract,
			utils.Replicat: utils.Replicat,
		},
	}
	t.TableList = make(map[OwnerTable]*ETL)

}

// 当没有参数时, 初始化此参数默认值
func (t *TableSets) InitDefault() error {
	t.Init()
	t.ParamPrefix = &utils.DBOptionsType
	return nil
}

func (t *TableSets) IsType(raw *string, dbType *string, processType *string) error {
	t.Init()
	_, ok := t.SupportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

// 新参数进入后, 第一次需要进入解析动作
func (t *TableSets) Parse(raw *string) error {
	matched, _ := regexp.MatchString(utils.TableRegular, *raw)
	if matched == true {
		tab := utils.TrimKeySpace(strings.Split(*raw, " "))
		for i := 0; i < len(tab); i++ {
			if strings.EqualFold(tab[i], utils.TableType) {
				t.ParamPrefix = &tab[i]
			} else {
				tabVal := strings.Split(tab[i][:len(tab[i])-1], ".")
				ownerTable := OwnerTable{tabVal[0], tabVal[1]}
				_, ok := t.TableList[ownerTable]
				if !ok {
					t.TableList[ownerTable] = nil
					t.TableListIndex = append(t.TableListIndex, ownerTable)
				}

			}
		}
		return nil
	}

	if ok := strings.HasSuffix(*raw, ";"); !ok {
		return errors.Errorf("%s parameter must end with a semicolon: %s", utils.TableType, *raw)
	}

	return errors.Errorf("Incorrect %s parameter user(or db) and table name rules: %s", utils.TableType, *raw)

	/*if ok := strings.HasSuffix(*raw, utils.TableRawDataSuffix); !ok {
		return errors.Errorf("%s parameter must end with a semicolon: %s", utils.TableType, *raw)
	}

	tab := strings.Split(*raw, " ")
	if len(tab) != 2 {
		return errors.Errorf("%s Parameter length mismatch: %s", utils.TableType, *raw)
	}

	tabLength := len(tab) - 1
	for i := 0; i < len(tab); i++ {
		if strings.EqualFold(tab[i], utils.TableType) {
			t.ParamPrefix = &tab[i]
			if i+1 > tabLength {
				return errors.Errorf("%s value must be specified", tab[i])
			}
		} else {
			tabVal := strings.Split(tab[i][:len(tab[i])-1], ".")
			ownerTable := OwnerTable{tabVal[0], tabVal[1]}
			_, ok := t.TableList[ownerTable]
			if !ok {
				t.TableList[ownerTable] = nil
				t.TableListIndex = append(t.TableListIndex, ownerTable)
			}

		}
	}*/

}

// 当出现第二次参数进入, 需要进入add动作
func (t *TableSets) Add(raw *string) error {
	matched, _ := regexp.MatchString(utils.TableRegular, *raw)
	if matched == true {
		tab := utils.TrimKeySpace(strings.Split(*raw, " "))
		for i := 0; i < len(tab); i++ {
			if strings.EqualFold(tab[i], utils.TableType) {
				t.ParamPrefix = &tab[i]
			} else {
				tabVal := strings.Split(tab[i][:len(tab[i])-1], ".")
				ownerTable := OwnerTable{tabVal[0], tabVal[1]}
				_, ok := t.TableList[ownerTable]
				if !ok {
					t.TableList[ownerTable] = nil
					t.TableListIndex = append(t.TableListIndex, ownerTable)
				}

			}
		}
		return nil
	}
	if ok := strings.HasSuffix(*raw, ";"); !ok {
		return errors.Errorf("%s parameter must end with a semicolon: %s", utils.TableType, *raw)
	}
	return errors.Errorf("Incorrect %s parameter user(or db) and table name rules: %s", utils.TableType, *raw)

	/*tab := strings.Split(*raw, " ")
	tabLength := len(tab) - 1
	if len(tab) != 2 {
		return errors.Errorf("%s Parameter length mismatch: %s", utils.TableType, *raw)
	}
	for i := 0; i < len(tab); i++ {
		if strings.EqualFold(tab[i], utils.TableType) {
			if i+1 > tabLength {
				return errors.Errorf("%s value must be specified", tab[i])
			}
		} else {
			tabVal := strings.Split(tab[i][:len(tab[i])-1], ".")
			ownerTable := OwnerTable{tabVal[0], tabVal[1]}
			_, ok := t.TableList[ownerTable]
			if !ok {
				t.TableList[ownerTable] = nil
				t.TableListIndex = append(t.TableListIndex, ownerTable)
			}

		}
	}*/

}

type TableSet struct {
	table *TableSets
}

var TableSetBus TableSet

func (t *TableSet) Init() {
	t.table = new(TableSets)
}

func (t *TableSet) Add(raw *string) error {
	return t.table.Add(raw)
}

func (t *TableSet) ListParamText() string {
	return t.table.Put()
}

func (t *TableSet) Registry() map[string]Parameter {
	t.Init()
	return map[string]Parameter{utils.TableType: t.table}
}
