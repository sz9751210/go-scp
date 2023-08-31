package actions

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
)

func RunStatus() {
	selectedHost, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	ssh.ExecuteRemoteCommand("docker ps", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
}
