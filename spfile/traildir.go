package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"strconv"
	"strings"
)

type TrailAttribute struct {
	SizeKey   *string
	SizeValue *int
	SizeUnit  *string
	KeepKey   *string
	KeepValue *int
	KeepUnit  *string
}

func (t *TrailAttribute) SetSizeKey(s *string) { t.SizeKey = s }
func (t *TrailAttribute) GetSizeKey() *string  { return t.SizeKey }

func (t *TrailAttribute) SetSizeValue(s *int) { t.SizeValue = s }
func (t *TrailAttribute) GetSizeValue() *int  { return t.SizeValue }

func (t *TrailAttribute) SetSizeUnit(s *string) { t.SizeUnit = s }
func (t *TrailAttribute) GetSizeUnit() *string  { return t.SizeUnit }

func (t *TrailAttribute) SetKeepKey(s *string) { t.KeepKey = s }
func (t *TrailAttribute) GetKeepKey() *string  { return t.KeepKey }

func (t *TrailAttribute) SetKeepVal(s *int) { t.KeepValue = s }
func (t *TrailAttribute) GetKeepVal() *int  { return t.KeepValue }

func (t *TrailAttribute) SetKeepUnit(s *string) { t.KeepUnit = s }
func (t *TrailAttribute) GetKeepUnit() *string  { return t.KeepUnit }

type TrailDir struct {
	SupportParams  map[string]map[string]string
	ParamPrefix    *string
	Dir            *string
	TrailAttribute *TrailAttribute
}

func (t *TrailDir) Put() {
	fmt.Println("traildir Info:", *t.ParamPrefix, *t.Dir, *t.TrailAttribute.GetSizeKey(), *t.TrailAttribute.GetSizeValue(), *t.TrailAttribute.GetKeepKey(), *t.TrailAttribute.GetKeepVal(), *t.TrailAttribute.GetKeepUnit())
}

// 初始化参数可以支持的数据库和进程
func (t *TrailDir) Init() {
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
}

func (t *TrailDir) IsType(raw *string, dbType *string, processType *string) error {
	t.Init()
	_, ok := t.SupportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

func (t *TrailDir) Parse(raw *string) error {
	trail := utils.TrimKeySpace(strings.Split(*raw, " "))
	trailLength := len(trail) - 1
	for i := 0; i < len(trail); i++ {
		switch {
		case strings.EqualFold(trail[i], utils.TrailDirType):
			t.ParamPrefix = &trail[i]
			if i+1 > trailLength {
				return errors.Errorf("%s value must be specified", trail[i])
			}
			NextVal := &trail[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if t.Dir != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			t.Dir = NextVal
			i += 1

		case strings.EqualFold(trail[i], utils.TrailSizeKey):
			if t.TrailAttribute == nil {
				t.TrailAttribute = new(TrailAttribute)
			}

			t.TrailAttribute.SetSizeKey(&trail[i])
			if i+1 > trailLength {
				return errors.Errorf("%s value must be specified", utils.TrailSizeKey)
			}
			NextVal := &trail[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if t.TrailAttribute.SizeValue != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s, err := strconv.Atoi(*NextVal)
			if err != nil {
				return errors.Errorf("%s value is not a numeric integer: %s", trail[i], *NextVal)
			}

			if s < utils.DefaultTrailMinSize {
				return errors.Errorf("%s value %s cannot be less than the minimum size %d", trail[i], *NextVal, utils.DefaultTrailMinSize)
			}

			t.TrailAttribute.SetSizeValue(&s)
			i += 1
		case strings.EqualFold(trail[i], utils.TrailKeepKey):
			if t.TrailAttribute == nil {
				t.TrailAttribute = new(TrailAttribute)
			}

			t.TrailAttribute.SetKeepKey(&trail[i])
			if i+1 > trailLength {
				return errors.Errorf("%s value must be specified", trail[i])
			}
			NextVal := &trail[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if t.TrailAttribute.KeepValue != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s, err := strconv.Atoi(*NextVal)
			if err != nil {
				return errors.Errorf("%s value is not a numeric integer: %s", trail[i], *NextVal)
			}
			t.TrailAttribute.SetKeepVal(&s)
			i += 1

			if i+1 > trailLength {
				return errors.Errorf("%s unit value must be specified", *t.TrailAttribute.GetKeepKey())
			}
			NextVal = &trail[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if t.TrailAttribute.GetKeepUnit() != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}

			if strings.EqualFold(*NextVal, utils.MB) || strings.EqualFold(*NextVal, utils.GB) || strings.EqualFold(*NextVal, utils.DAY) {
				t.TrailAttribute.SetKeepUnit(NextVal)
				i += 1
			} else {
				return errors.Errorf("%s unit value Only supported %s/%s/%s", *t.TrailAttribute.GetKeepKey(), utils.MB, utils.GB, utils.DAY)
			}
		default:
			return errors.Errorf("unknown parameter: %s", trail[i])
		}
	}

	if t.TrailAttribute == nil {
		t.TrailAttribute = &TrailAttribute{
			SizeKey:   &utils.TrailSizeKey,
			SizeValue: &utils.DefaultTrailMaxSize,
			SizeUnit:  &utils.MB,
			KeepKey:   &utils.TrailKeepKey,
			KeepValue: &utils.DefaultTrailKeepValue,
			KeepUnit:  &utils.DAY,
		}
	} else {
		if t.TrailAttribute.GetSizeValue() == nil {
			t.TrailAttribute.SetSizeKey(&utils.TrailSizeKey)
			t.TrailAttribute.SetSizeValue(&utils.DefaultTrailMaxSize)
			t.TrailAttribute.SetSizeUnit(&utils.MB)
		}
		if t.TrailAttribute.GetKeepVal() == nil {
			t.TrailAttribute.SetKeepKey(&utils.TrailKeepKey)
			t.TrailAttribute.SetKeepVal(&utils.DefaultTrailKeepValue)
			t.TrailAttribute.SetKeepUnit(&utils.DAY)

		}
	}

	return nil
}

type trailDirSet struct{}

var trailDirBus trailDirSet

func (t *trailDirSet) Init() {

}

func (t *trailDirSet) Registry() map[string]Parameter {
	return map[string]Parameter{utils.TrailDirType: &TrailDir{}}
}
