package gcp

import (
	"bufio"
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

type GCEInstanceConfig struct {
	Name        string
	Zone        string
	MachineType string
	// Add other configuration fields as needed
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

func RunCreateGCEInstance() {
	selectedHost, ExecutionMode, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// Retrieve the list of available zones
	zones, err := getAvailableZones()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Allow the user to choose a zone
	selectedZone, err := chooseZone(zones)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Retrieve the list of available machine type groups
	// groups, err := getMachineTypeGroups(selectedZone)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// 	return
	// }

	// Allow the user to choose a machine type group
	selectedGroup, err := chooseMachineTypeGroup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// List machine types for the selected group and zone
	machineTypes, err := listMachineTypes(selectedZone, selectedGroup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Allow the user to choose a machine type
	selectedMachineType, err := chooseMachineType(machineTypes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Prompt the user for GCE instance configuration
	config, err := promptForGCEInstanceConfig(selectedZone, selectedMachineType)
	fmt.Println(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Construct the 'gcloud' command to create the GCE instance
	command := fmt.Sprintf("gcloud compute instances create %s --zone=%s --machine-type=%s", config.Name, config.Zone, config.MachineType)
	// Add other flags and parameters as needed

	// Execute the 'gcloud' command
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

	fmt.Println("GCE instance created successfully.")
}

// Function to allow the user to choose a machine type group
func chooseMachineTypeGroup() (string, error) {
	prompt := promptui.Select{
		Label: "Select a machine type group:",
		Items: []string{"standard", "cpu", "mem", "gpu"},
		Size:  10,
	}

	// Show the prompt and get the selected group
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	fmt.Println(result)
	return result, nil
}

// Function to list machine types for the selected group and zone
func listMachineTypes(zone, group string) ([]string, error) {
	// Construct the 'gcloud' command to list machine types with the specified zone filter
	gcloudCmd := exec.Command("gcloud", "compute", "machine-types", "list", fmt.Sprintf("--filter=zone:%s", zone))

	// Create a pipe to capture the output of the 'gcloud' command
	gcloudOutput, err := gcloudCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// Start the 'gcloud' command
	if err := gcloudCmd.Start(); err != nil {
		return nil, err
	}

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(gcloudOutput)

	// Read the output of the 'gcloud' command line by line
	var gcloudOutputLines []string
	for scanner.Scan() {
		gcloudOutputLines = append(gcloudOutputLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Construct the 'grep' command to further filter the output
	grepCmd := exec.Command("grep", group)

	// Set the input of the 'grep' command to the output of the 'gcloud' command
	grepCmd.Stdin = strings.NewReader(strings.Join(gcloudOutputLines, "\n"))

	// Capture the output of the 'grep' command
	grepOutput, err := grepCmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse the output of the 'grep' command to extract machine type names
	var machineTypes []string
	for _, line := range strings.Split(string(grepOutput), "\n") {
		if strings.TrimSpace(line) != "" {
			machineTypes = append(machineTypes, strings.Fields(line)[0])
		}
	}

	return machineTypes, nil
}

// Function to allow the user to choose a machine type
func chooseMachineType(machineTypes []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select a machine type:",
		Items: machineTypes,
		Size:  10,
	}

	// Show the prompt and get the selected machine type
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// Function to get the list of available zones
func getAvailableZones() ([]string, error) {
	cmd := exec.Command("gcloud", "compute", "zones", "list")

	// Capture the output of the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// Parse the output to extract zone names
	var zones []string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) != "" {
			zones = append(zones, strings.Fields(line)[0])
		}
	}
	fmt.Println(zones)
	return zones, nil
}

// Function to allow the user to choose a zone
func chooseZone(zones []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select a zone:",
		Items: zones,
		Size:  10,
	}

	// Show the prompt and get the selected zone
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// Function to prompt the user for GCE instance configuration
func promptForGCEInstanceConfig(selectedZone, selectedMachineType string) (GCEInstanceConfig, error) {
	prompt := []*promptui.Prompt{
		{
			Label: "Enter instance name:",
		},
		// {
		// 	Label:   "Zone:",
		// 	Default: selectedZone, // Set the default zone to the selected one
		// },
		// {
		// 	Label:   "Choose machine type:",
		// 	Default: selectedMachineType,
		// },
		// Add prompts for other configuration fields as needed
	}

	config := GCEInstanceConfig{}

	for i, p := range prompt {
		result, err := p.Run()
		if err != nil {
			return GCEInstanceConfig{}, err
		}
		switch i {
		case 0:
			config.Name = result
			// Set other configuration fields based on prompts
		}
		config.Zone = selectedZone
		config.MachineType = selectedMachineType
	}
	fmt.Println("ok")
	return config, nil
}

// Function to execute a 'gcloud' command
// func execGCloudCommand(command string, ExecutionMode int, selectedHost config.SSHhost) error {
// 	if ExecutionMode == 1 {
// 		return ssh.ExecuteLocalCommand(command)
// 	}

// 	// Execute remotely using SSH
// 	remoteCommand := fmt.Sprintf("ssh -p %s %s@%s \"%s\"", selectedHost.Port, selectedHost.User, selectedHost.Host, command)
// 	return ssh.ExecuteRemoteCommand(remoteCommand)
// }
