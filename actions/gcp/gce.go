package gcp

import (
	"bufio"
	"context"
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/trzsz/promptui"
	"google.golang.org/api/compute/v1"
)

var (
	projectID = "splendid-window-398208"
	ctx       = context.Background()
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
	ImageInfo   ImageInfo
	Subnet      string
	DiskType    string
	DiskSize    string
	// Add other configuration fields as needed
}

type MachineTypeInfo struct {
	Name   string
	Zone   string
	CPU    int64
	Memory string
}

// 使用前請先驗證 -> https://matthung0807.blogspot.com/2023/02/gcp-setup-local-user-credential-to-adc.html
func RunGetVMs() {

	// 創建一個 Compute Engine 服務實例
	computeService, err := compute.NewService(ctx)
	if err != nil {
		fmt.Printf("Error creating compute service: %v\n", err)
		return
	}

	// 列出項目中的所有 GCE 實例
	instanceList, err := computeService.Instances.AggregatedList(projectID).Do()
	if err != nil {
		fmt.Printf("Error listing instances: %v\n", err)
		return
	}

	// 輸出所有 GCE 實例的名稱
	for _, itemList := range instanceList.Items {
		for _, instance := range itemList.Instances {
			zone := path.Base(instance.Zone)
			machine_type := path.Base(instance.MachineType)
			fmt.Printf("Instance Name: %s, Zone: %s, Status: %s,  Machine-Type: %s\n", instance.Name, zone, instance.Status, machine_type)
		}
	}
	// selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// 	return
	// }
	// command := "gcloud compute instances list"
	// if ExecutionMode == 1 {
	// 	ssh.ExecuteLocalCommand(command)
	// } else {
	// 	ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	// }

}

func RunStartVM() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
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
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
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

	// Create a Compute Engine service instance using the service account key file.
	computeService, err := compute.NewService(ctx)
	if err != nil {
		fmt.Printf("Error creating compute service: %v\n", err)
		return GCEInstance{}, err
	}

	// List all GCE instances in the project.
	instanceList, err := computeService.Instances.AggregatedList(projectID).Do()
	if err != nil {
		fmt.Printf("Error listing instances: %v\n", err)
		return GCEInstance{}, err
	}

	// Store instances in a slice of GCEInstance.
	var instances []GCEInstance
	for _, itemList := range instanceList.Items {
		for _, instance := range itemList.Instances {
			zone := path.Base(instance.Zone)
			machineType := path.Base(instance.MachineType)
			internalIP := instance.NetworkInterfaces[0].NetworkIP
			instances = append(instances, GCEInstance{Name: instance.Name, Zone: zone, InternalIP: internalIP, MachineType: machineType, Status: instance.Status})
		}
	}

	// cmd := exec.Command("gcloud", "compute", "instances", "list")

	// // Capture the output of the command
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return GCEInstance{}, fmt.Errorf("Error: %v\n", err)
	// }

	// // Convert the output to a string
	// outputStr := string(output)
	// // fmt.Println(outputStr)

	// // Parse the output to extract instance details
	// instances := parseGCEInstances(outputStr)

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
		return GCEInstance{}, fmt.Errorf("Error: %v\n", err)
	}

	// Get the selected instance by index
	return instances[index], nil
}

// Function to parse the output of 'gcloud compute instances list'
// func parseGCEInstances(output string) []GCEInstance {
// 	lines := strings.Split(output, "\n")
// 	var instances []GCEInstance

// 	// Skip the header line
// 	if len(lines) >= 1 {
// 		for _, line := range lines[1:] {
// 			fields := strings.Fields(line)
// 			// fmt.Println(fields)
// 			// fmt.Println(len(fields))
// 			if len(fields) >= 5 { // Make sure there are at least 6 fields
// 				instance := GCEInstance{
// 					Name:        fields[0],
// 					Zone:        fields[1],
// 					MachineType: fields[2],
// 					InternalIP:  fields[3],
// 				}

// 				// Check if there is an external IP field
// 				if len(fields) == 5 {
// 					instance.Status = fields[4]
// 				}

// 				// Check if there is a status field
// 				if len(fields) == 6 {
// 					instance.ExternalIP = fields[4]
// 					instance.Status = fields[5]
// 				}

// 				instances = append(instances, instance)
// 			}
// 		}
// 	}

// 	return instances
// }

func RunCreateGCEInstance() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// Retrieve the list of available zones
	regions, err := getAvailableRegions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Allow the user to choose a zone
	selectedRegion, err := chooseRegion(regions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	zones, err := getAvailableZones(selectedRegion)
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
	selectedSeries, err := chooseMachineTypeSeries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Allow the user to choose a machine type group
	selectedGroup, err := chooseMachineTypeGroup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// List machine types for the selected group and zone
	machineTypes, err := listMachineTypes(selectedZone, selectedSeries, selectedGroup)
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
	selectedImage, err := chooseImage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	selectedNetwork := chooseNetwork()
	selectedSubnet, err := chooseSubnetwork(selectedNetwork, selectedRegion)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	selectedDiskType, selectedDiskSize := chooseDiskType(selectedZone)

	// Prompt the user for GCE instance configuration
	config, err := promptForGCEInstanceConfig(selectedZone, selectedMachineType, selectedImage, selectedSubnet, selectedDiskType, selectedDiskSize)
	fmt.Println(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Create the command using strings.Builder
	var cmdBuilder strings.Builder
	cmdBuilder.WriteString("gcloud compute instances create ")
	cmdBuilder.WriteString(config.Name)
	cmdBuilder.WriteString(" --zone=")
	cmdBuilder.WriteString(config.Zone)
	cmdBuilder.WriteString(" --machine-type=")
	cmdBuilder.WriteString(config.MachineType)
	cmdBuilder.WriteString(" --image-project=")
	cmdBuilder.WriteString(config.ImageInfo.Project)
	cmdBuilder.WriteString(" --image=")
	cmdBuilder.WriteString(config.ImageInfo.Name)
	cmdBuilder.WriteString(" --subnet=")
	cmdBuilder.WriteString(config.Subnet)
	cmdBuilder.WriteString(" --boot-disk-device-name=")
	cmdBuilder.WriteString(config.Name)
	cmdBuilder.WriteString(" --boot-disk-type=")
	cmdBuilder.WriteString(config.DiskType)
	cmdBuilder.WriteString(" --boot-disk-size=")
	cmdBuilder.WriteString(config.DiskSize)

	// Get the final command string
	command := cmdBuilder.String()

	// Construct the 'gcloud' command to create the GCE instance
	// command := fmt.Sprintf("gcloud compute instances create %s --zone=%s --machine-type=%s --image-project=%s --image=%s", config.Name, config.Zone, config.MachineType, config.ImageInfo.Project, config.ImageInfo.Family)
	// Add other flags and parameters as needed

	// Execute the 'gcloud' command
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

	fmt.Println("GCE instance created successfully.")
}

func chooseMachineTypeSeries() (string, error) {
	series := []struct {
		Label       string
		Description string
		CPU         string
		Memory      string
		Platform    string
	}{
		{"c3", "具備穩定高效能", "4-176", "8 - 1408 GB", "Intel Sapphire Rapids"},
		{"e2", "低成本，適合日常運算", "0.25-32", "1 - 128 GB", "按照供應情形顯示"},
		{"n2", "兼顧價格和效能", "2 - 128", "2 - 864 GB", "Intel Cascade 和 Ice Lake"},
		{"n2d", "兼顧價格和效能", "2 - 224", "2 - 896 GB", "AMD EPYC"},
		{"t2a", "向外擴充工作負載", "1 - 48", "4 - 192 GB", "Ampere Altra Arm"},
		{"t2d", "向外擴充工作負載", "1 - 60", "4 - 240 GB", "AMD EPYC Milan"},
		{"n1", "兼顧價格和效能", "0.25 - 96", "0.6 - 624 GB", "Intel Skylake"},
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F4BB {{ .Label | cyan }} ({{ .Description | red }})",
		Inactive: "  {{ .Label | cyan }} ({{ .Description | red }})",
		Selected: "\U0001F4BB {{ .Label | red | cyan }}",
		Details: `
	--------- Detail ----------
	{{ "Series:" | faint }}	{{ .Label }}
	{{ "Description:" | faint }}	{{ .Description }}
	{{ "CPU:" | faint }}	{{ .CPU }}
	{{ "Memory:" | faint }}	{{ .Memory }}
	{{ "Platform:" | faint }}	{{ .Platform }}`,
	}

	prompt := promptui.Select{
		Label:     "Select a machine type series:",
		Items:     series,
		Templates: templates,
		Size:      10,
	}

	// Show the prompt and get the selected option
	selectedIndex, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	// Get the label of the selected option
	selectedLabel := series[selectedIndex].Label
	fmt.Println(selectedLabel)
	return selectedLabel, nil
}

// Function to allow the user to choose a machine type group
func chooseMachineTypeGroup() (string, error) {
	prompt := promptui.Select{
		Label: "Select a machine type group:",
		Items: []string{"standard", "highcpu", "highmem", "gpu"},
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
func listMachineTypes(zone, series, group string) ([]MachineTypeInfo, error) {
	// fmt.Printf("%s-%s", series, group)
	// // Construct the 'gcloud' command to list machine types with the specified zone filter
	// cmd := exec.Command("gcloud", "compute", "machine-types", "list", "--filter=zone:"+zone)
	// cmd.Stderr = os.Stderr
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	fmt.Println("Error creating stdout pipe:", err)
	// 	return nil, err
	// }

	// if err := cmd.Start(); err != nil {
	// 	fmt.Println("Error starting command:", err)
	// 	return nil, err
	// }
	// var machineTypes []MachineTypeInfo
	// scanner := bufio.NewScanner(stdout)

	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	if strings.Contains(line, fmt.Sprintf("%s-%s", series, group)) {
	// 		fields := strings.Fields(line)
	// 		if len(fields) >= 4 {
	// 			machineType := MachineTypeInfo{
	// 				Name:   fields[0],
	// 				Zone:   fields[1],
	// 				CPU:    fields[2],
	// 				Memory: fields[3],
	// 			}
	// 			machineTypes = append(machineTypes, machineType)
	// 		}
	// 	}
	// }

	// if err := cmd.Wait(); err != nil {
	// 	fmt.Println("Error waiting for command to finish:", err)
	// 	return nil, err
	// }

	// if len(machineTypes) == 0 {
	// 	fmt.Println("No 'n2-highcpu' machine types found in the specified zone.")
	// 	return nil, err
	// }

	// Create a Compute Engine service instance using the service account key file.
	computeService, err := compute.NewService(ctx)
	if err != nil {
		fmt.Printf("Error creating compute service: %v\n", err)
		return []MachineTypeInfo{}, err
	}

	// List all machine types in the specified zone of the project.
	machineTypeList, err := computeService.MachineTypes.List(projectID, zone).Do()
	if err != nil {
		fmt.Printf("Error listing machine types: %v\n", err)
		return []MachineTypeInfo{}, err
	}

	var machineTypes []MachineTypeInfo
	// Output the names and descriptions of the n2 machine types.
	for _, machineType := range machineTypeList.Items {
		if strings.HasPrefix(machineType.Name, fmt.Sprintf("%s-%s", series, group)) {
			cpu := machineType.GuestCpus
			memoryMb := machineType.MemoryMb
			memoryGb := float64(memoryMb) / 1024 // Convert memory from MB to GB.
			memoryGbFormatted := fmt.Sprintf("%.2f", memoryGb)
			machineTypes = append(machineTypes, MachineTypeInfo{Name: machineType.Name, Zone: machineType.Zone, CPU: cpu, Memory: memoryGbFormatted})
		}
	}

	return machineTypes, nil
}

// Function to allow the user to choose a machine type
func chooseMachineType(machineTypes []MachineTypeInfo) (MachineTypeInfo, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ .Name | cyan }} (Zone: {{ .Zone }}, CPU: {{ .CPU }}, Memory: {{ .Memory }})",
		Inactive: "  {{ .Name | cyan }} (Zone: {{ .Zone }}, CPU: {{ .CPU }}, Memory: {{ .Memory }})",
		Selected: "\U0001F449 {{ .Name | red | cyan }} (Zone: {{ .Zone | red }}, CPU: {{ .CPU | red }}, Memory: {{ .Memory | red }})",
		Details: `
	--------- Detail ----------
	{{ "Name:" | faint }}	{{ .Name }}
	{{ "Zone:" | faint }}	{{ .Zone }}
	{{ "CPU:" | faint }}	{{ .CPU }}
	{{ "Memory:" | faint }}	{{ .Memory }}`,
	}

	prompt := promptui.Select{
		Label:     "Select a Machine Type",
		Items:     machineTypes,
		Templates: templates,
		Size:      10,
	}

	// Show the prompt and get the selected machine type
	index, _, err := prompt.Run()
	if err != nil {
		return MachineTypeInfo{}, err
	}

	selectedMachineType := machineTypes[index]

	return selectedMachineType, nil
}

// Function to get the list of available zones
func getAvailableZones(region string) ([]string, error) {
	// cmd := exec.Command("gcloud", "compute", "zones", "list", fmt.Sprintf("--filter=region:%s", region))

	// // Capture the output of the command
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return nil, err
	// }

	// Create a Compute Engine service instance using the service account key file.
	computeService, err := compute.NewService(ctx)
	if err != nil {
		fmt.Printf("Error creating compute service: %v\n", err)
		return nil, err
	}

	// List all zones in the specified region of the project.
	zoneList, err := computeService.Zones.List(projectID).Filter(fmt.Sprintf("region eq .*%s", region)).Do()
	if err != nil {
		fmt.Printf("Error listing zones: %v\n", err)
		return nil, err
	}

	// Parse the output to extract zone names
	var zones []string
	for _, zone := range zoneList.Items {
		zones = append(zones, zone.Name)
	}
	// lines := strings.Split(string(output), "\n")
	// for i := 1; i < len(lines); i++ {
	// 	line := lines[i]
	// 	if strings.TrimSpace(line) != "" {
	// 		zones = append(zones, strings.Fields(line)[0])
	// 	}
	// }
	// fmt.Println(zones)
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

// Function to get the list of available zones
func getAvailableRegions() ([]string, error) {

	// Create a Compute Engine service instance using the service account key file.
	computeService, err := compute.NewService(ctx)
	if err != nil {
		fmt.Printf("Error creating compute service: %v\n", err)
		return nil, err
	}
	// List all regions in the project.
	regionList, err := computeService.Regions.List(projectID).Do()
	if err != nil {
		fmt.Printf("Error listing regions: %v\n", err)
		return nil, err
	}

	// cmd := exec.Command("gcloud", "compute", "regions", "list")

	// // Capture the output of the command
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return nil, err
	// }

	// Parse the output to extract zone names
	var regions []string
	for _, region := range regionList.Items {
		regions = append(regions, region.Name)
	}

	// lines := strings.Split(string(output), "\n")
	// for i := 1; i < len(lines); i++ {
	// 	line := lines[i]
	// 	if strings.TrimSpace(line) != "" {
	// 		regions = append(regions, strings.Fields(line)[0])
	// 	}
	// }
	// fmt.Println(regions)
	return regions, nil
}

// Function to allow the user to choose a zone
func chooseRegion(regions []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select a region:",
		Items: regions,
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
func promptForGCEInstanceConfig(selectedZone string, selectedMachineType MachineTypeInfo, selectedImage ImageInfo, selectedSubnet, selectedDiskType, selectedDiskSize string) (GCEInstanceConfig, error) {
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

	}
	config.Zone = selectedZone
	config.MachineType = selectedMachineType.Name
	config.ImageInfo = selectedImage
	config.Subnet = selectedSubnet
	config.DiskType = selectedDiskType
	config.DiskSize = selectedDiskSize
	return config, nil
}

type ImageInfo struct {
	Name    string
	Project string
	Family  string
}

func chooseImage() (ImageInfo, error) {

	// List of project IDs that host public images.
	// imageProjectIDs := []string{
	// 	"centos-cloud",
	// 	"coreos-cloud",
	// 	"cos-cloud",
	// 	"debian-cloud",
	// 	"rhel-cloud",
	// 	"suse-cloud",
	// 	"suse-sap-cloud",
	// 	"ubuntu-os-cloud",
	// 	"windows-cloud",
	// 	"windows-sql-cloud",
	// }

	cmd := exec.Command("gcloud", "compute", "images", "list")

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running gcloud command:", err)
		os.Exit(1)
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")

	// Extract project and family information into a slice of ImageInfo structs
	imageInfoList := []ImageInfo{}
	for _, line := range lines[1:] { // Skip the header
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			name := fields[0]
			project := fields[1]
			family := fields[2]
			imageInfo := ImageInfo{Name: name, Project: project, Family: family}
			imageInfoList = append(imageInfoList, imageInfo)
		}
	}

	// Create a prompt to select a PROJECT
	projectPrompt := promptui.Select{
		Label: "Select a PROJECT",
		Items: getUniqueProjects(imageInfoList), // Get unique projects from the struct slice
		Size:  10,
	}

	_, selectedProject, err := projectPrompt.Run() // Execute the prompt to select a project
	fmt.Println(selectedProject)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// selectedProject := strings.TrimSpace(result)

	// Display associated FAMILY values for the selected PROJECT
	selectedFamilies := getFamiliesForProject(imageInfoList, selectedProject)

	var selectedImageInfo ImageInfo
	if len(selectedFamilies) > 0 {
		familyPrompt := promptui.Select{
			Label: "Select a FAMILY for " + selectedProject,
			Items: selectedFamilies, // Display families associated with the selected project
			Size:  10,
		}

		_, selectedFamily, err := familyPrompt.Run() // Let the user select a FAMILY

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Find the corresponding NAME for the selected FAMILY
		selectedName := getNameForFamily(imageInfoList, selectedProject, selectedFamily)

		// Print the selected PROJECT and FAMILY
		selectedImageInfo = ImageInfo{Name: selectedName, Project: selectedProject, Family: selectedFamily}
		// fmt.Println("You selected PROJECT:", selectedProject)
		// fmt.Println("You selected FAMILY:", selectedFamily)
		// fmt.Println("You selected Name:", selectedName)
		return selectedImageInfo, nil
	} else {
		fmt.Println("No FAMILY values found for", selectedProject)
		return selectedImageInfo, err
	}

}

// getUniqueProjects extracts unique PROJECT values from the struct slice
func getUniqueProjects(imageInfoList []ImageInfo) []string {
	projectsMap := make(map[string]bool)
	var projects []string
	for _, imageInfo := range imageInfoList {
		if _, exists := projectsMap[imageInfo.Project]; !exists {
			projectsMap[imageInfo.Project] = true
			projects = append(projects, imageInfo.Project)
		}
	}
	return projects
}

// getFamiliesForProject returns a slice of FAMILY values associated with a selected PROJECT
func getFamiliesForProject(imageInfoList []ImageInfo, selectedProject string) []string {
	var families []string
	for _, imageInfo := range imageInfoList {
		if imageInfo.Project == selectedProject {
			families = append(families, imageInfo.Family)
		}
	}
	return families
}

// getNameForFamily returns the NAME associated with a selected FAMILY for a given PROJECT
func getNameForFamily(imageInfoList []ImageInfo, selectedProject, selectedFamily string) string {
	for _, imageInfo := range imageInfoList {
		if imageInfo.Project == selectedProject && imageInfo.Family == selectedFamily {
			return imageInfo.Name
		}
	}
	return "No matching NAME found"
}

type SubnetInfo struct {
	Name      string
	Region    string
	Network   string
	IPRange   string
	StackType string
}

func chooseNetwork() (network string) {
	cmd := exec.Command("gcloud", "compute", "networks", "list")

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running gcloud command:", err)
		os.Exit(1)
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")
	// Extract network names into a slice of NetworkInfo structs
	networkNames := []string{}
	for _, line := range lines[1:] { // Skip the header
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			name := fields[0]
			networkNames = append(networkNames, name)
		}
	}

	// Create a custom template for displaying network information
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "-> {{ .Name | cyan }}",
		Inactive: "   {{ .Name | faint }}",
		Selected: "\U0001F4E1 {{ .Name | green }}",
	}

	// Create a prompt to select a network
	networkPrompt := promptui.Select{
		Label:     "Select a Network",
		Items:     networkNames,
		Templates: templates,
	}

	_, selectedNetwork, err := networkPrompt.Run() // Execute the prompt to select a network

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("You selected network:%s", selectedNetwork)
	return selectedNetwork
}

func chooseSubnetwork(network, region string) (subnet string, err error) {
	command := fmt.Sprintf("gcloud compute networks subnets list --network=%s --filter=region:%s", network, region)
	// cmd := exec.Command("gcloud", "compute", "networks", "subnets", "list", fmt.Sprintf(" --network=%s --filter=region:%s", network, region))
	// cmd := exec.Command(command)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Convert the output to a string
	outputStr := string(output)

	// Split the output into lines
	lines := strings.Split(string(outputStr), "\n")

	// Extract subnet information into a slice of SubnetInfo structs
	subnetInfoList := []SubnetInfo{}
	for _, line := range lines[1:] { // Skip the header
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			name := fields[0]
			region := fields[1]
			network := fields[2]
			ip_range := fields[3]
			stack_type := fields[4]
			subnetInfo := SubnetInfo{Name: name, Region: region, Network: network, IPRange: ip_range, StackType: stack_type}
			subnetInfoList = append(subnetInfoList, subnetInfo)
		}
	}
	// fmt.Println(subnetInfoList)
	// return subnetInfoList, nil
	// }

	// Create a custom template for displaying network information
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F4E1 {{ .Name | cyan }}",
		Inactive: "   {{ .Name | faint }}",
		Selected: "\U0001F4E1 {{ .Name | green }}",
		Details: `
--------- Detail ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Region:" | faint }}	{{ .Region }}
{{ "Network:" | faint }}	{{ .Network }}
{{ "IPRange:" | faint }}	{{ .IPRange }}
{{ "StackType:" | faint }}	{{ .StackType }}`,
	}

	// Create a prompt to select a network
	subnetPrompt := promptui.Select{
		Label:     "Select a Subnet",
		Items:     subnetInfoList,
		Templates: templates,
	}

	index, _, err := subnetPrompt.Run() // Execute the prompt to select a network

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	selectedSubnetInfo := subnetInfoList[index]
	return fmt.Sprintf("projects/wacare-alpha/regions/%s/subnetworks/%s", region, selectedSubnetInfo.Name), nil
	// return selectedSubnet, nil
}

// DiskTypeInfo represents information about a disk type
type DiskTypeInfo struct {
	Name           string
	Zone           string
	ValidDiskSizes string
}

func chooseDiskType(selectedZone string) (diskType, diskSize string) {
	// Run gcloud compute disk-types list command and capture its output
	command := fmt.Sprintf("gcloud compute disk-types list --filter=zone:%s", selectedZone)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running gcloud command: %v\n", err)
		os.Exit(1)
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")

	// Create a slice to store DiskTypeInfo objects
	var diskTypes []DiskTypeInfo

	// Parse each line and create DiskTypeInfo objects
	for _, line := range lines[1:] {
		if line != "" {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				diskType := DiskTypeInfo{
					Name:           parts[0],
					Zone:           parts[1],
					ValidDiskSizes: parts[2],
				}
				diskTypes = append(diskTypes, diskType)
			}
		}
	}

	// Create a prompt to select a disk type
	prompt := promptui.Select{
		Label: "Select a Disk Type",
		Items: diskTypes,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Name }}",
			Active:   "\U0001F4E1 {{ .Name | cyan  }}",
			Inactive: "   {{ .Name | faint  }}",
			Selected: "\U0001F4E1 {{ .Name | green }}",
			Details: `
--------- Detail ----------
{{ "Disk Type:" | faint }}	{{ .Name }}
{{ "Zone:" | faint }}	{{ .Zone }}
{{ "Valid Disk Sizes:" | faint }}	{{ .ValidDiskSizes }}
`,
		},
	}

	// Prompt the user to select a disk type
	index, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		os.Exit(1)
	}

	// Get the selected disk type
	selectedDiskType := diskTypes[index]
	diskType = selectedDiskType.Name
	// Prompt the user to enter a disk size
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter Disk Size , Range -> %s: ", diskTypes[index].ValidDiskSizes)
	scanner.Scan()
	diskSize = scanner.Text()
	// Print the selected disk type and entered disk size
	fmt.Printf("\nSelected Disk Type:\nDisk Type: %s\nZone: %s\nValid Disk Sizes: %s\nEntered Disk Size: %sGB\n",
		selectedDiskType.Name, selectedDiskType.Zone, selectedDiskType.ValidDiskSizes, diskSize)
	return diskType, diskSize

}
