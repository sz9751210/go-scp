package main

import (
	"fmt"
	"github.com/kevinburke/ssh_config"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type sshHost struct {
	Alias string
	Host  string
	Port  string
	User  string
}

func chooseAlias() (sshHost, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return sshHost{}, fmt.Errorf("user home dir failed: %v", err)
	}
	fmt.Println(home)
	path := filepath.Join(home, ".ssh", "config")
	f, err := os.Open(path)
	if err != nil {
		return sshHost{}, fmt.Errorf("open config [%s] failed: %v", path, err)
	}
	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return sshHost{}, fmt.Errorf("decode config [%s] failed: %v", path, err)
	}
	hosts := []sshHost{}
	fmt.Println(cfg.Hosts)
	for _, host := range cfg.Hosts {
		alias := host.Patterns[0].String()
		if alias == "*" {
			continue
		}
		host := ssh_config.Get(alias, "HostName")
		port := ssh_config.Get(alias, "Port")
		user := ssh_config.Get(alias, "User")
		hosts = append(hosts, sshHost{alias, host, port, user})
	}
	fmt.Println(hosts)
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Alias | cyan }} ({{ .Host | red }})",
		Inactive: "  {{ .Alias | cyan }} ({{ .Host | red }})",
		Selected: "\U0001F336 {{ .Alias | red | cyan }}",
		Details: `
--------- SSH Alias ----------
{{ "Alias:" | faint }}	{{ .Alias }}
{{ "Host:" | faint }}	{{ .Host }}
{{ "Port:" | faint }}	{{ .Port }}
{{ "User:" | faint }}	{{ .User }}`,
	}

	searcher := func(input string, index int) bool {
		h := hosts[index]
		alias := strings.Replace(strings.ToLower(h.Alias), " ", "", -1)
		host := strings.Replace(strings.ToLower(h.Host), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(alias, input) || strings.Contains(host, input)
	}

	prompt := promptui.Select{
		Label:     "SSH Alias",
		Items:     hosts,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return sshHost{}, err
	}
	return hosts[i], nil
}

func main() {
	selectedHost, err := chooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// fmt.Printf("You chose Alias: %s, Host: %s, Port: %s\n", selectedHost.Alias, selectedHost.Host, selectedHost.Port)

	selectedFile, err := chooseFileInteractive()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	remoteDestination := promptForRemoteDestination()
	userInfo := fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host)

	copyFileUsingSCP(selectedFile, userInfo, selectedHost.Port, remoteDestination)
}

func promptForRemoteDestination() string {
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

func copyFileUsingSCP(sourceFile, userInfo, port, remoteDestination string) {
	cmdString := fmt.Sprintf("scp -P %s %s %s:%s", port, sourceFile, userInfo, remoteDestination)

	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func chooseFileInteractive() (string, error) {
	dirPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return chooseFileRecursive(dirPath)
}

func chooseFileRecursive(currentDir string) (string, error) {
	files, err := getFilesAndDirectoriesInDirectory(currentDir)
	if err != nil {
		return "", err
	}

	files = append([]string{".."}, files...)

	prompt := promptui.Select{
		Label: "Choose a file or directory:",
		Items: files,
		Templates: &promptui.SelectTemplates{
			Selected: "{{ . | cyan }}",
		},
	}
	fileIndex, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	fileChoice := files[fileIndex]
	if fileChoice == ".." {
		parentDir := filepath.Dir(currentDir)
		return chooseFileRecursive(parentDir)
	}

	selectedPath := filepath.Join(currentDir, fileChoice)
	if isDirectory(selectedPath) {
		return chooseFileRecursive(selectedPath)
	}

	return selectedPath, nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func getFilesAndDirectoriesInDirectory(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var options []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += string(filepath.Separator)
		}
		options = append(options, name)
	}
	return options, nil
}
