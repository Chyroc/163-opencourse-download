package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "163 open course download"
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) == 0 {
			return cli.ShowAppHelp(c)
		}
		return run(c.Args().Get(0))
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(url string) error {
	courses, err := getCourseList(url)
	if err != nil {
		return err
	}

	for _, course := range courses {
		fmt.Println(course)
	}
	return nil
}
