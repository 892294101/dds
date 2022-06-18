package main

import (
	"context"
	"fmt"
	"io"

	"httk/cli/terminal/api"
)

const LibType = "EXTRACT"

var Commands testCmds

type helloCmd string

func (t helloCmd) Name() string {
	return string(t)
}
func (t helloCmd) Usage() string {
	return `hello`
}
func (t helloCmd) ShortDesc() string {
	return `prints greeting "hello there"`
}
func (t helloCmd) LongDesc() string {
	return t.ShortDesc()
}
func (t helloCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	out := ctx.Value(api.ShellStdout).(io.Writer)
	fmt.Fprintln(out, "hello there")
	return ctx, nil
}

type goodbyeCmd string

func (t goodbyeCmd) Name() string {
	return string(t)
}
func (t goodbyeCmd) Usage() string {
	return t.Name()
}
func (t goodbyeCmd) ShortDesc() string {
	return `prints message "bye bye"`
}
func (t goodbyeCmd) LongDesc() string {
	return t.ShortDesc()
}
func (t goodbyeCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	out := ctx.Value(api.ShellStdout).(io.Writer)
	fmt.Fprintln(out, "bye bye")
	return ctx, nil
}

// command module
type testCmds struct{}

func (t *testCmds) Init(ctx context.Context) error {
	return nil
}

func (t *testCmds) Registry() (re map[string]map[string]api.Command) {
	re = make(map[string]map[string]api.Command)
	re[LibType] = make(map[string]api.Command)
	re[LibType][api.Hello] = helloCmd(api.Hello)
	re[LibType][api.GoodBye] = goodbyeCmd(api.GoodBye)
	return re
}
