package main

import (
	"fmt"
	"go-copy-tool/config"
	"go-copy-tool/file"
	"go-copy-tool/ssh"
	"os"
)

func main() {
	selectedHost, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// fmt.Printf("You chose Alias: %s, Host: %s, Port: %s\n", selectedHost.Alias, selectedHost.Host, selectedHost.Port)

	selectedFile, err := file.ChooseFileInteractive()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	remoteDestination := file.PromptForRemoteDestination()
	userInfo := fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host)

	ssh.CopyFileUsingSCP(selectedFile, userInfo, selectedHost.Port, remoteDestination)
}
