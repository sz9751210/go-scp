package actions

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
)

func RunCheckCPU() {
	selectedHost, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	ssh.ExecuteRemoteCommand("lscpu", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
}

func RunCheckMemory() {
	selectedHost, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	ssh.ExecuteRemoteCommand("free -h", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
}

func RunCheckDisk() {
	selectedHost, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	ssh.ExecuteRemoteCommand("df -h", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
}

func RunCheckSwap() {
	selectedHost, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	ssh.ExecuteRemoteCommand("cat /proc/swaps", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
}

func RunCheckNetwork() {
	selectedHost, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	ssh.ExecuteRemoteCommand("ip", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
}
