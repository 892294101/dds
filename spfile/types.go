package spfile

import "github.com/892294101/dds/utils"

const (
	IpV4Reg  = "^((0|[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(0|[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-5])$"
	IpV4Port = "^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$"
	/*	SourceDBRegular = "(^)" +
		"(?i:(" + SourceDBType + "))  (\\s+) (((0|[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(0|[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-5])) (\\s+)" +
		"(?i:(" + Port + ")) (\\s+) (\\d+) (\\s+)" +
		"(?i:(" + DataBase + ")) (*) (\\s+) " +
		"(?i:(" + Types + ")) (" + MySQL + "|" + "MariaDB" + ") (\\s+) " +
		"(?i:(" + UserId + ")) (*) (\\s+) " +
		"(?i:(" + PassWord + ")) (*) (\\s+) " +
		"($)"*/
)

func GetMySQLName() string {
	return utils.MySQL
}

func GetMariaDBName() string {
	return utils.MariaDB
}
func GetOracleName() string {
	return utils.Oracle
}

func GetExtractName() string {
	return utils.Extract
}

func GetDBOptionsName() string {
	return utils.DBOptionsType
}

type Module interface {
	Init()
	Add(raw *string) error
	ListParamText() string
	GetParam() interface{}
}

type Parameter interface {
	put() string
	init()
	add(raw *string) error
	initDefault() error
	isType(raw *string, dbType *string, processType *string) error
	parse(raw *string) error
}

type Parameters interface {
	Module
	Registry() map[string]Parameter
}
