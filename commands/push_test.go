package commands

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/openfaas/faas-cli/builder"
	"github.com/openfaas/faas-cli/test"
)

func TestExecCommandHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// println("Mocked stdout:", os.Getenv("STDOUT"))
	fmt.Fprintf(os.Stdout, os.Getenv("STDOUT"))
	i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
	os.Exit(i)
}

func TestPrintDate(t *testing.T) {
	expectedOut := "Sun Aug 201"
	builder.MockExec(1, expectedOut)

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"push",
			"--yaml=testdata/sample_stack.yml",
		})
		faasCmd.Execute()
	})

	if stdOut != expectedOut {
		t.Errorf("Expected %q, got %q", expectedOut, stdOut)
	}
}
