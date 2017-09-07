// Copyright (c) Alex Ellis, Eric Stoekl 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package commands

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"os"

)

func Test_addTemplate(t *testing.T) {
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
	defer ts.Close()

	URL = ts.URL
	faasCmd.SetArgs([]string{"add-template"})
	faasCmd.Execute()

	// Remove existing master.zip file if it exists
	if _, err := os.Stat("master.zip"); err == nil {
		t.Log("Found a master.zip file, removing it...")

		err := os.Remove("master.zip")
		if err != nil {
			t.Fatal(err)
		}
	}

	// Remove existing templates folder, if it exist
	if _, err := os.Stat("template/"); err == nil {
		t.Log("Found a template/ directory, removing it...")

		err := os.RemoveAll("template/")
		if err != nil {
			t.Fatal(err)
		}
	}
}
