package python

import (
	"context"
	"os/exec"

	stashExec "github.com/stashapp/stash/pkg/exec"
)

type Python string

func (p *Python) Command(ctx context.Context, args []string) *exec.Cmd {
	return stashExec.CommandContext(ctx, string(*p), args...)
}

// New returns a new Python instance at the given path.
func New(path string) Python {
	return Python(path)
}

// Resolve tries to find the python executable in the system.
// It first checks for python3, then python.
// Returns an empty string and an exec.ErrNotFound error if not found.
func Resolve() (Python, error) {
	_, err := exec.LookPath("python3")

	if err != nil {
		_, err = exec.LookPath("python")
		if err != nil {
			return "", err
		}
		return "python", nil
	}
	return "python3", nil
}

// IsPythonCommand returns true if arg is "python" or "python3"
func IsPythonCommand(arg string) bool {
	return arg == "python" || arg == "python3"
}
