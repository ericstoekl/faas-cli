// Copyright (c) Alex Ellis, Eric Stoekl 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package commands

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var smallestZipFile = []byte{80, 75, 05, 06, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00}

func Test_PullTemplates(t *testing.T) {
	defer tearDown_fetch_templates(t)

	// Create fake server for testing.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Write out the minimum number of bytes to make the response a valid .zip file
		w.Write(smallestZipFile)

	}))
	defer ts.Close()

	err := fetchTemplates(ts.URL+"/owner/repo", false)
	if err != nil {
		t.Error(err)
	}

	// Verify the downloaded archive
	archive := ".cache/template-owner-repo.zip"
	if _, err := os.Stat(archive); err != nil {
		t.Fatalf("The archive %s was not downloaded", archive)
	}
}

// tearDown_fetch_templates_test cleans all files and directories created by the test
func tearDown_fetch_templates(t *testing.T) {

	// Remove existing archive file if it exists
	if _, err := os.Stat(".cache/"); err == nil {
		t.Log("Found a .cache/ directory, removing it...")

		err := os.RemoveAll(".cache")
		if err != nil {
			t.Log(err)
		}
	} else {
		t.Log("The archive was not downloaded: %s", err)
	}

	// Remove existing templates folder, if it exist
	if _, err := os.Stat("template/"); err == nil {
		t.Log("Found a template/ directory, removing it...")

		err := os.RemoveAll("template/")
		if err != nil {
			t.Log(err)
		}
	} else {
		t.Log("Directory template was not created: %s", err)
	}
}
