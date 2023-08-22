package main

import (
	"fmt"
	"go-copy-tool/config"
	"go-copy-tool/file"
	"go-copy-tool/ssh"
	"os"
	"path/filepath"
)

func main() {
	selectedHost, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// fmt.Printf("You chose Alias: %s, Host: %s, Port: %s\n", selectedHost.Alias, selectedHost.Host, selectedHost.Port)

	selectedFile, isDirectoryMode, err := file.ChooseFileInteractive()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	userInfo := fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host)
	remoteDest := file.PromptForRemoteDestination()

	if isDirectoryMode {
		remoteDest = filepath.Join(remoteDest, filepath.Base(selectedFile))
		ssh.CopyUsingSCP(selectedFile, remoteDest, userInfo, selectedHost.Port, true) // Recursive copy
	} else {
		ssh.CopyUsingSCP(selectedFile, remoteDest, userInfo, selectedHost.Port, false) // Single file copy
	}
}
