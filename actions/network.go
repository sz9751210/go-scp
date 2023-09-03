package actions

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
)

func RunPing() {
	selectedHost, _, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	ssh.ExecuteRemoteCommand("ping 8.8.8.8", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
}
