package file

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"path/filepath"
)

func ChooseFileInteractive() (string, error) {
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

func PromptForRemoteDestination() string {
	prompt := promptui.Prompt{
		Label: "Enter the remote destination path ",
	}
	remoteDest, err := prompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return remoteDest
}
