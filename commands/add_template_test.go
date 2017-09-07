// Copyright (c) Alex Ellis, Eric Stoekl 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package commands

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

func Test_addTemplate(t *testing.T) {
	defer tearDown_fetch_templates(t)

	ts := httpTestServer(t)
	defer ts.Close()

	repository = ts.URL + "/owner/repo"
	faasCmd.SetArgs([]string{"add-template", repository})
	faasCmd.Execute()

	// Verify created directories
	if _, err := os.Stat("template"); err != nil {
		t.Fatalf("The directory %s was not created", "template")
	}
}

func Test_addTemplate_with_overwriting(t *testing.T) {
	defer tearDown_fetch_templates(t)

	ts := httpTestServer(t)
	defer ts.Close()

	repository = ts.URL + "/owner/repo"
	faasCmd.SetArgs([]string{"add-template", repository})
	faasCmd.Execute()

	// reset cacheCanWriteLanguage
	cacheCanWriteLanguage = make(map[string]bool)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	r := regexp.MustCompile(`(?m:overwriting is not allowed)`)

	faasCmd.SetArgs([]string{"add-template", repository})
	faasCmd.Execute()

	// reset cacheCanWriteLanguage
	cacheCanWriteLanguage = make(map[string]bool)

	if !r.MatchString(buf.String()) {
		t.Fatal(buf.String())
	}

	buf.Reset()

	faasCmd.SetArgs([]string{"add-template", repository, "--overwrite"})
	faasCmd.Execute()

	if r.MatchString(buf.String()) {
		t.Fatal()
	}

	// Verify created directories
	if _, err := os.Stat("template"); err != nil {
		t.Fatalf("The directory %s was not created", "template")
	}
}

func Test_addTemplate_error_no_arg(t *testing.T) {
	var buf bytes.Buffer

	faasCmd.SetArgs([]string{"add-template"})
	faasCmd.SetOutput(&buf)
	faasCmd.Execute()

	if !strings.Contains(buf.String(), "Error: A repository URL must be specified") {
		t.Fatal("Output does not contain the required string")
	}
}

func Test_addTemplate_error_not_valid_url(t *testing.T) {
	var buf bytes.Buffer

	faasCmd.SetArgs([]string{"add-template", "git@github.com:alexellis/faas-cli.git"})
	faasCmd.SetOutput(&buf)
	faasCmd.Execute()

	if !strings.Contains(buf.String(), "Error: The repository URL must be in the format https://github.com/<owner>/<repository>") {
		t.Fatal("Output does not contain the required string", buf.String())
	}
}

// httpTestServer returns a testing http server
func httpTestServer(t *testing.T) *httptest.Server {
	const sampleMasterZipPath string = "testdata/master_test.zip"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, err := os.Stat(sampleMasterZipPath); os.IsNotExist(err) {
			t.Error(err)
		}

		fileData, err := ioutil.ReadFile(sampleMasterZipPath)
		if err != nil {
			t.Error(err)
		}

		w.Write(fileData)
	}))

	return ts
}

func Test_repositoryUrlRegExp(t *testing.T) {
	var url string
	r := regexp.MustCompile(repositoryRegexpGithub)

	url = "http://github.com/owner/repo"
	if r.MatchString(url) {
		t.Errorf("Url %s must start with https", url)
	}

	url = "https://github.com/owner/repo.git"
	if r.MatchString(url) {
		t.Errorf("Url %s must not end with .git or must start with https", url)
	}

	url = "https://github.com/owner/repo"
	if !r.MatchString(url) {
		t.Errorf("Url %s must be valid", url)
	}
}
