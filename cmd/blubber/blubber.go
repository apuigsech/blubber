package main

import (
	"os"
	"github.com/urfave/cli"
	"github.com/apuigsech/blubber"
)

func main() {
	app := cli.NewApp()

	app.Name = "blubber"
	app.Version = "0.0.1"
	app.Usage = "consolidate information of IP reputation"

	app.Author = "Albert Puigsech Galicia"
	app.Email = "albert@puigsech.com"

	app.EnableBashCompletion = false

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "providers, p",
		},
	}

	app.Action = runBlubber

	app.Run(os.Args)
}

func runBlubber(c *cli.Context) error {
	providers := c.String("providers")

	bl := blubber.NewBlubber()
	bl.LoadProvidersFromFile(providers)
	bl.MakeRequest(c.Args().Get(0), 4)

	return nil
}