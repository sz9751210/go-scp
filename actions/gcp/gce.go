package gcp

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
	"os/exec"
	"strings"

	"github.com/trzsz/promptui"
)

type GCEInstance struct {
	Name        string
	Zone        string
	MachineType string
	InternalIP  string
	ExternalIP  string
	Status      string
}

func RunGetVMs() {
	selectedHost, ExecutionMode, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := "gcloud compute instances list"
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

}

func RunStartVM() {
	selectedHost, ExecutionMode, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	selectedGCE, err := ChooseGCE()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := fmt.Sprintf("gcloud compute instances start %s --zone=%s", selectedGCE.Name, selectedGCE.Zone)
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

}

func RunStopVM() {
	selectedHost, ExecutionMode, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	selectedGCE, err := ChooseGCE()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := fmt.Sprintf("gcloud compute instances stop %s --zone=%s", selectedGCE.Name, selectedGCE.Zone)
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

}

func ChooseGCE() (GCEInstance, error) {
	cmd := exec.Command("gcloud", "compute", "instances", "list")

	// Capture the output of the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return GCEInstance{}, fmt.Errorf("Error:", err)
	}

	// Convert the output to a string
	outputStr := string(output)
	// fmt.Println(outputStr)

	// Parse the output to extract instance details
	instances := parseGCEInstances(outputStr)

	// Create a prompt for selecting an instance
	prompt := promptui.Select{
		Label: "Select an instance:",
		Items: instances,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Name }} ({{ .Status }})",
			Active:   "\U0001F4BB {{ .Name | cyan }} ({{ .Status | red }})",
			Inactive: "  {{ .Name | cyan }} ({{ .Status | red }})",
			Selected: "{{ .Name | red | cyan }}",
			Details: `
	--------- Detail ----------
	{{ "Name:" | faint }}	{{ .Name }}
	{{ "Type:" | faint }}	{{ .MachineType }}
	{{ "Zone:" | faint }}	{{ .Zone }}
	{{ "IP:" | faint }}	{{ .InternalIP }}
	{{ "Status:" | faint }}	{{ .Status }}`,
		},
		Size: 10,
	}

	// Show the prompt and get the selected instance
	index, _, err := prompt.Run()
	if err != nil {
		return GCEInstance{}, fmt.Errorf("Error:", err)
	}

	// Get the selected instance by index
	return instances[index], nil
}

// Function to parse the output of 'gcloud compute instances list'
func parseGCEInstances(output string) []GCEInstance {
	lines := strings.Split(output, "\n")
	var instances []GCEInstance

	// Skip the header line
	if len(lines) >= 1 {
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			// fmt.Println(fields)
			// fmt.Println(len(fields))
			if len(fields) >= 5 { // Make sure there are at least 6 fields
				instance := GCEInstance{
					Name:        fields[0],
					Zone:        fields[1],
					MachineType: fields[2],
					InternalIP:  fields[3],
				}

				// Check if there is an external IP field
				if len(fields) == 5 {
					instance.Status = fields[4]
				}

				// Check if there is a status field
				if len(fields) == 6 {
					instance.ExternalIP = fields[4]
					instance.Status = fields[5]
				}

				instances = append(instances, instance)
			}
		}
	}

	return instances
}
