package ffmpeg

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/logger"
)

type Encoder string

var (
	runningEncoders      = make(map[string][]*os.Process)
	runningEncodersMutex = sync.RWMutex{}
)

func registerRunningEncoder(path string, process *os.Process) {
	runningEncodersMutex.Lock()
	processes := runningEncoders[path]

	runningEncoders[path] = append(processes, process)
	runningEncodersMutex.Unlock()
}

func deregisterRunningEncoder(path string, process *os.Process) {
	runningEncodersMutex.Lock()
	defer runningEncodersMutex.Unlock()
	processes := runningEncoders[path]

	for i, v := range processes {
		if v == process {
			runningEncoders[path] = append(processes[:i], processes[i+1:]...)
			return
		}
	}
}

func waitAndDeregister(path string, cmd *exec.Cmd) error {
	err := cmd.Wait()
	deregisterRunningEncoder(path, cmd.Process)

	return err
}

func KillRunningEncoders(path string) {
	runningEncodersMutex.RLock()
	processes := runningEncoders[path]
	runningEncodersMutex.RUnlock()

	for _, process := range processes {
		// assume it worked, don't check for error
		logger.Infof("Killing encoder process for file: %s", path)
		if err := process.Kill(); err != nil {
			logger.Warnf("failed to kill process %v: %v", process.Pid, err)
		}

		// wait for the process to die before returning
		// don't wait more than a few seconds
		done := make(chan error)
		go func() {
			_, err := process.Wait()
			done <- err
		}()

		select {
		case <-done:
			return
		case <-time.After(5 * time.Second):
			return
		}
	}
}
