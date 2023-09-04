package menu

import (
	"fmt"
	"go-ssh-util/actions/gcp"

	"github.com/trzsz/promptui"
)

func showGCPMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"GCE", "Return to Main Menu"},
	}

	GCPIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return GCPIndex
}

func RunGCP() {
	GCPActive := true

	for GCPActive {
		GCPIndex := showGCPMenu()

		switch GCPIndex {
		case 0:
			RunGCE()
		default:
			GCPActive = false
		}
	}
}

func showGCEMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"Get VMs", "Start VM", "Stop VM", "Create VM", "Return to Main Menu"},
	}

	GCPIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return GCPIndex
}

func RunGCE() {
	GCPActive := true

	for GCPActive {
		GCPIndex := showGCEMenu()

		switch GCPIndex {
		case 0:
			gcp.RunGetVMs()
		case 1:
			gcp.RunStartVM()
		case 2:
			gcp.RunStopVM()
		case 3:
			gcp.RunCreateGCEInstance()
		default:
			GCPActive = false
		}
	}
}
