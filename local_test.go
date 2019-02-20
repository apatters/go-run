// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package run_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/apatters/go-run"
	"github.com/stretchr/testify/assert"
)

func TestLocal_RunSuccess(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err := l.Run("/bin/true")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.NoError(t, err)
}

func TestLocal_RunFail(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err := l.Run("/bin/false")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NotZero(t, code)
	assert.NoError(t, err)
}

func TestLocal_RunExit(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err := l.Run("/bin/sh", "-c", "exit 6")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 6)
	assert.NoError(t, err)
}

func TestLocal_RunOutput(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err := l.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false",
		"/xyzzy")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Equal(t, stdout, "/bin/false\n/bin/true\n")
	assert.Regexp(
		t,
		regexp.MustCompile(`/bin/ls: cannot access .*xyzzy.*: No such file or directory`),
		stderr)
	assert.Equal(t, code, 2)
	assert.NoError(t, err)
}

func TestLocal_RunStdin(t *testing.T) {
	stdinStr := "Hello, world"
	l := run.NewLocal(run.LocalConfig{
		Stdin: strings.NewReader(stdinStr),
	})
	stdout, stderr, code, err := l.Run(
		"/usr/bin/tr",
		"[:upper:]",
		"[:lower:]")
	t.Logf("stdin = %q", stdinStr)
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Equal(t, strings.ToLower(stdinStr), stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.NoError(t, err)
}

func TestLocal_RunStdout(t *testing.T) {
	var b bytes.Buffer
	l := run.NewLocal(run.LocalConfig{
		Stdout: bufio.NewWriter(&b),
	})
	stdout, stderr, code, err := l.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false",
		"/xyzzy")
	t.Logf("b = %q", b.String())
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Equal(t, b.String(), "/bin/false\n/bin/true\n")
	assert.Empty(t, stdout)
	assert.Regexp(
		t,
		regexp.MustCompile(`/bin/ls: cannot access .*xyzzy.*: No such file or directory`),
		stderr)
	assert.Equal(t, code, 2)
	assert.NoError(t, err)
}

func TestLocal_RunStderr(t *testing.T) {
	var b bytes.Buffer
	l := run.NewLocal(run.LocalConfig{
		Stderr: bufio.NewWriter(&b),
	})
	stdout, stderr, code, err := l.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false",
		"/xyzzy")
	t.Logf("b = %q", b.String())
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Equal(t, stdout, "/bin/false\n/bin/true\n")
	assert.Regexp(
		t,
		regexp.MustCompile(`/bin/ls: cannot access .*xyzzy.*: No such file or directory`),
		b.String())
	assert.Empty(t, stderr)
	assert.Equal(t, code, 2)
	assert.NoError(t, err)
}

func TestLocal_RunEnv(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{
		Env: []string{"FIRST=1st", "SECOND=2nd"},
	})
	stdout, stderr, code, err := l.Run("/bin/env")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	assert.Equal(t, "FIRST=1st\nSECOND=2nd\n", stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 0)
	assert.NoError(t, err)
}

func TestLocal_RunDir(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{
		Dir: "/",
	})
	savedDir, _ := os.Getwd()
	stdout, stderr, code, err := l.Run("/bin/pwd")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	curDir, _ := os.Getwd()
	assert.Equal(t, savedDir, curDir, "Working directory not restored.")
	assert.Equal(t, "/\n", stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 0)
	assert.NoError(t, err)
}

func TestLocal_ShellSuccess(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err := l.Shell("exit 0")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.NoError(t, err)
}

func TestLocal_ShellFail(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err := l.Shell("exit 1")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 1)
	assert.NoError(t, err)
}

func TestLocal_ShellOutput(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err := l.Shell("cd /bin && ls true false xyzzy")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Equal(t, stdout, "false\ntrue\n")
	assert.Regexp(
		t,
		regexp.MustCompile(`.*ls: cannot access .*xyzzy.*: No such file or directory`),
		stderr)
	assert.Equal(t, code, 2)
	assert.NoError(t, err)
}

func TestLocal_NewShell(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	l.ShellExecutable = "/bin/bash"
	stdout, stderr, code, err := l.Shell("xyzzy")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Regexp(
		t,
		regexp.MustCompile(`/bin/bash: xyzzy: command not found\n`),
		stderr)
	assert.Equal(t, code, 127)
	assert.NoError(t, err)
}

func TestLocal_FormatRun(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})

	msg := l.FormatRun("uname")
	t.Logf("cmd = %q", "uname")
	t.Logf("msg = %q", msg)
	assert.Equal(t, msg, "uname")

	msg = l.FormatRun("uname", "-a")
	t.Logf("cmd = %q", "uname -a")
	t.Logf("msg = %q", msg)
	assert.Equal(t, msg, "uname -a")
}

func TestLocal_ShellEnv(t *testing.T) {
	envVars := []string{"FIRST=1st", "SECOND=2nd"}
	l := run.NewLocal(run.LocalConfig{
		Env: envVars,
	})
	stdout, stderr, code, err := l.Shell("/bin/env")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	stdoutLines := strings.Split(stdout, "\n")
	assert.Subset(t, stdoutLines, envVars)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 0)
	assert.NoError(t, err)
}

func TestLocal_ShellDir(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{
		Dir: "/",
	})
	savedDir, _ := os.Getwd()
	stdout, stderr, code, err := l.Shell("pwd")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	curDir, _ := os.Getwd()
	assert.Equal(t, savedDir, curDir, "Working directory not restored.")
	assert.Equal(t, "/\n", stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 0)
	assert.NoError(t, err)
}

func TestLocal_FormatShell(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})

	cmd := fmt.Sprintf(`%s -c "%s"`, l.ShellExecutable, "uname")
	msg := l.FormatShell("uname")
	t.Logf("cmd = %q", cmd)
	t.Logf("msg = %q", msg)
	assert.Equal(t, msg, `/bin/sh -c "uname"`)

	cmd = fmt.Sprintf(`%s -c "%s"`, l.ShellExecutable, "uname -a")
	msg = l.FormatShell("uname -a")
	t.Logf("cmd = %q", cmd)
	t.Logf("msg = %q", msg)
	assert.Equal(t, msg, `/bin/sh -c "uname -a"`)
}

func TestLocal_TarFailure(t *testing.T) {
	l := run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err := l.Run(
		"/usr/bin/tar",
		"--create",
		"--absolute-names",
		"--file",
		"/dev/null",
		"/xyzzy")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Regexp(
		t,
		regexp.MustCompile(`.*tar: /xyzzy: Cannot stat: No such file or directory.*`),
		stderr)
	assert.NotZero(t, code)
	assert.NoError(t, err)
}
