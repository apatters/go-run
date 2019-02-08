// This example shows how to use the run.Runner interface to create
// functions to take either Local or Remote objects. In this example,
// we "log" the command run by the Run() and Shell() methods.

package run_test

import (
	"fmt"

	"github.com/apatters/go-run"
)

// logRun outputs the command as if run on the command line and then
// "runs" the command using the Run() method. It can take either a
// run.Local or run.Remote object as they both fulfill the run.Runner
// interface.
func logRun(r run.Runner, cmd string, cmdArgs ...string) (stdout string, stderr string, exitCode int, err error) {
	fmt.Printf("%s\n", r.FormatRun(cmd, cmdArgs...))
	return r.Run(cmd, cmdArgs...)
}

// logShell outputs the shell command as if run on the command line
// and then "runs" the command using the Shell() method. It can take
// either a run.Local or run.Remote object as they both fulfill the
// run.Runner interface.
func logShell(r run.Runner, cmd string) (stdout string, stderr string, exitCode int, err error) {
	fmt.Printf("%s\n", r.FormatShell(cmd))
	return r.Shell(cmd)
}

func ExampleRunner() {
	l := run.NewLocal(run.LocalConfig{})
	r, _ := run.NewRemote(run.RemoteConfig{
		Credentials: run.Credentials{
			Hostname: "localhost",
			Username: "buildman"},
	})

	stdout, stderr, code, _ := logRun(l, "/bin/ls", "-1", "/bin/true", "/bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	stdout, stderr, code, _ = logRun(r, "/bin/ls", "-1", "/bin/true", "/bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	stdout, stderr, code, _ = logShell(l, "/bin/ls -1 /bin/true /bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	stdout, stderr, code, _ = logShell(r, "/bin/ls -1 /bin/true /bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	// Output:
	// /bin/ls -1 /bin/true /bin/false
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = ""
	// exit code = 0
	//
	// ssh buildman@localhost /bin/ls -1 /bin/true /bin/false
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = ""
	// exit code = 0
	//
	// /bin/sh -c "/bin/ls -1 /bin/true /bin/false"
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = ""
	// exit code = 0
	//
	// ssh buildman@localhost /bin/sh -c "/bin/ls -1 /bin/true /bin/false"
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = ""
	// exit code = 0
}
