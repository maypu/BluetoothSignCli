package main

import (
	"BluetoothSignCli/utils"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
)

func main()  {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "show list of all hospital",
				Action:  func(c *cli.Context) error {
					fmt.Println("hospital list:")
					utils.List()
					return nil
				},
			},
			{
				Name:    "now",
				Aliases: []string{"n"},
				Usage:   "get current hospital by Mac Address",
				Action:  func(c *cli.Context) error {
					utils.Now()
					return nil
				},
			},
			{
				Name:    "set",
				Aliases: []string{"s"},
				Usage:   "set hospital option by index",
				Action:  func(c *cli.Context) error {
					if c.Args().First() == "" {
						fmt.Println("-s must has a index argument, please add a int argument by --list")
						return nil
					}
					argFirst,_ := strconv.Atoi(c.Args().First())
					utils.Set(argFirst)
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
