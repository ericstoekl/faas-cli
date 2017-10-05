// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package builder

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openfaas/faas-cli/utils"
)

// BuildImage construct Docker image from function parameters
func BuildImage(image string, handler string, functionName string, language string, nocache bool, squash bool) {

	switch language {
	case "node", "python", "ruby", "csharp", "python3":
		tempPath := createBuildTemplate(functionName, handler, language)

		fmt.Printf("Building: %s with %s template. Please wait..\n", image, language)

		flagStr := buildFlagString(nocache, squash, os.Getenv("http_proxy"), os.Getenv("https_proxy"))

		builder := strings.Split(fmt.Sprintf("docker build %s-t %s .", flagStr, image), " ")
		fmt.Println(strings.Join(builder, " "))
		ExecCommand(tempPath, builder)
		fmt.Printf("Image: %s built.\n", image)

		break
	case "Dockerfile", "dockerfile":
		tempPath := handler
		if _, err := os.Stat(handler); err != nil {
			fmt.Printf("Unable to build %s, %s is an invalid path\n", image, handler)
			fmt.Printf("Image: %s not built.\n", image)

			break
		}
		fmt.Printf("Building: %s with Dockerfile. Please wait..\n", image)

		flagStr := buildFlagString(nocache, squash, os.Getenv("http_proxy"), os.Getenv("https_proxy"))

		builder := strings.Split(fmt.Sprintf("docker build %s-t %s .", flagStr, image), " ")
		fmt.Println(strings.Join(builder, " "))
		ExecCommand(tempPath, builder)
		fmt.Printf("Image: %s built.\n", image)

	default:
		log.Fatalf("Language template: %s not supported. Build a custom Dockerfile instead.", language)
	}
}

// createBuildTemplate creates temporary build folder to perform a Docker build with language template
func createBuildTemplate(functionName string, handler string, language string) string {
	tempPath := fmt.Sprintf("./build/%s/", functionName)
	fmt.Printf("Clearing temporary build folder: %s\n", tempPath)

	clearErr := os.RemoveAll(tempPath)
	if clearErr != nil {
		fmt.Printf("Error clearing temporary build folder %s\n", tempPath)
	}

	fmt.Printf("Preparing %s %s\n", handler+"/", tempPath+"function")

	functionPath := tempPath + "/function"
	mkdirErr := os.MkdirAll(functionPath, 0700)
	if mkdirErr != nil {
		fmt.Printf("Error creating path %s - %s.\n", functionPath, mkdirErr.Error())
	}

	// Drop in directory tree from template
	utils.CopyFiles("./template/"+language, tempPath, true)

	// Overlay in user-function
	utils.CopyFiles(handler, tempPath+"function/", true)

	return tempPath
}

func buildFlagString(nocache bool, squash bool, httpProxy string, httpsProxy string) string {

	buildFlags := ""

	if nocache {
		buildFlags += "--no-cache "
	}
	if squash {
		buildFlags += "--squash "
	}

	if len(httpProxy) > 0 {
		buildFlags += fmt.Sprintf("--build-arg http_proxy=%s ", httpProxy)
	}

	if len(httpsProxy) > 0 {
		buildFlags += fmt.Sprintf("--build-arg https_proxy=%s ", httpsProxy)
	}

	return buildFlags
}
