package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type Target struct {
	Output string   `yaml:"output"`
	SQL    string   `yaml:"sql"`
	Deps   []string `yaml:"deps"`
}

type PlummetFile struct {
	Targets map[string]Target `yaml:"targets"`
}

func main() {
	app := &cli.App{
		Name:  "plummet",
		Usage: "A build system that runs SQL against a database",
		Action: func(c *cli.Context) error {
			file, err := os.ReadFile("plummet.yml")
			if err != nil {
				log.Fatalf("Unable to read plummet.yml: %v", err)
			}

			var plummetFile PlummetFile
			err = yaml.Unmarshal(file, &plummetFile)
			if err != nil {
				log.Fatalf("Unable to parse plummet.yml: %v", err)
			}

			fmt.Println("Available targets:")
			for target := range plummetFile.Targets {
				fmt.Println(target)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
