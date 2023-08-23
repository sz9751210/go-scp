package main

import (
	"fmt"
	"github.com/trzsz/promptui"
	"go-ssh-util/config"
	"go-ssh-util/file"
	"go-ssh-util/ssh"
	"os"
	"path/filepath"
)

func main() {
	for {
		// Prompt for SCP or 'free' command
		scpPrompt := promptui.Select{
			Label: "Choose an operation:",
			Items: []string{"Copy File/Directory (SCP)", "🖥️ Check Systeminfo", "Quit"},
		}

		scpIndex, _, err := scpPrompt.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Check if the system info loop should be active
		systemInfoActive := false

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
		case 1: // systeminfo
			systemInfoActive = true // Set the flag to true

		case 2:
			return
		}

		if systemInfoActive {
			// Prompt for system info options
			systeminfoPrompt := promptui.Select{
				Label: "Choose an operation:",
				Items: []string{"Check CPU", "Check Memory", "Check Disk", "Check Swap", "Check Network"},
			}

			systeminfoIndex, _, err := systeminfoPrompt.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			switch systeminfoIndex {
			case 0: // cpu
				selectedHost, err := config.ChooseAlias()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					return
				}

				ssh.ExecuteRemoteCommand("lscpu", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
			case 1: // memory
				selectedHost, err := config.ChooseAlias()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					return
				}

				ssh.ExecuteRemoteCommand("free -h", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
			case 2: //disk
				selectedHost, err := config.ChooseAlias()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					return
				}

				ssh.ExecuteRemoteCommand("df -h", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
			case 3: //swap
				selectedHost, err := config.ChooseAlias()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					return
				}

				ssh.ExecuteRemoteCommand("cat /proc/swaps", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
			case 4: //network
				selectedHost, err := config.ChooseAlias()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					return
				}

				ssh.ExecuteRemoteCommand("ip", fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
			}
			systemInfoActive = false // Reset the flag
		}
	}
}
