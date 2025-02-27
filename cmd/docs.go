// Copyright © 2019 IBM Corporation and others.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	flag "github.com/spf13/pflag"
)

//generate Doc file (.md) for cmds in package

func generateDoc(commandDocFile string) error {

	if commandDocFile == "" {
		return errors.New("no docFile specified")
	}
	dir := filepath.Dir(commandDocFile)

	if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
		mkdirErr := os.MkdirAll(dir, 0755)
		if mkdirErr != nil {
			Error.log("Could not create doc file directory: ", mkdirErr)
			return mkdirErr
		}
	}
	docFile, createErr := os.Create(commandDocFile)
	if createErr != nil {
		Error.log("Could not create doc file (.md): ", createErr)
		return createErr
	}

	defer docFile.Close()

	preAmble := "---\ntitle: Appsody CLI Reference\npath: /docs/using-appsody/cli-commands\nsection: Using Appsody\n---\n# Appsody CLI\n"
	preAmbleBytes := []byte(preAmble)
	_, preambleErr := docFile.Write(preAmbleBytes)
	if preambleErr != nil {
		Error.log("Could not write to markdown file:", preambleErr)
		return preambleErr
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		newbase := strings.ReplaceAll(base, "_", "-")
		return "#" + newbase
	}
	commandArray := []*cobra.Command{rootCmd, buildCmd, bashCompletionCmd, debugCmd, deployCmd, extractCmd, initCmd, listCmd, repoCmd, addCmd, repoListCmd, removeCmd, runCmd, stopCmd, testCmd, versionCmd}
	for _, cmd := range commandArray {

		markdownGenErr := doc.GenMarkdownCustom(cmd, docFile, linkHandler)

		if markdownGenErr != nil {
			Error.log("Doc file generation failed: ", markdownGenErr)
			return markdownGenErr
		}
	}
	return nil

}

// docs command is used to generate markdown file for all the appsody commands
var docsCmd = &cobra.Command{
	Use:    "docs",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {

		Debug.log("Running appsody docs command.")
		err := generateDoc(docFile)
		if err != nil {
			Error.log("appsody docs command failed with error: ", err)
			os.Exit(1)
		}
		Debug.log("appsody docs command completed successfully.")
	},
}

var docFile string

func init() {
	rootCmd.AddCommand(docsCmd)
	docFlags := flag.NewFlagSet("", flag.ContinueOnError)

	docFlags.StringVar(&docFile, "docFile", "", "Specify the file to contain the generated documentation.")
	docsCmd.PersistentFlags().AddFlagSet(docFlags)
}
