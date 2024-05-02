package main

import (
	"context"
	"errors"
	"log"
	"os"
	"runtime/debug"
	"slices"

	"github.com/ffenix113/zigbee_home/cmd/zigbee/firmware"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(log.Lmsgprefix | log.LstdFlags | log.Lshortfile)

	app := &cli.App{
		Name:  "zigbee",
		Usage: "Zigbee Home CLI application",
		Commands: []*cli.Command{
			firmware.RootCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "zigbee.yaml",
			},
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Bool("version") {
				return printVersion()
			}

			cli.ShowAppHelpAndExit(ctx, 0)

			return nil
		},
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func printVersion() error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("could not read build information")
	}

	vcsProps := make([]debug.BuildSetting, 0, 2)

	for _, setting := range buildInfo.Settings {
		if !slices.Contains([]string{"vcs.revision", "vcs.modified"}, setting.Key) {
			continue
		}

		vcsProps = append(vcsProps, setting)
		if len(vcsProps) == cap(vcsProps) {
			break
		}
	}

	slices.SortFunc(vcsProps, func(a, b debug.BuildSetting) int {
		if a.Key > b.Key {
			return 1
		}

		return -1
	})

	// Information that will help with investigation of
	log.Printf("%s, tag:%s, version:%s", buildInfo.GoVersion, buildInfo.Main.Version, vcsProps)

	return nil
}
