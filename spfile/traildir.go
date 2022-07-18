package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/utils"
	"strconv"
	"strings"
)

type TrailAttribute struct {
	dir       *string
	sizeKey   *string
	sizeValue *int
	sizeUnit  *string
	keepKey   *string
	keepValue *int
	keepUnit  *string
}

func (t *TrailAttribute) setDir(s *string) { t.dir = s }
func (t *TrailAttribute) GetDir() *string  { return t.dir }

func (t *TrailAttribute) setSizeKey(s *string) { t.sizeKey = s }
func (t *TrailAttribute) GetSizeKey() *string  { return t.sizeKey }

func (t *TrailAttribute) setSizeValue(s *int) { t.sizeValue = s }
func (t *TrailAttribute) GetSizeValue() *int  { return t.sizeValue }

func (t *TrailAttribute) setSizeUnit(s *string) { t.sizeUnit = s }
func (t *TrailAttribute) GetSizeUnit() *string  { return t.sizeUnit }

func (t *TrailAttribute) setKeepKey(s *string) { t.keepKey = s }
func (t *TrailAttribute) GetKeepKey() *string  { return t.keepKey }

func (t *TrailAttribute) setKeepVal(s *int) { t.keepValue = s }
func (t *TrailAttribute) GetKeepVal() *int  { return t.keepValue }

func (t *TrailAttribute) setKeepUnit(s *string) { t.keepUnit = s }
func (t *TrailAttribute) GetKeepUnit() *string  { return t.keepUnit }

type TrailDir struct {
	supportParams     map[string]map[string]string
	paramPrefix       *string
	DirTrailAttribute *TrailAttribute
}

func (t *TrailDir) put() string {
	return fmt.Sprintf("%s %s %s %d %s %d %s", *t.paramPrefix, *t.DirTrailAttribute.GetDir(), *t.DirTrailAttribute.GetSizeKey(), *t.DirTrailAttribute.GetSizeValue(), *t.DirTrailAttribute.GetKeepKey(), *t.DirTrailAttribute.GetKeepVal(), *t.DirTrailAttribute.GetKeepUnit())
}

// 初始化参数可以支持的数据库和进程
func (t *TrailDir) init() {
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
}

func (t *TrailDir) initDefault() error {
	return nil
}

func (t *TrailDir) isType(raw *string, dbType *string, processType *string) error {
	t.init()
	_, ok := t.supportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

func (t *TrailDir) parse(raw *string) error {
	trail := utils.TrimKeySpace(strings.Split(*raw, " "))
	trailLength := len(trail) - 1
	for i := 0; i < len(trail); i++ {
		switch {
		case strings.EqualFold(trail[i], utils.TrailDirType):
			t.paramPrefix = &trail[i]
			if i+1 > trailLength {
				return errors.Errorf("%s value must be specified", trail[i])
			}
			NextVal := &trail[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if t.DirTrailAttribute.dir != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			t.DirTrailAttribute.dir = NextVal
			i += 1

		case strings.EqualFold(trail[i], utils.TrailSizeKey):

			t.DirTrailAttribute.setSizeKey(&trail[i])
			if i+1 > trailLength {
				return errors.Errorf("%s value must be specified", utils.TrailSizeKey)
			}
			NextVal := &trail[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if t.DirTrailAttribute.sizeValue != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s, err := strconv.Atoi(*NextVal)
			if err != nil {
				return errors.Errorf("%s value is not a numeric integer: %s", trail[i], *NextVal)
			}

			if s < utils.DefaultTrailMinSize {
				return errors.Errorf("%s value %s cannot be less than the minimum size %d", trail[i], *NextVal, utils.DefaultTrailMinSize)
			}

			t.DirTrailAttribute.setSizeValue(&s)
			i += 1
		case strings.EqualFold(trail[i], utils.TrailKeepKey):

			t.DirTrailAttribute.setKeepKey(&trail[i])
			if i+1 > trailLength {
				return errors.Errorf("%s value must be specified", trail[i])
			}
			NextVal := &trail[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if t.DirTrailAttribute.keepValue != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s, err := strconv.Atoi(*NextVal)
			if err != nil {
				return errors.Errorf("%s value is not a numeric integer: %s", trail[i], *NextVal)
			}
			t.DirTrailAttribute.setKeepVal(&s)
			i += 1

			if i+1 > trailLength {
				return errors.Errorf("%s unit value must be specified", *t.DirTrailAttribute.GetKeepKey())
			}
			NextVal = &trail[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if t.DirTrailAttribute.GetKeepUnit() != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}

			if strings.EqualFold(*NextVal, utils.MB) || strings.EqualFold(*NextVal, utils.GB) || strings.EqualFold(*NextVal, utils.DAY) {
				t.DirTrailAttribute.setKeepUnit(NextVal)
				i += 1
			} else {
				return errors.Errorf("%s unit value Only supported %s/%s/%s", *t.DirTrailAttribute.GetKeepKey(), utils.MB, utils.GB, utils.DAY)
			}
		default:
			return errors.Errorf("unknown parameter: %s", trail[i])
		}
	}

	if t.DirTrailAttribute == nil {
		t.DirTrailAttribute = &TrailAttribute{
			sizeKey:   &utils.TrailSizeKey,
			sizeValue: &utils.DefaultTrailMaxSize,
			sizeUnit:  &utils.MB,
			keepKey:   &utils.TrailKeepKey,
			keepValue: &utils.DefaultTrailKeepValue,
			keepUnit:  &utils.DAY,
		}
	} else {
		if t.DirTrailAttribute.GetSizeValue() == nil {
			t.DirTrailAttribute.setSizeKey(&utils.TrailSizeKey)
			t.DirTrailAttribute.setSizeValue(&utils.DefaultTrailMaxSize)
			t.DirTrailAttribute.setSizeUnit(&utils.MB)
		}
		if t.DirTrailAttribute.GetKeepVal() == nil {
			t.DirTrailAttribute.setKeepKey(&utils.TrailKeepKey)
			t.DirTrailAttribute.setKeepVal(&utils.DefaultTrailKeepValue)
			t.DirTrailAttribute.setKeepUnit(&utils.DAY)

		}
	}

	return nil
}

func (t *TrailDir) add(raw *string) error {

	return nil
}

type trailDirSet struct {
	trailDir *TrailDir
}

var trailDirBus trailDirSet

func (t *trailDirSet) Init() {
	t.trailDir = new(TrailDir)
	t.trailDir.DirTrailAttribute = new(TrailAttribute)
}

func (t *trailDirSet) Add(raw *string) error {
	return nil
}

func (t *trailDirSet) ListParamText() string {
	return t.trailDir.put()
}

func (t *trailDirSet) GetParam() interface{} {
	return t.trailDir
}

func (t *trailDirSet) Registry() map[string]Parameter {
	t.Init()
	return map[string]Parameter{utils.TrailDirType: t.trailDir}
}
