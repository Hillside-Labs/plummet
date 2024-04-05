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

func executeTarget(targetName string, plummetFile *PlummetFile, visited map[string]bool) error {
	if visited[targetName] {
		return fmt.Errorf("circular dependency detected on target '%s'", targetName)

	}
	visited[targetName] = true

	target, ok := plummetFile.Targets[targetName]
	if !ok {
		return fmt.Errorf("target '%s' not found", targetName)
	}

	// Execute dependencies first
	for _, dep := range target.Deps {
		err := executeTarget(dep, plummetFile, visited)
		if err != nil {
			return fmt.Errorf("failed to execute dependency '%s' for target '%s': %v", dep, targetName, err)
		}
	}

	// Here you would add the logic to execute the SQL against the database
	// and handle the output, for now we just print the SQL to be executed.
	fmt.Printf("Executing SQL for target %s: %s\n", targetName, target.SQL)
	visited[targetName] = false

	return nil
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

			if c.Args().Len() > 0 {
				targetName := c.Args().First()
				visited := make(map[string]bool)
				err := executeTarget(targetName, &plummetFile, visited)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Println("Available targets:")
				for target := range plummetFile.Targets {
					fmt.Println(target)
				}
			}

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
