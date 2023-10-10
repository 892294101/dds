package main

import (
	"flag"
	"github.com/892294101/dds-spfile"
	"github.com/892294101/dds-utils"
	oramysql "github.com/892294101/dds/extract/mysql"
	_ "net/http/pprof"
	"os"
	"strings"
)

var processName = flag.String("processid", "", "Please enter the process id name")

//var bg = flag.String("background", "n", "Put it to run in the background")

func main() {
	/*dat := []byte("123123  dfsadf")
	dat2 := []byte("h 45")
	var b utils.BytesBufferPool
	b.Init()
	b.BytePut(dat)
	b.BytePut([]byte("123123"))
	fmt.Println(b.ByteGet())

	fmt.Println(b.ByteGet())

	b.BytePut(dat2)
	fmt.Println(b.ByteGet())

	fmt.Println(b.ByteGet())*/

	/*go func() {
		_ = http.ListenAndServe("0.0.0.0:8081", nil)
	}()*/

	/*	os.Exit(1)*/

	flag.Parse()
	if processName == nil || len(*processName) == 0 {
		os.Exit(1)
	}
	dds_utils.GlobalProcessID = strings.ToUpper(*processName)
	canal := oramysql.NewMySQLSync()
	canal.InitSyncerConfig(*processName, dds_spfile.GetMySQLName(), dds_spfile.GetExtractName())
	canal.StartSyncToStream()
}

/*func init() {
	flag.Parse()
	if processName == nil || len(*processName) == 0 {
		os.Exit(1)
	}
	dir, _ := utils.GetHomeDirectory()
	lockFile := filepath.Join(*dir, "tmp", strings.ToUpper(*processName)+".lock")
	if utils.IsFileExist(lockFile) {
		fmt.Printf("Process group is starting\n")
		os.Exit(1)
	} else {
		switch {
		case strings.EqualFold(*bg, "y"):
			_, err := os.Create(lockFile)
			if err != nil {
				fmt.Fprintf(os.Stdout, "Lock file creation failed: %v\n", err)
				os.Exit(1)
			} else {
				cmd := exec.Command(os.Args[0], os.Args[1], os.Args[2])
				if err := cmd.Start(); err != nil {
					fmt.Fprintf(os.Stdout, "Process group start failed: %v\n", err)
				}
				cmd.Process.Release()
				os.Remove(lockFile)
				os.Exit(0)
			}
		case strings.EqualFold(*bg, "n"):
			// 当为n时要继续执行，以进入到mian函数
		default:
			fmt.Fprintf(os.Stdout, "The background parameter only supports y and n\n")
			os.Exit(1)
		}

	}
}*/
