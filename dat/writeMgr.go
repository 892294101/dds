package dat

import (
	"encoding/binary"
	"fmt"
	"github.com/892294101/cache-mmap/mmap"
	"github.com/892294101/dds-metadata"
	"github.com/892294101/dds-spfile"
	"github.com/892294101/dds-utils"
	"github.com/892294101/dds/serialize"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func (w *WriteCache) Init(s *ddsspfile.Spfile, dbType *string, mdh ddsmetadata.MetaData, log *logrus.Logger) error {
	trail := s.GetTrail()
	w.MaxSize = *trail.GetSizeValue() * 1024 * 1024
	w.md = mdh
	w.log = log
	w.log.Debugf("Get the program home directory")
	home, err := ddsutils.GetHomeDirectory()
	if err != nil {
		return err
	}
	w.log.Debugf("program home directory %v", *home)
	w.ProcName = *s.GetProcessName()
	dir := *trail.GetDir()
	ok := strings.HasPrefix(dir, "./")
	if ok {
		ind := strings.LastIndex(dir, "/")
		if ind == -1 {
			return errors.Errorf("Trail directory extraction error: %s", dir)
		}
		w.DatDir = path.Join(*home, dir[:ind])
		w.Prefix = dir[ind+1:]
	} else {
		ok := strings.HasPrefix(dir, "/")
		if ok {
			ind := strings.LastIndex(dir, "/")
			if ind == -1 {
				return errors.Errorf("Trail directory extraction error: %s", dir)
			}
			w.DatDir = dir[:ind]
			w.Prefix = dir[ind+1:]
		}
	}

	if len(w.DatDir) == 0 || len(w.Prefix) == 0 {
		return errors.Errorf("Failed to load trail directory when loading writer: %s", dir)
	}

	return nil
}

func (w *WriteCache) CreateDatFile() error {
	w.log.Debugf("open file %s%09d", w.Prefix, w.Seq)

	file, err := mmap.NewMmap(w.CurrentFile, os.O_CREATE|os.O_RDWR, int64(w.MaxSize))
	if err != nil {
		return errors.Errorf("file creation failed: %s%09d", w.Prefix, w.Seq)
	} else {
		w.file = file
	}

	return nil
}

func (w *WriteCache) NextDatFile() error {
	w.log.Infof("switch file %s%09d due to EOF", w.Prefix, w.Seq)
	if err := w.WriteFileEndMark(); err != nil {
		return err
	}

	if err := w.Sync(); err != nil {
		return err
	}

	if w.file != nil {
		if err := w.Close(); err != nil {
			return errors.Errorf("Closing file %s%09d error: %s", w.Prefix, w.Seq, err)
		} else {
			w.log.Infof("Closing file %s%09d succeeded", w.Prefix, w.Seq)
		}
	}

	w.Seq = w.Seq + 1
	w.CurrentFile = filepath.Join(w.DatDir, fmt.Sprintf("%s%09d", w.Prefix, w.Seq))
	if err := w.CreateDatFile(); err != nil {
		return err
	}
	w.Rba = 0
	w.log.Infof("Switching to next trail file %s%09d", w.Prefix, w.Seq)
	if err := w.WriteFileBeginMark(); err != nil {
		return err
	}
	w.log.Debugf("current file sequence %v rba %v", w.Seq, w.Rba)
	return nil
}

func (w *WriteCache) OpenFile() error {
	w.log.Debugf("open file %s%09d", w.Prefix, w.Seq)
	file, err := mmap.NewMmap(w.CurrentFile, os.O_RDWR, int64(w.MaxSize))
	if err != nil {
		return errors.Errorf("open file failed: %s", err)
	}
	w.file = file
	return nil
}

func (w *WriteCache) WriteFileBeginMark() error {
	w.log.Debugf("write file begin flag")
	var dfs serialize.Serialize

	dfs = &serialize.FileMarkV1{}
	b, err := dfs.Encode(&serialize.FileMark{WriteMark: serialize.FileBegin, WriteTime: uint64(time.Now().UnixNano())}, false)
	if err != nil {
		return errors.Errorf("Error writing header information of file: %v", err)
	}

	if err := w.writeOffset(uint64(len(b))); err != nil {
		return err
	}

	if err := w.writeData(b); err != nil {
		return err
	}
	w.log.Debugf("file begin flag size %v", len(b)+8)
	return w.UpdateFileRBA()
}

func (w *WriteCache) WriteFileEndMark() error {
	w.log.Debugf("write file end flag")
	var dfs serialize.Serialize
	dfs = &serialize.FileMarkV1{}
	b, err := dfs.Encode(&serialize.FileMark{WriteMark: serialize.FileEnd, WriteTime: uint64(time.Now().UnixNano())}, false)
	if err != nil {
		return errors.Errorf("Error writing header information of file: %s", err)
	}

	if err := w.writeOffset(uint64(len(b))); err != nil {
		return err
	}

	if err := w.writeData(b); err != nil {
		return err
	}

	w.log.Debugf("file end flag size %v", len(b)+8)
	return w.UpdateFileRBA()
}

// 写入数据偏移
func (w *WriteCache) writeOffset(offset uint64) error {
	w.Dirty = true
	w.lock.Lock()
	defer w.lock.Unlock()
	ost := make([]byte, 8)
	binary.LittleEndian.PutUint64(ost, offset)
	n, err := w.file.WriteAt(ost, int64(w.Rba))
	if err != nil {
		return err
	}
	w.Rba += uint64(n)
	return nil
}

// 写入数据
func (w *WriteCache) writeData(b []byte) error {
	w.Dirty = true
	w.lock.Lock()
	defer w.lock.Unlock()
	n, err := w.file.WriteAt(b, int64(w.Rba))
	if err != nil {
		return err
	}
	w.Rba += uint64(n)
	return nil
}

func (w *WriteCache) UpdateFileRBA() error {
	return w.md.SetFilePosition(w.Seq, w.Rba)
}

func (w *WriteCache) GetFileRBA() (*uint64, *uint64, error) {
	fn, rba, err := w.md.GetFilePosition()
	if err != nil {
		return nil, nil, err
	}
	w.log.Debugf("sequence: %v rba: %v", *fn, *rba)
	return fn, rba, err
}

func (w *WriteCache) LoadDatFile() error {
	fn, rba, err := w.GetFileRBA()
	if err != nil {
		return err
	}
	w.flushPeriodTime = time.Second
	if (fn != nil && rba != nil) && (*fn == 0 && *rba == 0) {
		w.CurrentFile = filepath.Join(w.DatDir, fmt.Sprintf("%s%09d", w.Prefix, *fn))
		if !ddsutils.PathExists(w.DatDir) {
			return errors.Errorf("directory does not exist: %s", w.CurrentFile)
		}

		if ddsutils.IsFileExist(w.CurrentFile) {
			return errors.Errorf("file already exists: %s", w.CurrentFile)
		} else {
			if err := w.CreateDatFile(); err != nil {
				return err
			}
			w.Seq = *fn
			w.Rba = *rba
			return w.WriteFileBeginMark()
		}

	} else if (fn != nil && rba != nil) && (*fn >= 0 && *rba > 0) {
		w.CurrentFile = filepath.Join(w.DatDir, fmt.Sprintf("%s%09d", w.Prefix, *fn))
		w.Seq = *fn
		w.Rba = *rba
	}
	if err := w.OpenFile(); err != nil {
		return err
	}
	w.wg.Add(1)
	go w.periodicFlush()
	return nil
}

func (w *WriteCache) CheckFull(b []byte) error {
	np := int(w.Rba) + len(b) + 8 + 22
	if np > w.MaxSize {
		return w.NextDatFile()
	}
	return nil
}

func (w *WriteCache) Write(b []byte) error {
	if err := w.CheckFull(b); err != nil {
		return err
	}
	if err := w.writeOffset(uint64(len(b))); err != nil {
		return err
	}

	if err := w.writeData(b); err != nil {
		return err
	}

	// 更新RBA地址
	if err := w.UpdateFileRBA(); err != nil {
		return err
	}

	return nil
}

func (w *WriteCache) periodicFlush() {
	defer w.wg.Done()
	timer := &time.Timer{C: make(chan time.Time)}
	timer = time.NewTimer(w.flushPeriodTime)
	var drainFlag bool
	for {
		if w.flushPeriodTime != 0 {
			if !drainFlag && !timer.Stop() {
				<-timer.C
			}
			timer.Reset(w.flushPeriodTime)
			drainFlag = false
		}

		select {
		case <-w.quit:
			return
		case <-w.drain:
			_ = w.file.Flush()
		case <-timer.C:
			drainFlag = true
			_ = w.file.Flush()
		}
	}
}

func (w *WriteCache) Sync() error {
	if err := w.file.Flush(); err != nil {
		return err
	}
	w.Dirty = false
	return nil
}

func (w *WriteCache) Close() error {
	w.log.Infof("Closing file %s%09d rba %v", w.Prefix, w.Seq, w.Rba)
	return w.file.Close()
}

func (w *WriteCache) CloseDat() error {
	w.log.Infof("Synchronizing %s%09d to disk", w.Prefix, w.Seq)
	if err := w.Sync(); err != nil {
		return err
	}
	return w.Close()
}

func NewWriteMgr() *WriteCache {
	return new(WriteCache)
}
