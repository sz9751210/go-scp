package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"path/filepath"
)

type FileCopyConfig struct {
	SourceFile  string
	Username    string
	RemoteHost  string
	RemotePort  string
	Destination string
}

func main() {
	selectedFile, err := chooseFileInteractive()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	config, err := getCopyConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	destination := fmt.Sprintf("%s@%s:%s", config.Username, config.RemoteHost, config.Destination)
	if config.RemotePort != "" {
		destination = fmt.Sprintf("-P %s %s", config.RemotePort, destination)
	}
	copyFileUsingSCP(selectedFile, destination, config)
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

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func getCopyConfig() (FileCopyConfig, error) {
	var config FileCopyConfig

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Enter your username:",
		Templates: templates,
	}
	config.Username, _ = prompt.Run()

	prompt = promptui.Prompt{
		Label:     "Enter the remote host:",
		Templates: templates,
	}
	config.RemoteHost, _ = prompt.Run()

	prompt = promptui.Prompt{
		Label:     "Enter the remote port (leave empty for default):",
		Templates: templates,
	}
	config.RemotePort, _ = prompt.Run()

	prompt = promptui.Prompt{
		Label:     "Enter the destination path:",
		Templates: templates,
	}
	config.Destination, _ = prompt.Run()

	return config, nil
}

func getFilesInDirectory(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil
}

func copyFileUsingSCP(sourceFile, destination string, config FileCopyConfig) {
	cmdString := fmt.Sprintf("scp -P %s %s %s@%s:%s", config.RemotePort, sourceFile, config.Username, config.RemoteHost, config.Destination)

	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
