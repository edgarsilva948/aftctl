/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package docs

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var args struct {
	dir    string
	format string
}

// Cmd represents the Cobra command for generating documentation files.
var Cmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generates documentation files",
	Hidden: true,
	RunE:   Run,
}

func init() {
	flags := Cmd.Flags()

	flags.StringVarP(
		&args.dir,
		"dir",
		"d",
		"docs",
		"The directory where to save the documentation to",
	)

	flags.StringVarP(
		&args.format,
		"format",
		"f",
		"markdown",
		"The output format of the documentation. Valid options are 'markdown', 'man', 'restructured'",
	)
}

// CreateDir creates a directory if it doesn't already exist.
func CreateDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("falha ao criar o diretório: %s", err)
		}
	}
	return nil
}

// Run is the function executed when the Cobra command is run.
// It generates the documentation files.
func Run(cmd *cobra.Command, _ []string) (err error) {
	cmd.Root().DisableAutoGenTag = true

	if err := CreateDir(args.dir); err != nil {
		return err
	}

	switch args.format {
	case "markdown":
		err = doc.GenMarkdownTree(cmd.Root(), args.dir)
	case "man":
		year := time.Now().Year()
		header := &doc.GenManHeader{
			Title:   "aftctl",
			Section: "1",
			Source:  fmt.Sprintf("Copyright (c) %d Edgar Costa edgarsilva948@gmail.com", year),
		}
		err = doc.GenManTree(cmd.Root(), header, args.dir)
	}

	if err != nil {
		return err
	}

	fmt.Println("Documents generated successfully on", args.dir)

	return
}
