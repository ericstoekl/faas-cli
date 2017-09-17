// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const defaultTemplateRepository = "https://github.com/alexellis/faas-cli/archive/master.zip"
const cacheDirectory = "./.cache/"
const templateDirectory = "./template/"

var cacheCanWriteLanguage = make(map[string]bool)

// fetchTemplates fetch code templates from GitHub master zip file.
func fetchTemplates(templateURL string, overwrite bool) error {
	if len(templateURL) == 0 {
		templateURL = defaultTemplateRepository
	} else {
		templateURL = templateURL + "/archive/master.zip"
	}

	archive, err := fetchMasterZip(templateURL)
	if err != nil {
		return err
	}

	zipFile, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	log.Printf("Attempting to expand templates from %s\n", archive)

	for _, z := range zipFile.File {
		var rc io.ReadCloser

		relativePath := z.Name[strings.Index(z.Name, "/")+1:]
		if strings.Index(relativePath, "template/") != 0 {
			// Process only directories inside "template" at root
			continue
		}

		var language string
		if languageSplit := strings.Split(relativePath, "/"); len(languageSplit) > 2 {
			language = languageSplit[1]

			if len(languageSplit) == 3 && relativePath[len(relativePath)-1:] == "/" {
				// template/language/

				if !canWriteLanguage(language, overwrite) {
					fmt.Printf("Directory %s exits, overwriting is not allowed\n", relativePath)
					continue
				}
			} else {
				// template/language/*

				if !canWriteLanguage(language, overwrite) {
					continue
				}
			}
		} else {
			// template/
			continue
		}

		if strings.Index(relativePath, "template") == 0 {
			fmt.Printf("Found \"%s\"\n", relativePath)
			if rc, err = z.Open(); err != nil {
				break
			}

			if err = createPath(relativePath, z.Mode()); err != nil {
				break
			}

			// If relativePath is just a directory, then skip expanding it.
			if len(relativePath) > 1 && relativePath[len(relativePath)-1:] != "/" {
				if err = writeFile(rc, z.UncompressedSize64, relativePath, z.Mode()); err != nil {
					break
				}
			}
		}
	}

	fmt.Println("")

	return err
}

func fetchMasterZip(templateURL string) (string, error) {
	var err error

	if err = createPath(cacheDirectory, os.FileMode(0755)); err != nil {
		log.Fatalf("Could not create %s directory", cacheDirectory)
	}

	templateURLSplit := strings.Split(templateURL, "/")
	templateURLSplit = templateURLSplit[len(templateURLSplit)-4 : len(templateURLSplit)-2]
	templateRepository := strings.Join(templateURLSplit, "-")

	archive := fmt.Sprintf("%stemplate-%s.zip", cacheDirectory, templateRepository)

	if _, err = os.Stat(archive); err != nil {
		c := http.Client{}

		req, err := http.NewRequest("GET", templateURL, nil)
		if err != nil {
			log.Println(err.Error())
			return "", err
		}
		log.Printf("HTTP GET %s\n", templateURL)
		res, err := c.Do(req)
		if err != nil {
			log.Println(err.Error())
			return "", err
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
			return "", err
		}

		log.Printf("Writing %dKb to %s\n", len(bytesOut)/1024, archive)
		err = ioutil.WriteFile(archive, bytesOut, 0700)
		if err != nil {
			log.Println(err.Error())
			return "", err
		}
	}
	return archive, err
}

func writeFile(rc io.ReadCloser, size uint64, relativePath string, perms os.FileMode) error {
	var err error

	defer rc.Close()
	fmt.Printf("Writing %d bytes to %s\n", size, relativePath)
	f, err := os.OpenFile(relativePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.CopyN(f, rc, int64(size))

	return err
}

func createPath(relativePath string, perms os.FileMode) error {
	dir := filepath.Dir(relativePath)
	err := os.MkdirAll(dir, perms)
	return err
}

// canWriteLanguage tells whether the language can be processed or not
// if overwrite is activated, the directory template/language/ is removed before to keep it in sync
func canWriteLanguage(language string, overwrite bool) bool {
	if len(language) > 0 {
		if _, ok := cacheCanWriteLanguage[language]; ok {
			return cacheCanWriteLanguage[language]
		}
		dir := templateDirectory + language

		if _, err := os.Stat(dir); err == nil {
			// The directory template/language/ exists
			if overwrite == false {
				log.Printf("Directory %s exists, overwriting is not allowed\n", dir)
				cacheCanWriteLanguage[language] = false
			} else {
				// Clean up the directory to keep in sync with new updates
				if err := os.RemoveAll(dir); err != nil {
					log.Fatalf("Directory %s cannot be removed", dir)
					cacheCanWriteLanguage[language] = true
				}
			}
		} else {
			cacheCanWriteLanguage[language] = true
		}

		return cacheCanWriteLanguage[language]
	}

	return false
}
