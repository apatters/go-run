package run_test

import (
	"fmt"
	"os"

	"github.com/apatters/go-run"
)

func ExampleRun() {
	// There is no need for a constructor when running local
	// commands when using the standard runner.
	fmt.Println("Run ls command.")
	stdout, stderr, code, err := run.Run(
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

	fmt.Println("Run id command.")
	stdout, stderr, code, err = run.Run("/bin/id", "root")
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
	// Run id command.
	// stdout = "uid=0(root) gid=0(root) groups=0(root)\n"
	// stderr = ""
	// exit code = 0
}

func ExampleShell() {
	// There is no need for a constructor when running local
	// commands when using the standard runner.
	fmt.Println("Run ls command using shell.")
	stdout, stderr, code, err := run.Shell("/bin/ls -1 /bin/true /bin/false")
	if err != nil {
		fmt.Printf("Internal error executing ls: %s.\n", err)
		os.Exit(1)
	}
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	// Print out the command to run, then run it.
	fmt.Println("Run id command.")
	stdout, stderr, code, err = run.Shell("/bin/id root | cut -f1 -d ' '")
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
	// Run id command.
	// stdout = "uid=0(root)\n"
	// stderr = ""
	// exit code = 0
}
