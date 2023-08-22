package main

import (
	"fmt"
	"go-copy-tool/config"
	"go-copy-tool/file"
	"go-copy-tool/ssh"
	"os"
	"path/filepath"

	"github.com/trzsz/promptui"
)

func main() {
	// Prompt for SCP or 'free' command
	scpPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"Copy File/Directory (SCP)", "Execute 'free' Command"},
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
	}
}
