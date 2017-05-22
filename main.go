package main

import (
	"github.com/jpbetz/cellularautomata/apps/conway"
	"github.com/jpbetz/cellularautomata/apps/guardduty"
	"github.com/jpbetz/cellularautomata/apps/langton"
	"github.com/jpbetz/cellularautomata/apps/wireworld"
	"github.com/jpbetz/cellularautomata/io"
	"github.com/jpbetz/cellularautomata/sdlui"
	"github.com/mitchellh/cli"
	"log"
	"os"
	"runtime"
)

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

func main() {
	input := make(chan io.InputEvent, 100)
	ui := sdlui.NewSdlUi(input)
	defer ui.Close()

	c := cli.NewCLI("cellular", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"conway": func() (cli.Command, error) {
			return &conway.ConwayCommand{
				UI: ui,
			}, nil
		},
		"guardduty": func() (cli.Command, error) {
			return &guardduty.GuardDutyCommand{
				UI: ui,
			}, nil
		},
		"wireworld": func() (cli.Command, error) {
			return &wireworld.WireWorldCommand{
				UI: ui,
			}, nil
		},
		"langton": func() (cli.Command, error) {
			return &langton.LangtonCommand{
				UI: ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
