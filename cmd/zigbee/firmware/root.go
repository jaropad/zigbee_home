package firmware

import (
	"errors"
	"log"
	"runtime/debug"
	"slices"

	"github.com/urfave/cli/v2"
)

func RootCmd() *cli.Command {
	return &cli.Command{
		Name:  "firmware",
		Usage: "firmware operations like build & flash",
		Subcommands: []*cli.Command{
			buildCmd(),
			flashCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "workdir",
				Usage: "change the working directory for the build process (currently does not do anything)",
			},
		},
		Before: func(ctx *cli.Context) error {
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
				if len(vcsProps) == 3 {
					break
				}
			}

			slices.SortFunc(vcsProps, func(a, b debug.BuildSetting) int {
				if a.Key > b.Key {
					return 1
				}

				return -1
			})

			// Information that will help with investigation of issues
			log.Printf("build info: %s/%s, %s", buildInfo.GoVersion, buildInfo.Main.Version, vcsProps)

			return nil
		},
	}
}
