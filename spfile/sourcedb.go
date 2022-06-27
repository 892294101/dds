package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"regexp"
	"strings"
)

type PortModel struct {
	Key   *string
	Value *string
}

type DatabaseModel struct {
	Key   *string
	Value *string
}

type TypeModel struct {
	Key   *string
	Value *string
}

type UserIdModel struct {
	Key   *string
	Value *string
}

type PassWordModel struct {
	Key   *string
	Value *string
}

type sourceDB struct {
	SupportParams map[string]map[string]string // 参数支持吃数据库和进程
	ParamPrefix   *string                      // 参数前缀
	Address       *string                      // 数据库地址
	Port          *PortModel                   // 数据库端口
	Database      *DatabaseModel               // 连接的数据库
	Type          *TypeModel                   // 连接数据库类型, mysql或 mariadb
	UserId        *UserIdModel                 // 用户名
	PassWord      *PassWordModel               // 密码
}

func (s *sourceDB) Put() {
	fmt.Println("sourceDB Info:", *s.ParamPrefix, *s.Address, *s.Port.Key, *s.Port.Value, *s.Database.Key, *s.Database.Value, *s.Type.Key, *s.Type.Value, *s.UserId.Key, *s.UserId.Value, *s.PassWord.Key, *s.PassWord.Value)
}

// 初始化参数可以支持的数据库和进程

func (s *sourceDB) Init() {
	s.SupportParams = map[string]map[string]string{
		utils.MySQL: {
			utils.Extract: utils.Extract,
		},
		utils.MariaDB: {
			utils.Extract: utils.Extract,
		},
	}
	/*s.Port = new(PortModel)
	s.Database = new(DatabaseModel)
	s.Type = new(TypeModel)
	s.UserId = new(UserIdModel)
	s.PassWord = new(PassWordModel)
	*/
}

func (s *sourceDB) IsType(raw *string, dbType *string, processType *string) error {
	s.Init()
	_, ok := s.SupportParams[*dbType][*processType]
	if ok {
		return nil
	}
	return errors.Errorf("The %s %s process does not support this parameter: %s", *dbType, *processType, *raw)
}

func (s *sourceDB) Parse(raw *string) error {
	sdb := utils.TrimKeySpace(strings.Split(*raw, " "))
	sdbLength := len(sdb) - 1

	for i := 0; i < len(sdb); i++ {
		switch {
		case strings.EqualFold(sdb[i], utils.SourceDBType):
			s.ParamPrefix = &sdb[i]
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", utils.SourceDBType)
			}
			NextVal := &sdb[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.Address != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}

			match, _ := regexp.MatchString(IpV4Reg, *NextVal)
			if !match {
				return errors.Errorf("%s is an illegal IPV4 address\n", *NextVal)
			}

			s.Address = NextVal
			i += 1
		case strings.EqualFold(sdb[i], utils.Port):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", utils.Port)
			}
			NextVal := &sdb[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.Port != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			match, _ := regexp.MatchString(IpV4Port, *NextVal)
			if !match {
				return errors.Errorf("%s is an illegal IPV4 Port\n", *NextVal)
			}

			s.Port = &PortModel{Key: &sdb[i], Value: NextVal}
			i += 1
		case strings.EqualFold(sdb[i], utils.DataBase):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", utils.DataBase)
			}
			NextVal := &sdb[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.Database != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s.Database = &DatabaseModel{Key: &sdb[i], Value: NextVal}
			i += 1
		case strings.EqualFold(sdb[i], utils.Types):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", utils.Types)
			}
			NextVal := &sdb[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.Type != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s.Type = &TypeModel{Key: &sdb[i], Value: NextVal}
			i += 1
		case strings.EqualFold(sdb[i], utils.UserId):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", utils.UserId)
			}
			NextVal := &sdb[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.UserId != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s.UserId = &UserIdModel{Key: &sdb[i], Value: NextVal}
			i += 1
		case strings.EqualFold(sdb[i], utils.PassWord):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", utils.PassWord)
			}
			NextVal := &sdb[i+1]
			if utils.KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.PassWord != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s.PassWord = &PassWordModel{Key: &sdb[i], Value: NextVal}
			i += 1
		default:
			return errors.Errorf("unknown parameter: %s", sdb[i])
		}

	}

	if s.Port == nil {
		s.Port = &PortModel{Key: &utils.Port, Value: &utils.DefaultPort}
	}
	if s.Database == nil {
		s.Database = &DatabaseModel{Key: &utils.DataBase, Value: &utils.DefaultDataBase}
	}
	if s.Type == nil {
		s.Type = &TypeModel{Key: &utils.Types, Value: &utils.DefaultTypes}
	}
	if s.UserId == nil {
		s.UserId = &UserIdModel{Key: &utils.UserId, Value: &utils.DefaultUserId}
	}
	if s.PassWord == nil {
		return errors.Errorf("%s Password must be specified", utils.SourceDBType)
	}

	return nil
}

type sourceDBSet struct{}

var sourceDBSetBus sourceDBSet

func (sd *sourceDBSet) Init() {

}

func (sd *sourceDBSet) Registry() map[string]Parameter {
	return map[string]Parameter{utils.SourceDBType: &sourceDB{}}
}
