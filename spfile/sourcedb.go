package spfile

import (
	"fmt"
	"github.com/pkg/errors"
	"myGithubLib/dds/extract/mysql/utils"
	"regexp"
	"strings"
)

var (
	SourceDBType = "SOURCEDB"
	Port         = "PORT"
	DataBase     = "DATABASE"
	Types        = "TYPE"
	UserId       = "USERID"
	PassWord     = "PASSWORD"
)

var (
	DefaultPort     = "3306"
	DefaultDataBase = "test"
	DefaultTypes    = "mysql"
	DefaultUserId   = "root"
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
	fmt.Println("sourceDB Info: ", *s.ParamPrefix, *s.Address, *s.Port.Key, *s.Port.Value, *s.Database.Key, *s.Database.Value, *s.Type.Key, *s.Type.Value, *s.UserId.Key, *s.UserId.Value, *s.PassWord.Key, *s.PassWord.Value)
}

// 初始化参数可以支持的数据库和进程

func (s *sourceDB) Init() {
	s.SupportParams = map[string]map[string]string{
		MySQL: {
			Extract: Extract,
		},
		MariaDB: {
			Extract: Extract,
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
		case strings.EqualFold(sdb[i], SourceDBType):
			s.ParamPrefix = &sdb[i]
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", SourceDBType)
			}
			NextVal := &sdb[i+1]
			if KeyCheck(NextVal) {
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
		case strings.EqualFold(sdb[i], Port):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", Port)
			}
			NextVal := &sdb[i+1]
			if KeyCheck(NextVal) {
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
		case strings.EqualFold(sdb[i], DataBase):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", DataBase)
			}
			NextVal := &sdb[i+1]
			if KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.Database != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s.Database = &DatabaseModel{Key: &sdb[i], Value: NextVal}
			i += 1
		case strings.EqualFold(sdb[i], Types):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", Types)
			}
			NextVal := &sdb[i+1]
			if KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.Type != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s.Type = &TypeModel{Key: &sdb[i], Value: NextVal}
			i += 1
		case strings.EqualFold(sdb[i], UserId):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", UserId)
			}
			NextVal := &sdb[i+1]
			if KeyCheck(NextVal) {
				return errors.Errorf("keywords cannot be used: %s", *NextVal)
			}
			if s.UserId != nil {
				return errors.Errorf("Parameters cannot be repeated: %s", *NextVal)
			}
			s.UserId = &UserIdModel{Key: &sdb[i], Value: NextVal}
			i += 1
		case strings.EqualFold(sdb[i], PassWord):
			if i+1 > sdbLength {
				return errors.Errorf("%s value must be specified", PassWord)
			}
			NextVal := &sdb[i+1]
			if KeyCheck(NextVal) {
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
		s.Port = &PortModel{Key: &Port, Value: &DefaultPort}
	}
	if s.Database == nil {
		s.Database = &DatabaseModel{Key: &DataBase, Value: &DefaultDataBase}
	}
	if s.Type == nil {
		s.Type = &TypeModel{Key: &Types, Value: &DefaultTypes}
	}
	if s.UserId == nil {
		s.UserId = &UserIdModel{Key: &UserId, Value: &DefaultUserId}
	}
	if s.PassWord == nil {
		return errors.Errorf("%s Password must be specified", SourceDBType)
	}

	return nil
}

func KeyCheck(s *string) bool {
	key := map[string]string{
		strings.ToUpper(SourceDBType): SourceDBType,
		strings.ToUpper(Port):         Port,
		strings.ToUpper(DataBase):     DataBase,
		strings.ToUpper(Types):        Types,
		strings.ToUpper(UserId):       UserId,
		strings.ToUpper(PassWord):     PassWord,
	}
	_, ok := key[strings.ToUpper(*s)]
	return ok
}

type sourceDBSet struct{}

var sourceDBSetBus sourceDBSet

func (sd *sourceDBSet) Init() {

}

func (sd *sourceDBSet) Registry() map[string]Parameter {
	return map[string]Parameter{SourceDBType: &sourceDB{}}
}
