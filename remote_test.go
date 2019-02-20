// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package run_test

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/apatters/go-run"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemote_RunSuccess(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Run("/bin/true")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.NoError(t, err)
}

func TestRemote_RunFail(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Run("/bin/false")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NotZero(t, code)
	assert.NoError(t, err)
}

func TestRemote_RunOutput(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Run(
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

func TestRemote_RunStdin(t *testing.T) {
	stdinStr := "Hello, world"
	r, err := run.NewRemote(run.RemoteConfig{
		Stdin: strings.NewReader(stdinStr),
	})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Run(
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

func TestRemote_RunStdout(t *testing.T) {
	var b bytes.Buffer
	r, err := run.NewRemote(run.RemoteConfig{
		Stdout: bufio.NewWriter(&b),
	})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Run(
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

func TestRemote_RunStderr(t *testing.T) {
	var b bytes.Buffer
	r, err := run.NewRemote(run.RemoteConfig{
		Stderr: bufio.NewWriter(&b),
	})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Run(
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

func TestRemote_ShellSuccess(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Shell("exit 0")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.NoError(t, err)
}

func TestRemote_ShellFail(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Shell("exit 1")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 1)
	assert.NoError(t, err)
}

func TestRemote_ShellOutput(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)
	stdout, stderr, code, err := r.Shell("cd /bin && ls true false xyzzy")
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

func TestRemote_NewShell(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)
	r.ShellExecutable = "/bin/bash"
	stdout, stderr, code, err := r.Shell("xyzzy")
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

func TestRemote_FormatRun(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)

	msg := r.FormatRun("uname")
	t.Logf("cmd = %q", "uname")
	t.Logf("msg = %q", msg)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* uname`),
		msg)

	msg = r.FormatRun("uname", "-a")
	t.Logf("cmd = %q", "uname -a")
	t.Logf("msg = %q", msg)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* uname -a`),
		msg)
}

func TestRemote_FormatShell(t *testing.T) {
	r, err := run.NewRemote(run.RemoteConfig{})
	require.NoError(t, err)

	cmd := fmt.Sprintf(`%s -c "%s"`, r.ShellExecutable, "uname")
	msg := r.FormatShell("uname")
	t.Logf("cmd = %q", cmd)
	t.Logf("msg = %q", msg)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "uname"`),
		msg)

	cmd = fmt.Sprintf(`%s -c "%s"`, r.ShellExecutable, "uname -a")
	msg = r.FormatShell("uname -a")
	t.Logf("cmd = %q", cmd)
	t.Logf("msg = %q", msg)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "uname -a"`),
		msg)
}
