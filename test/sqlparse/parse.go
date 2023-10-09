package main

import (
	"context"
	"fmt"
	"github.com/jiunx/xsqlparser/parser"
	"github.com/jiunx/xsqlparser/util"
	"os"
)

func main() {
	sql := ` insert into WTEST
		values ('T','测试1','测试2','测试3','测试4','测试5','测试6',32.02,33, 23432849.2134231, 234234.123123123123213E+32, 2142342134.1232134E-12, sysdate,
		   systimestamp, TO_TIMESTAMP_TZ('2020-12-12 11:12:13.123456789 +8:00', 'YYYY-MM-DD HH24:MI:SS.XFF TZH:TZM'),
		   systimestamp, INTERVAL '2-6' YEAR TO MONTH,  INTERVAL '7' DAY, utl_raw.cast_to_raw('ASDFKLSADF-SAD-SADF-SADF-SADF-SDAF'), HEXTORAW('0e78e8be4b880e6af9b111'))
		   `

	parserFactory := parser.GetParserFactoryInstance()
	mysqlParser, err := parserFactory.GetParser(util.Oracle)
	if err != nil {
		fmt.Println("mysqlParser: ",err)
		os.Exit(1)
	}

	_, err2 := mysqlParser.Parse(context.Background(), sql)
	if err2 != nil {
		fmt.Println("smt: ",err2)
		os.Exit(1)
	}


}
