//go:build darwin
// +build darwin

package desktop

import (
	"os/exec"

	"github.com/kermieisinthehouse/gosx-notifier"
	"github.com/stashapp/stash/pkg/logger"
)

func isService() bool {
	// MacOS /does/ support services, using launchd, but there is no straightforward way to check if it was used.
	return false
}

func isServerDockerized() bool {
	return false
}

func hideExecShell(cmd *exec.Cmd) {

}

func sendNotification(notificationTitle string, notificationText string) {
	notification := gosxnotifier.NewNotification(notificationText)
	notification.Title = notificationTitle
	notification.AppIcon = getIconPath()
	notification.Link = getServerURL("")
	err := notification.Push()

	if err != nil {
		logger.Errorf("Could not send MacOS notification: %s", err.Error())
	}
}
