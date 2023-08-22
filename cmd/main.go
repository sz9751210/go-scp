package main

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/file"
	"go-ssh-util/ssh"
	"os"
	"path/filepath"

	"github.com/trzsz/promptui"
)

func main() {
	// Prompt for SCP or 'free' command
	scpPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"Copy File/Directory (SCP)", "Check Memory", "Check Disk", "Check Swap"},
	}

	scpIndex, _, err := scpPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch scpIndex {
	case 0: // Copy File/Directory (SCP)
		selectedHost, err := config.ChooseAlias()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		selectedFile, isDirectoryMode, err := file.ChooseFileInteractive()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		remoteDest := file.PromptForRemoteDestination()

		userInfo := fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host)
		if isDirectoryMode {
			remoteDest = filepath.Join(remoteDest, filepath.Base(selectedFile))
			ssh.CopyUsingSCP(selectedFile, remoteDest, userInfo, selectedHost.Port, true) // Recursive copy
		} else {
			ssh.CopyUsingSCP(selectedFile, remoteDest, userInfo, selectedHost.Port, false) // Single file copy
		}

	case 1: // Execute 'free' Command
		selectedHost, err := config.ChooseAlias()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		ssh.ExecuteRemoteCommand("free -h", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	case 2: // Execute 'df' Command
		selectedHost, err := config.ChooseAlias()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		ssh.ExecuteRemoteCommand("df -h", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	case 3: // Execute 'df' Command
		selectedHost, err := config.ChooseAlias()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		ssh.ExecuteRemoteCommand("cat /proc/swaps", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}
}
