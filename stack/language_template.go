// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package stack

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// ParseYAML parse a YAML file into a LanguageTemplate struct.
func ParseYAMLForLanguageTemplate(yamlFile string) (*LanguageTemplate, error) {
	var langTemplate LanguageTemplate
	var err error
	var fileData []byte

	fileData, err = ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileData, &langTemplate)
	if err != nil {
		fmt.Printf("Error with YAML file\n")
		return nil, err
	}

	return &langTemplate, err
}

func IsValidTemplate(lang string) bool {
	var found bool
	if strings.ToLower(lang) == "dockerfile" {
		found = true
	}
	if _, err := os.Stat("./template/" + lang); err == nil {
		found = true
	}

	return found
}
