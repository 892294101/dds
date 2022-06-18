package log

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

//判断文件是否存在
func IsFileExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// 判断文件夹是否存在
func PathExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return fi.IsDir()
	}
	return false
}

//获取文件Size
func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

//根据执行文件路径获取程序的HOME路径
func GetHomeDirectory() (homedir string) {
	file, _ := exec.LookPath(os.Args[0])
	ExecFilePath, _ := filepath.Abs(file)

	os := runtime.GOOS
	switch os {
	case "windows":
		execfileslice := strings.Split(ExecFilePath, `\`)
		HomeDirectory := execfileslice[:len(execfileslice)-2]
		for i, v := range HomeDirectory {
			if v != "" {
				if i > 0 {
					homedir += `\` + v
				} else {
					homedir += v
				}
			}
		}
	case "linux":
		execfileslice := strings.Split(ExecFilePath, "/")
		HomeDirectory := execfileslice[:len(execfileslice)-2]
		for _, v := range HomeDirectory {
			if v != "" {
				homedir += `/` + v
			}
		}
	default:
		logrus.Error(fmt.Sprintf("Unknown operation type: %s", os))
	}

	if homedir == "" {
		logrus.Error(fmt.Sprintf("Get program home directory failed: %s", homedir))
	}
	return homedir
}

//统一格式化输出
func Format(column_name string, value string) string {
	var cn string
	length := 31
	if len(column_name) < length {
		for i := 0; i < length-len(column_name); i++ {
			cn += " "
		}
		cn += " " + column_name + ": "
	} else {
		cn += column_name
	}
	return strings.ToUpper(cn) + value + "\n"
}

//切片转为字符类型
func SliceToString(kv []string) *string {
	var kwsb strings.Builder
	var kw string
	if len(kv) > 0 {
		for i, v := range kv {
			if i == len(kv)-1 {
				kwsb.WriteString(v)
			} else {
				kwsb.WriteString(v)
				kwsb.WriteString(" ")
			}
		}
		kw = kwsb.String()
		return &kw
	}
	return nil
}

func StringSliceEqualBCE(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

//字节转为Int64类型
func BytesToInt64(bys []byte) int {
	iv := strings.TrimSpace(string(bys))
	byteToInt, err := strconv.Atoi(iv)
	if err != nil {
		logrus.Warningf("Byte conversion to Int64 failed: %s [%s]", err, iv)
	}
	return byteToInt
}

func StringToInt(res string) (int, error) {
	iv := strings.TrimSpace(string(res))
	r, err := strconv.Atoi(iv)
	if err != nil {
		logrus.Warningf("string conversion to Int failed: %s [%s]", err, iv)
		return 0, err
	}
	return r, nil
}

//=======================================================================================================
//自定义Panic异常处理,调用方式: 例如Test()函数, 指定defer ErrorCheckOfRecover(Test)
func GetFunctionName(i interface{}, seps ...rune) string {
	u := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Entry()
	f, _ := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).FileLine(u)
	return f
}

func ErrorCheckOfRecover(n interface{}) {
	if err := recover(); err != nil {
		logrus.Errorf("Panic Message: %s", err)
		logrus.Errorf("Exception File: %s", GetFunctionName(n))
		logrus.Errorf("Print Stack Message: %s", string(debug.Stack()))
		logrus.Fatal("Abnormal exit of program")
	}
}

func GenerateIntHex() (uint8, error) {
	rs, err := strconv.Atoi(fmt.Sprintf("rd: 0x%x", fmt.Sprintf("%4v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(999))))
	return uint8(rs), err
}

func UInt64ToBytes(vn *uint64) []byte {
	byteBuf := bytes.NewBuffer([]byte{})
	binary.Write(byteBuf, binary.BigEndian, vn)
	return byteBuf.Bytes()
}

func BytesTouInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func Split(indata string, separate string) (o1, o2 string, err error) {
	res := strings.Split(indata, separate)
	if len(res) == 2 {
		return res[0], res[1], nil
	}
	return "", "", errors.Errorf("Format error: %s", indata)
}

func DeleteExtraSpace(s string) string {
	s1 := strings.Replace(s, "  ", " ", -1)      //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}

func InPlaceholder(ip int) (agentId *string) {
	var ips string
	for i := 1; i <= ip; i++ {
		if i == 1 {
			ips = ":" + strconv.Itoa(i)
		} else {
			ips += ",:" + strconv.Itoa(i)
		}
	}
	return &ips
}
