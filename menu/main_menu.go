package menu

import (
	"fmt"
	"go-ssh-util/actions"

	"github.com/trzsz/promptui"
)

func showMainMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"Copy File/Directory (SCP)", "🖥️ Check Systeminfo", "Quit"},
	}

	scpIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return scpIndex
}

func RunMainLoop() {
	for {
		scpIndex := showMainMenu()

		switch scpIndex {
		case 0: // Copy File/Directory (SCP)
			actions.RunCopyFiles()
		case 1: // systeminfo
			RunSystemInfo()
		case 2:
			fmt.Println("Goodbye!")
			return
		}

	}
}
