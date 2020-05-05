package main

import (
	"env-generator/adapter"
	"env-generator/model"
	"env-generator/usecase/service"
	"github.com/urfave/cli"
	"log"
	"os"
)

func run(args []string) error {
	envRepo := adapter.NewEnv()
	envService := service.NewEnv(envRepo)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "profile",
			Value: "default",
		},
		cli.StringFlag{
			Name:  "region",
			Value: "ap-northeast-1",
		},
		cli.StringFlag{
			Name:     "secret",
			Required: true,
		},
		cli.StringFlag{
			Name:  "version",
			Value: "AWSCURRENT",
			Usage: "secrets manager's version",
		},
	}

	app.Action = func(c *cli.Context) error {
		args := model.Args{
			Version: c.String("version"),
			Secret:  c.String("secret"),
			Region:  c.String("region"),
			Profile: c.String("profile"),
		}
		return envService.Exec(args)
	}

	return app.Run(args)
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
