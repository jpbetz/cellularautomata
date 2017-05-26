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
	"strconv"
	"fmt"
)

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

func main() {
	c := cli.NewCLI("cellular", "1.0.0")
	c.Args = os.Args[1:]

	input := make(chan io.InputEvent, 10)
	ui := sdlui.NewSdlUi(
		input,
		intEnvOrDefault("WIDTH", 60),
		intEnvOrDefault("HEIGHT", 40),
		intEnvOrDefault("CWIDTH", 15),
		intEnvOrDefault("CHEIGHT", 15),
		intEnvOrDefault("CBORDER", 1),
	)
	defer ui.Close()

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

func intEnvOrDefault(envName string, def int32) int32 {
	if wstr := os.Getenv(envName); len(wstr) > 0 {
		i, err := strconv.Atoi(wstr)
		if err != nil {
			panic(fmt.Sprintf("%s environment variable must be an integer, but was %s", envName, wstr))
		}
		return int32(i)
	} else {
		return def
	}
}
