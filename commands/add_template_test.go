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
	ts := httpTestServer(t)
	defer ts.Close()

	repository = ts.URL + "/owner/repo"
	faasCmd.SetArgs([]string{"add-template", repository})
	faasCmd.Execute()

	// Remove existing master.zip file if it exists
	if _, err := os.Stat(".cache/template-owner-repo.zip"); err == nil {
		t.Log("Found the archive file, removing it...")

		err := os.RemoveAll(".cache")
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("The archive was not downloaded: %s", err)
	}

	// Remove existing templates folder, if it exist
	if _, err := os.Stat("template/"); err == nil {
		t.Log("Found a template/ directory, removing it...")

		err := os.RemoveAll("template/")
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("Directory template was not created: %s", err)
	}
}

func Test_addTemplate_with_overwriting(t *testing.T) {
	ts := httpTestServer(t)
	defer ts.Close()

	repository = ts.URL + "/owner/repo"
	faasCmd.SetArgs([]string{"add-template", repository})
	faasCmd.Execute()

	// reset cache
	cache = make(map[string]bool)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	r := regexp.MustCompile(`(?m:overwriting is not allowed)`)

	faasCmd.SetArgs([]string{"add-template", repository})
	faasCmd.Execute()

	// reset cache
	cache = make(map[string]bool)

	if !r.MatchString(buf.String()) {
		t.Fatal(buf.String())
	}

	buf.Reset()

	faasCmd.SetArgs([]string{"add-template", repository, "--overwrite"})
	faasCmd.Execute()

	if r.MatchString(buf.String()) {
		t.Fatal()
	}

	// Remove existing master.zip file if it exists
	if _, err := os.Stat(".cache/template-owner-repo.zip"); err == nil {
		t.Log("Found the archive file, removing it...")

		err := os.RemoveAll(".cache")
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("The archive was not downloaded: %s", err)
	}

	// Remove existing templates folder, if it exist
	if _, err := os.Stat("template/"); err == nil {
		t.Log("Found a template/ directory, removing it...")

		err := os.RemoveAll("template/")
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("Directory template was not created: %s", err)
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

	if !strings.Contains(buf.String(), "Error: The given URL does not begin with http") {
		t.Fatal("Output does not contain the required string")
	}
}

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
