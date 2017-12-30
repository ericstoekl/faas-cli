// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/morikuni/aec"
)

var ExecCommand = execCommand
var mockedExitStatus = 0
var mockedStdout string

func MockExec(exitStatus int, output string) {
	mockedExitStatus = exitStatus
	mockedStdout = output
	ExecCommand = fakeExecCommand
}

// ExecCommand run a system command
func execCommand(tempPath string, builder []string) {
	targetCmd := exec.Command(builder[0], builder[1:]...)
	targetCmd.Dir = tempPath
	targetCmd.Stdout = os.Stdout
	targetCmd.Stderr = os.Stderr
	targetCmd.Start()
	err := targetCmd.Wait()
	if err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", builder)
		log.Fatalf(aec.RedF.Apply(errString))
	}
}

func fakeExecCommand(tempPath string, builder []string) {
	cs := []string{"-test.run=TestExecCommandHelper", "--", builder[0]}
	cs = append(cs, builder[1:]...)
	cmd := exec.Command(os.Args[0], cs...)
	es := strconv.Itoa(mockedExitStatus)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + mockedStdout,
		"EXIT_STATUS=" + es}
	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", builder)
		log.Fatalf(aec.RedF.Apply(errString))
	}
}
