package main

import (
	"fmt"
	mymap "github.com/892294101/cache-mmap/mmap"
	"github.com/892294101/dds/dbs/dat"
	"github.com/892294101/dds/dbs/ddslog"
	"github.com/892294101/dds/dbs/metadata"
	"github.com/892294101/dds/dbs/spfile"
	"github.com/grandecola/mmap"
	"os"
	"syscall"
	"time"
)

func main() {

	log, err := ddslog.InitDDSlog("HTTK_0002")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("1")

	md, err := metadata.InitMetaData("HTTK_0002", spfile.GetMySQLName(), spfile.GetExtractName(), log, metadata.LOAD)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer md.Close()
	fmt.Println("2")
	/*for i := 0; i < 1000000; i++ {
		x := uint64(i)
		if err := md.SetFilePosition(x, x); err != nil {
			panic(err)
		}

		if err := md.SetPosition(x, x); err != nil {
			panic(err)
		}

		if err := md.SetLastUpdateTime(uint64(time.Now().Unix())); err != nil {
			panic(err)
		}

	}*/

	pfile, err := spfile.LoadSpfile(fmt.Sprintf("%s.desc", "HTTK_0002"), spfile.UTF8, log, spfile.GetMySQLName(), spfile.GetExtractName())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("3")

	if err := pfile.Production(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("4")

	w := dat.NewWriteMgr()
	dbs := spfile.GetMySQLName()
	if err := w.Init(pfile, &dbs, md, log); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := w.LoadDatFile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := []byte("lskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;llskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;lksadjfllskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;lksadjfllskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;lksadjfllskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;lksadjfllskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;lksadjfllskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;lksadjfllskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;lksadjfllskdflsadjflsadjf;lksajdlkfjsaldjf;lsadjf;lsadjf;lsadjf;lksadjflksadjfl")
	for i := 0; i < 10000000; i++ {
		err := w.Write(b)
		if err != nil {
			fmt.Println(err)
			fmt.Println("I: ", i*len(b))
			os.Exit(1)
		}
	}

}

func t() {
	dstFile, err := os.OpenFile("/tmp/test.txt", os.O_CREATE|os.O_WRONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st := time.Now()
	defer func() {

		dstFile.Close()
		fmt.Println("文件写入耗时：", time.Now().Sub(st).Seconds(), "s")
	}()

	for i := 0; i < 100000; i++ {
		dstFile.WriteAt([]byte("1111111"), 0)
		dstFile.WriteAt([]byte("1111111"), 0)
		dstFile.WriteAt([]byte("1111111"), 0)
	}

}

func t2() {

	fd, err := os.OpenFile("/tmp/dat.test", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0740)
	if err != nil {
		panic(err)
	}
	_, err = fd.Stat()
	if err != nil {
		panic(err)
	}

	err = fd.Truncate(10000)

	x, err := mmap.NewSharedFileMmap(fd, 0, 10000, syscall.PROT_READ|syscall.PROT_WRITE)
	if err != nil {
		panic(err)
	}
	x.Lock()
	for i := 0; i < 1000000; i++ {
		x.WriteAt([]byte("1111111"), 0)
	}
	x.Unlock()

}

func t3() {

	f, err := mymap.NewMmap("/tmp/dat.test", os.O_RDWR|os.O_CREATE|os.O_EXCL, 10000)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < 1000000; i++ {
		f.WriteAt([]byte("1111111"), 0)
		f.WriteAt([]byte("1111111"), 9)
		f.WriteAt([]byte("1111111"), 10)

	}

}
