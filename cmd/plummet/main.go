package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/template"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type Target struct {
	Output string                 `yaml:"output"`
	SQL    string                 `yaml:"sql"`
	Deps   []string               `yaml:"deps"`
	Config map[string]interface{} `yaml:"config"`
}

type PlummetFile struct {
	Targets map[string]Target `yaml:"targets"`
}

func executeTarget(targetName string, plummetFile *PlummetFile, visited map[string]bool, db *sql.DB, outputs map[string]interface{}) error {
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
		err := executeTarget(dep, plummetFile, visited, db, outputs)
		if err != nil {
			return fmt.Errorf("failed to execute dependency '%s' for target '%s': %v", dep, targetName, err)
		}
	}

	// Merge target-specific config with the outputs from dependencies
	config := make(map[string]interface{})
	for k, v := range target.Config {
		config[k] = v
	}
	for k, v := range outputs {
		config[k] = v
	}

	tmpl, err := template.New("sql").Parse(target.SQL)
	if err != nil {
		return fmt.Errorf("failed to parse SQL template for target '%s': %v", targetName, err)
	}

	var sqlBuffer bytes.Buffer
	err = tmpl.Execute(&sqlBuffer, config)
	if err != nil {
		return fmt.Errorf("failed to execute SQL template for target '%s': %v", targetName, err)
	}

	executedSQL := sqlBuffer.String()
	if target.Output == "" {
		_, err = db.Exec(executedSQL)
		if err != nil {
			return fmt.Errorf("failed to execute SQL for target '%s' with SQL: %s, error: %v", targetName, executedSQL, err)
		}
	} else {
		rows, err := db.Query(executedSQL)
		if err != nil {
			return fmt.Errorf("failed to query SQL for target '%s' with SQL: %s, error: %v", targetName, executedSQL, err)
		}
		defer rows.Close()
		var result interface{}
		if rows.Next() {
			err = rows.Scan(&result)
			if err != nil {
				return fmt.Errorf("failed to scan result for target '%s': %v", targetName, err)
			}
			outputs[targetName+"."+target.Output] = result
		}
		if err = rows.Err(); err != nil {
			return fmt.Errorf("error iterating through results for target '%s': %v", targetName, err)
		}
	}
	fmt.Printf("Successfully executed SQL for target %s\n", targetName)
	visited[targetName] = false

	return nil
}

func main() {
	app := &cli.App{
		Name:  "plummet",
		Usage: "A build system that runs SQL against a database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dbfile",
				Value:   "plummet.db",
				Usage:   "DuckDB database file",
				Aliases: []string{"d"},
			},
		},
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

			dbFile := c.String("dbfile")

			db, err := sql.Open("duckdb", dbFile)
			if err != nil {
				log.Fatalf("Unable to open database file: %v", err)
			}
			defer db.Close()

			if c.Args().Len() > 0 {
				targetName := c.Args().First()
				visited := make(map[string]bool)
				outputs := make(map[string]interface{})
				err := executeTarget(targetName, &plummetFile, visited, db, outputs)
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
