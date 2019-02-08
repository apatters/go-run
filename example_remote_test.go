package run_test

import (
	"fmt"
	"os"

	"github.com/apatters/go-run"
)

func ExampleRemote_Run() {
	// Initialize Remote object using defaults.
	runner, _ := run.NewRemote(run.RemoteConfig{
		Credentials: run.Credentials{
			Hostname: "localhost"},
	})

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

	fmt.Println("Run the ls command with an internal error (bad user/passwd).")
	runner, _ = run.NewRemote(run.RemoteConfig{
		Credentials: run.Credentials{
			Hostname: "localhost",
			Username: "bad_user",
			Password: "xyzzy"},
	})
	_, _, _, err = runner.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false")
	fmt.Printf("Internal error executing ls: %s.\n", err)
	fmt.Println()

	// Output:
	// Run ls command.
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = ""
	// exit code = 0
	//
	// Run ls command with an expected error.
	// stdout = "/bin/false\n/bin/true\n"
	// stderr = "/bin/ls: cannot access '/xyzzy': No such file or directory\n"
	// exit code = 2
	//
	// Run the ls command with an internal error (bad user/passwd).
	// Internal error executing ls: run: connection to bad_user@localhost failed: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain.
}

func ExampleRemote_Shell() {
	// Initialize Remote object using defaults.
	runner, _ := run.NewRemote(run.RemoteConfig{})

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

	fmt.Println("Run the ls command using a shell with an internal error (bad user/passwd).")
	runner, _ = run.NewRemote(run.RemoteConfig{
		Credentials: run.Credentials{
			Hostname: "localhost",
			Username: "bad_user",
			Password: "xyzzy"},
	})
	_, _, _, err = runner.Shell("/bin/ls -1 /bin/true /bin/false")
	fmt.Printf("Internal error executing ls: %s.\n", err)
	fmt.Println()

	fmt.Println("Run complex shell command.")
	runner, _ = run.NewRemote(run.RemoteConfig{})
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
	// stderr = "/bin/ls: cannot access '/xyzzy': No such file or directory\n"
	// exit code = 2
	//
	// Run the ls command using a shell with an internal error (bad user/passwd).
	// Internal error executing ls: run: connection to bad_user@localhost failed: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain.
	//
	// Run complex shell command.
	// stdout = "false\n"
	// stderr = ""
	// exit code = 0
}
