// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"fmt"
	"os"

	"github.com/morikuni/aec"
	"github.com/openfaas/faas-cli/builder"
	"github.com/openfaas/faas-cli/stack"
	"github.com/spf13/cobra"
)

var (
	lang string
	list bool
)

func init() {
	newFunctionCmd.Flags().StringVar(&functionName, "name", "", "Name for your function")
	newFunctionCmd.Flags().StringVar(&lang, "lang", "", "Language or template to use")
	newFunctionCmd.Flags().StringVar(&gateway, "gateway", defaultGateway,
		"Gateway URL to store in YAML stack file")
	newFunctionCmd.Flags().StringVar(&image, "image", "", "Name for docker image")

	newFunctionCmd.Flags().BoolVar(&list, "list", false, "List available languages")

	faasCmd.AddCommand(newFunctionCmd)
}

// newFunctionCmd displays newFunction information
var newFunctionCmd = &cobra.Command{
	Use:   "new (--name=FUNCTION_NAME --lang=FUNCTION_LANGUAGE [--gateway=http://domain:port] | --list)",
	Short: "Create a new template in the current folder with the name given as name",
	Long: `The new command creates a new function based upon hello-world in the given
language or type in --list for a list of languages available.`,
	Example: `faas-cli new --name chatbot --lang node
  faas-cli new --name textparser --lang python --gateway http://mydomain:8080
  faas-cli new --list`,
	Run: runNewFunction,
}

func runNewFunction(cmd *cobra.Command, args []string) {
	if list == true {
		fmt.Printf(`Languages available as templates:
- node
- python
- python3
- ruby
- csharp
- Dockerfile

`)
		return
	}
	if len(functionName) == 0 {
		fmt.Println("You must supply a function name with the --name flag")
		return
	}

	if len(lang) == 0 {
		fmt.Println("You must supply a function language with the --lang flag")
		return
	}

	var stackFileName string
	var services stack.Services
	if len(yamlFile) == 0 {
		// We will create a new YAML file for this function
		stackFileName = functionName + ".yml"
	} else {
		// YAML file was passed in, so parse to see if it is valid
		parsedServices, err := stack.ParseYAMLFile(yamlFile, "", "")
		if err != nil {
			fmt.Printf("Specified file (" + yamlFile + ") is not valid YAML\n")
			return
		}
		services = *parsedServices

		stackFileName = yamlFile
	}

	PullTemplates("")

	if _, err := os.Stat(functionName); err == nil {
		fmt.Printf("Folder: %s already exists\n", functionName)
		return
	}

	if err := os.Mkdir("./"+functionName, 0700); err == nil {
		fmt.Printf("Folder: %s created.\n", functionName)
	}

	if lang != "Dockerfile" && lang != "dockerfile" {
		builder.CopyFiles("./template/"+lang+"/function/", "./"+functionName+"/", true)
	} else {
		if _, err := os.Create("./" + functionName + "/Dockerfile"); err != nil {
			fmt.Printf("Couldn't create Dockerfile")
			return
		}
	}

	if len(yamlFile) == 0 {
		services.Provider = stack.Provider{Name: "faas", GatewayURL: gateway}
		services.Functions = make(map[string]stack.Function)
	}

	if len(image) > 0 {
		services.Functions[functionName] = stack.Function{Language: lang, Image: image, Handler: "./" + functionName}
	} else {
		services.Functions[functionName] = stack.Function{Language: lang, Image: functionName, Handler: "./" + functionName}
	}

	fmt.Printf(aec.BlueF.Apply(figletStr))
	fmt.Println()
	fmt.Printf("Function created in folder: %s\n", functionName)

	stackWriteErr := stack.WriteYAMLData(&services, stackFileName)

	if stackWriteErr != nil {
		fmt.Printf("Error writing stack file %v\n", stackWriteErr)
	} else {
		fmt.Printf("Stack file written: %s\n", stackFileName)
	}

	return
}

func AppendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
