// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package run_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/apatters/go-run"
	"github.com/stretchr/testify/assert"
)

func TestStdRunner_RunSuccess(t *testing.T) {
	stdout, stderr, code, err := run.Run("/bin/true")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.NoError(t, err)
}

func TestStdRunner_RunFail(t *testing.T) {
	stdout, stderr, code, err := run.Run("/bin/false")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NotZero(t, code)
	assert.NoError(t, err)
}

func TestStdRunner_RunExit(t *testing.T) {
	stdout, stderr, code, err := run.Run("/bin/sh", "-c", "exit 6")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 6)
	assert.NoError(t, err)
}

func TestStdRunner_RunOutput(t *testing.T) {
	stdout, stderr, code, err := run.Run(
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

func TestStdRunner_FormatRun(t *testing.T) {
	msg := run.FormatRun("uname", "-a")
	t.Logf("cmd = %q", "uname -a")
	t.Logf("msg = %q", msg)

	assert.Equal(t, msg, "uname -a")
}

func TestStdRunner_ShellSuccess(t *testing.T) {
	stdout, stderr, code, err := run.Shell("exit 0")
	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.NoError(t, err)
}

func TestStdRunner_ShellFail(t *testing.T) {
	stdout, stderr, code, err := run.Shell("exit 1")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 1)
	assert.NoError(t, err)
}

func TestStdRunner_ShellOutput(t *testing.T) {
	stdout, stderr, code, err := run.Shell("cd /bin && ls true false xyzzy")
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

func TestStdRunner_FormatShell(t *testing.T) {
	cmd := fmt.Sprintf(`%s -c "%s"`, run.DefaultShellExecutable, "uname -a")
	msg := run.FormatShell("uname -a")
	t.Logf("cmd = %q", cmd)
	t.Logf("msg = %q", msg)

	assert.Equal(t, msg, `/bin/sh -c "uname -a"`)
}
