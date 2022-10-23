package main

import (
	"fmt"
	"github.com/892294101/dds/serialize"
	"github.com/892294101/go-mysql/canal"
	"github.com/892294101/go-mysql/replication"
	"os"
	"time"
)

func main() {
	re := new(canal.RowsEvent)
	re.Action = "insert"
	re.Header = &replication.EventHeader{HeadRaws: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}

	tme := &replication.TableMapEvent{Schema: []byte{'H', 'T', 'T', 'K'}, Table: []byte{'T', 'E', 'S', 'T'}}

	var rs []byte

	for i := 0; i < 50; i++ {
		rs = append(rs, []byte{1, 1, 1, 2, 4, 4, 4, 5, 9, 9, 1, 3, 9, 8, 2, 1, 6}...)

	}

	re.RowsEvent = &replication.RowsEvent{Table: tme, RowRaws: rs}

	d := new(serialize.DataEventV1)
	fmt.Println(time.Now().Format(time.RFC3339Nano))
	for i := 0; i < 5000000; i++ {
		_, err := d.Encode(re, true)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
	}
	fmt.Println(time.Now().Format(time.RFC3339Nano))

}
