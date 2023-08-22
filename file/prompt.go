package file

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

func PromptForRemoteDestination() string {
	prompt := promptui.Prompt{
		Label: "Enter the remote destination path:",
	}
	remoteDest, err := prompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return remoteDest
}
