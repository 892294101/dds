package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"regexp"
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
	for i, index := range t.TableListIndex {
		_, ok := t.TableList[index]
		if ok {
			if i > 0 {
				msg += fmt.Sprintf("\n")
			}
			msg += fmt.Sprintf("%s %s.%s", *t.ParamPrefix, index.OwnerValue, index.TableValue)
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
	reg, err := regexp.Compile(utils.TableRegular)
	if reg == nil || err != nil {
		return errors.Errorf("%s parameter Regular compilation error: %s", utils.TableType, *raw)
	}

	result := reg.FindStringSubmatch(*raw)
	if len(result) < 1 {
		return errors.Errorf("%s parameter Regular get substring error: %s", utils.TableType, *raw)
	}
	result = utils.TrimKeySpace(result)

	if t.ParamPrefix == nil {
		t.ParamPrefix = &result[1]
	}

	ownerTable := OwnerTable{result[3], result[5]}
	_, ok := t.TableList[ownerTable]
	if !ok {
		t.TableList[ownerTable] = nil
		t.TableListIndex = append(t.TableListIndex, ownerTable)
	}

	return nil
	/*matched, _ := regexp.MatchString(utils.TableRegular, *raw)
	if matched == true {
		rawText := *raw
		rawText = rawText[:len(rawText)-1]

		tab := utils.TrimKeySpace(strings.Split(rawText, " "))
		for i := 0; i < len(tab); i++ {
			if strings.EqualFold(tab[i], utils.TableType) {
				t.ParamPrefix = &tab[i]
			} else {
				tabVal := strings.Split(tab[i], ".")
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

	return errors.Errorf("Incorrect %s parameter user(or db) and table name rules: %s", utils.TableType, *raw)*/
}

// 当出现第二次参数进入, 需要进入add动作
func (t *TableSets) Add(raw *string) error {
	reg, err := regexp.Compile(utils.TableRegular)
	if reg == nil || err != nil {
		return errors.Errorf("%s parameter Regular compilation error: %s", utils.TableType, *raw)
	}

	result := reg.FindStringSubmatch(*raw)
	if len(result) < 1 {
		return errors.Errorf("%s parameter Regular get substring error: %s", utils.TableType, *raw)
	}
	result = utils.TrimKeySpace(result)

	ownerTable := OwnerTable{result[3], result[5]}
	_, ok := t.TableList[ownerTable]
	if !ok {
		t.TableList[ownerTable] = nil
		t.TableListIndex = append(t.TableListIndex, ownerTable)
	}

	/*matched, _ := regexp.MatchString(utils.TableRegular, *raw)
	if matched == true {
		rawText := *raw
		rawText = rawText[:len(rawText)-1]

		tab := utils.TrimKeySpace(strings.Split(rawText, " "))
		for i := 0; i < len(tab); i++ {
			if strings.EqualFold(tab[i], utils.TableType) {
				t.ParamPrefix = &tab[i]
			} else {
				tabVal := strings.Split(tab[i], ".")
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
	return t.table.Add(raw)
}

func (t *TableSet) ListParamText() string {
	return t.table.Put()
}

func (t *TableSet) Registry() map[string]Parameter {
	t.Init()
	return map[string]Parameter{utils.TableType: t.table}
}
