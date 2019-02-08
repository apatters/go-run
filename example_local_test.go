package run_test

import (
	"fmt"
	"os"

	"github.com/apatters/go-run"
)

func ExampleLocal_Run() {
	// Initialize Local object using defaults.
	runner := run.NewLocal(run.LocalConfig{})

	fmt.Println("Run ls command.")
	stdout, stderr, code, err := runner.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false")
	if err != nil {
		fmt.Printf("Internal error executing ls: %s.\n", err)
		os.Exit(1)
	}
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run ls command with an expected error.")
	stdout, stderr, code, err = runner.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false",
		"/xyzzy")
	if err != nil {
		fmt.Printf("Internal error executing ls: %s.\n", err)
		os.Exit(1)
	}
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run the ls command with an internal error (bad path).")
	_, _, _, err = runner.Run(
		"/not_a_good_path/ls",
		"-1",
		"/bin/true",
		"/bin/false")
	fmt.Printf("Internal error executing ls: %s.\n", err)
	fmt.Println()

	fmt.Println("Run ls command after changing directory.")
	runner = run.NewLocal(run.LocalConfig{Dir: "/bin"})
	stdout, stderr, code, err = runner.Run(
		"/bin/ls",
		"-1",
		"true",
		"false")
	if err != nil {
		fmt.Printf("Internal error executing ls: %s.\n", err)
		os.Exit(1)
	}
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	// Output:
	// Run ls command.
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = ""
	// exit code = 0
	//
	// Run ls command with an expected error.
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = "/bin/ls: cannot access /xyzzy: No such file or directory\n"
	// exit code = 2
	//
	// Run the ls command with an internal error (bad path).
	// Internal error executing ls: fork/exec /not_a_good_path/ls: no such file or directory.
	//
	// Run ls command after changing directory.
	// stdout = "false\ntrue\n"
	// stderr = ""
	// exit code = 0
}

func ExampleLocal_Shell() {
	// Initialize Local object using defaults.
	runner := run.NewLocal(run.LocalConfig{})

	fmt.Println("Run ls command using shell.")
	stdout, stderr, code, err := runner.Shell("/bin/ls -1 /bin/true /bin/false")
	if err != nil {
		fmt.Printf("Internal error executing ls: %s.\n", err)
		os.Exit(1)
	}
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run ls command using shell with an expected error.")
	stdout, stderr, code, err = runner.Shell("/bin/ls -1 /bin/true /bin/false /xyzzy")
	if err != nil {
		fmt.Printf("Internal error executing ls: %s.\n", err)
		os.Exit(1)
	}
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run the ls command using shell with an internal error (bad path to shell).")
	runner = run.NewLocal(run.LocalConfig{ShellExecutable: "/bin/badsh"})
	_, _, _, err = runner.Shell("/bin/ls -1 /bin/true /bin/false")
	fmt.Printf("Internal error executing ls: %s.\n", err)
	fmt.Println()

	fmt.Println("Run ls command using shell after changing directory.")
	runner = run.NewLocal(run.LocalConfig{Dir: "/bin"})
	stdout, stderr, code, err = runner.Shell("/bin/ls -1 true false")
	if err != nil {
		fmt.Printf("Internal error executing ls: %s.\n", err)
		os.Exit(1)
	}
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run complex shell command.")
	runner = run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, err = runner.Shell("cd /bin && /bin/ls -1 true false | head -n 1")
	if err != nil {
		fmt.Printf("Internal error executing ls: %s.\n", err)
		os.Exit(1)
	}
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	// Output:
	// Run ls command using shell.
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = ""
	// exit code = 0
	//
	// Run ls command using shell with an expected error.
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = "/bin/ls: cannot access /xyzzy: No such file or directory\n"
	// exit code = 2
	//
	// Run the ls command using shell with an internal error (bad path to shell).
	// Internal error executing ls: fork/exec : no such file or directory.
	//
	// Run ls command using shell after changing directory.
	// stdout = "false\ntrue\n"
	// stderr = ""
	// exit code = 0
	//
	// Run complex shell command.
	// stdout = "false\n"
	// stderr = ""
	// exit code = 0
}
