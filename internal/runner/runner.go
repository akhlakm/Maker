package runner

import (
	"os"
	"os/exec"
)

func Execute(exe string, args []string, cwd string) (bool, error) {
	cmd := exec.Command(exe, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
	cmd.Dir = cwd

    err := cmd.Run()
    if err != nil {
        return false, err
    }

    return true, nil
}
