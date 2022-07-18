package terminal

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"plugin"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"dds/cli/terminal/api"
)

var (
	reCmd = regexp.MustCompile(`\S+`)
)

type shell struct {
	ctx        context.Context
	pluginsDir string
	commands   map[string]map[string]api.Command
	closed     chan struct{}
}

// newShell returns a new shell
func newShell() *shell {
	return &shell{
		pluginsDir: api.PluginsDir,
		commands:   make(map[string]map[string]api.Command),
		closed:     make(chan struct{}),
	}
}

// Init initializes the shell with the given context
func (gosh *shell) Init(ctx context.Context) error {
	gosh.ctx = ctx
	return gosh.loadCommands()
}

func (gosh *shell) loadCommands() error {

	if _, err := os.Stat(filepath.Join(api.GetHomeDir(), gosh.pluginsDir)); err != nil {
		return err
	}

	plugins, err := listFiles(filepath.Join(api.GetHomeDir(), gosh.pluginsDir), `.*_command.so`)
	if err != nil {
		return err
	}

	for _, cmdPlugin := range plugins {
		plug, err := plugin.Open(path.Join(filepath.Join(api.GetHomeDir(), gosh.pluginsDir), cmdPlugin.Name()))
		if err != nil {
			fmt.Printf("failed to open plugin %s: %v\n", cmdPlugin.Name(), err)
			continue
		}
		cmdSymbol, err := plug.Lookup(api.CmdSymbolName)
		if err != nil {
			fmt.Printf("plugin %s does not export symbol \"%s\"\n",
				cmdPlugin.Name(), api.CmdSymbolName)
			continue
		}
		commands, ok := cmdSymbol.(api.Commands)
		if !ok {
			fmt.Printf("Symbol %s (from %s) does not implement Commands interface\n", api.CmdSymbolName, cmdPlugin.Name())
			continue
		}
		if err := commands.Init(gosh.ctx); err != nil {
			fmt.Printf("%s initialization failed: %v\n", cmdPlugin.Name(), err)
			continue
		}
		for name, cmd := range commands.Registry() {
			gosh.commands[name] = cmd
		}
		gosh.ctx = context.WithValue(gosh.ctx, api.ShellCommand, gosh.commands)
	}
	return nil
}

// Open opens the shell for the given reader
func (gosh *shell) Open(r *bufio.Reader, sgl chan os.Signal) {
	loopCtx := gosh.ctx
	line := make(chan string)
	for {
		// start a goroutine to get input from the user
		go func(ctx context.Context, input chan<- string) {
			for {
				// TODO: future enhancement is to capture input key by key
				// to give command granular notification of key events.
				// This could be used to implement command autocompletion.
				fmt.Fprintf(ctx.Value(api.ShellStdout).(io.Writer), "%s ", api.GetPrompt(loopCtx))
				line, err := r.ReadString('\n')
				if err != nil {
					fmt.Fprintf(ctx.Value(api.ShellStderr).(io.Writer), "%v\n", err)
					continue
				}

				input <- line
				return
			}
		}(loopCtx, line)

		// wait for input or cancel
		select {
		case <-gosh.ctx.Done():
			close(gosh.closed)
			return
		case input := <-line:
			// 捕获终止信号，并丢弃
			if len(sgl) > 0 {
				sig := <-sgl
				switch sig {
				case syscall.SIGINT:
					fmt.Printf("\n")
				case syscall.SIGTSTP:
					fmt.Printf("\n")
				case syscall.SIGQUIT:
					fmt.Printf("\n")
				}
				break
			}
			var err error
			loopCtx, err = gosh.handle(loopCtx, input)
			if err != nil {
				fmt.Fprintf(loopCtx.Value(api.ShellStderr).(io.Writer), "%v\n", err)
			}
		}
	}
}

// Closed returns a channel that closes when the shell has closed
func (gosh *shell) Closed() <-chan struct{} {
	return gosh.closed
}

func (gosh *shell) handle(ctx context.Context, cmdLine string) (context.Context, error) {
	line := strings.TrimSpace(cmdLine)
	if line == "" {
		return ctx, nil
	}
	args := reCmd.FindAllString(line, -1)
	if args != nil {
		cmdName := strings.ToUpper(args[0])
		var NotFound bool
		for s := range gosh.commands {
			cmd, ok := gosh.commands[s][cmdName]
			if ok {
				NotFound = true
				return cmd.Exec(ctx, args)
			}
		}
		if !NotFound {
			return ctx, errors.New(fmt.Sprintf("command not found: %s\n", cmdName))
		}
	}
	return ctx, errors.New(fmt.Sprintf("unable to parse command line: %s", line))
}

func listFiles(dir, pattern string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			return nil, err
		}
		if matched {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

func OpenShell() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, api.ShellPrompt, api.SetPrompt())
	ctx = context.WithValue(ctx, api.ShellStdout, os.Stdout)
	ctx = context.WithValue(ctx, api.ShellStderr, os.Stderr)
	ctx = context.WithValue(ctx, api.ShellStdin, os.Stdin)

	shell := newShell()
	if err := shell.Init(ctx); err != nil {
		return errors.Errorf("\n\nfailed to initialize:\n", err)
	}

	cmdCount := len(shell.commands)
	var libraryNum int
	var cmdNum int
	if cmdCount > 0 {
		for _, m := range shell.commands {
			libraryNum += 1
			cmdNum = cmdNum + len(m)
		}
		fmt.Printf("\nLoad %d command library, %d commands in total...\n", libraryNum, cmdNum)
		fmt.Printf("Type help for available commands\n")
		fmt.Printf("\n")

	} else {
		fmt.Printf("Command set not found\n\n")
	}

	// 捕获终止信号
	SignalSet := make(chan os.Signal, 3)
	signal.Notify(SignalSet, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(w *sync.WaitGroup, sgl chan os.Signal) {
		defer w.Done()
		shell.Open(bufio.NewReader(os.Stdin), sgl)
	}(&wg, SignalSet)

	wg.Wait()
	return nil
}
