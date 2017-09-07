// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Flags that are to be added to commands

var (
	URL string
)

func init() {
	// Setup flags that are used only by this command (variables defined above)
	addTemplateCmd.Flags().StringVar(&URL, "url", "http://github.com/alexellis/faas-cli", "URL from which to pull git repo to grab 'template' dir from")

	faasCmd.AddCommand(addTemplateCmd)
}

// addTemplateCmd represents the addTemplate command
var addTemplateCmd = &cobra.Command{
	Use:   "add-template [--url URL]",
	Short: "Downloads templates from the specified github repo",
	Long: `Downloads the compressed github repo specified by [URL], and extracts the 'template'
directory from the root of the repo, if it exists.`,
	Example: "faas-cli add-template --url https://github.com/alexellis/faas-cli",
	Run: runAddTemplate,
}

func runAddTemplate(cmd *cobra.Command, args []string) {
	URL = strings.TrimRight(URL, "/")
	URL = URL + "/archive/master.zip"

	err := os.Setenv("templateUrl", URL)
	if err != nil {
		fmt.Printf("Error setting templateUrl env var: %v\n", err)
		os.Exit(1)
	}

	err = fetchTemplates()
	if err != nil {
		fmt.Printf("Error getting templates from URL: %v\n", err)
		os.Exit(1)
	}
}
