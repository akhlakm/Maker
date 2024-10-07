package runner

import (
	"os"
	"os/exec"
)

func Execute(exe string, args []string, cwd string) (bool, error) {
	cmd := exec.Command(exe, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
	cmd.Dir = cwd

    // Start the command, will not wait for it to finish.
    err := cmd.Start()

    if err != nil {
        return false, err
    }

    return true, nil
}
