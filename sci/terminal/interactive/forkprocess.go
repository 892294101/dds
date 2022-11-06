package interactive

import (
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"path/filepath"
)

type ForkProcess struct {
	dir   string
	args0 string
	args  []string
}

func NewForkProcess() *ForkProcess {
	return new(ForkProcess)
}

func (f *ForkProcess) InitFork(execDir string, args0 string, args []string) error {
	f.dir = execDir
	f.args0 = args0
	f.args = args

	return nil
}

func (f *ForkProcess) Start() (int, error) {
	binaryFile := filepath.Join(f.dir, "bin", f.args0)
	/*	p, err := os.StartProcess(binaryFile, f.args, &os.ProcAttr{Dir: path.Join(f.dir),
		Env: append(os.Environ(), f.args...),
		Sys: &syscall.SysProcAttr{Setsid: true}})*/
	/*pid, err := syscall.ForkExec(binaryFile, f.args, &syscall.ProcAttr{
		Dir:   path.Join(f.dir),
		Env:   append(os.Environ(), f.args...),
		Files: []uintptr{0, 1, 2},
		Sys:   &syscall.SysProcAttr{Setsid: true},
	})*/

	//cmd := exec.Command(binaryFile, f.args...)
	/*var cred = &syscall.Credential{Uid: uint32(os.Getuid()), Gid: uint32(os.Getgid())}
	SysProcAttr := &syscall.SysProcAttr{Noctty: true, Credential: cred}*/

	/*var attr = os.ProcAttr{Env: os.Environ(), Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	process, err := os.StartProcess(binaryFile, f.args, &attr)
	if err != nil {
		return 0, errors.Errorf("StartProcess %v", err)
	}*/

	cmd := exec.Command(binaryFile, f.args...)
	//cred := &syscall.Credential{Uid: uint32(os.Getuid()), Gid: uint32(os.Getgid())}
	//cmd.SysProcAttr = &syscall.SysProcAttr{Ptrace: true, Setpgid: true, Foreground:false}

	if err := cmd.Start(); err != nil {
		return 0, err
	}
	// defer cmd.Process.Release()
	go func(c *exec.Cmd) {
		c.Wait()
	}(cmd)

	if cmd.Process == nil {
		count := 0
		for {
			t := time.NewTicker(time.Second)
			select {
			case <-t.C:
				if count == 8 {
					t.Stop()
					if cmd.Process != nil {
						t.Stop()
						return cmd.Process.Pid, nil
					} else {
						return 0, errors.Errorf("process group start timeout")
					}

				} else {
					if cmd.Process != nil {
						t.Stop()
						return cmd.Process.Pid, nil
					}
				}
				count++
			}
			t.Stop()
		}
	}
	return cmd.Process.Pid, nil
}