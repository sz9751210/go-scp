package file

import (
	"fmt"
	"go-copy-tool/types"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/manifoldco/promptui"
)

func ChooseFileInteractive() (string, bool, error) {
	dirPath, err := os.Getwd()
	if err != nil {
		return "", false, err
	}

	return chooseModeAndFile(dirPath)
}

func chooseModeAndFile(currentDir string) (string, bool, error) {
	modePrompt := promptui.Select{
		Label: "Choose SCP mode:",
		Items: []string{"File Mode", "Directory Mode"},
	}
	modeIndex, _, err := modePrompt.Run()
	if err != nil {
		return "", false, err
	}

	isDirectoryMode := modeIndex == 1

	return chooseFileRecursive(currentDir, isDirectoryMode)
}

func chooseFileRecursive(currentDir string, isDirectoryMode bool) (string, bool, error) {
	files, err := getFilesAndDirectoriesInDirectory(currentDir)
	if err != nil {
		return "", false, err
	}

	files = append([]string{".."}, files...)

	var fileDetails []types.FileDetail
	for _, file := range files {
		if file == ".." {
			fileDetails = append(fileDetails, types.FileDetail{Name: ".."})
		} else {
			fullPath := filepath.Join(currentDir, file)
			fileInfo, err := os.Stat(fullPath)
			if err == nil {
				isDir := fileInfo.IsDir()
				if (isDirectoryMode && isDir) || (!isDirectoryMode && !isDir) {
					details := types.FileDetail{
						Name:        file,
						Size:        fileInfo.Size(),
						Mode:        fileInfo.Mode(),
						ModTime:     fileInfo.ModTime(),
						IsDirectory: isDir,
						Owner:       getFileOwner(fileInfo),
						Group:       getFileGroup(fileInfo),
					}
					fileDetails = append(fileDetails, details)
				}
			}
		}
	}

	prompt := promptui.Select{
		Label: "Choose a file or directory:",
		Items: fileDetails,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Name }}?",
			Active:   "\U0001F4C4 {{ .Name | cyan }}",
			Inactive: "  {{ .Name | cyan }}",
			Selected: "\U0001F4C4 {{ .Name | red | cyan }}",
			Details: `
	--------- File Detail ----------
	{{ "Name:" | faint }}       {{ .Name }}
	{{ "Size:" | faint }}       {{ .Size }} bytes
	{{ "Mode:" | faint }}       {{ .Mode }}
	{{ "ModTime:" | faint }}    {{ .ModTime }}
	{{ "IsDirectory:" | faint }} {{ .IsDirectory }}
	{{ "Owner:" | faint }}      {{ .Owner }}
	{{ "Group:" | faint }}      {{ .Group }}
	`,
		},
		Size: 10,
	}

	fileIndex, _, err := prompt.Run()
	if err != nil {
		return "", false, err
	}

	fileChoice := fileDetails[fileIndex]
	if fileChoice.Name == ".." {
		parentDir := filepath.Dir(currentDir)
		return chooseModeAndFile(parentDir)
	}

	selectedPath := filepath.Join(currentDir, fileChoice.Name)

	return selectedPath, fileChoice.IsDirectory, nil
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

func getFileOwner(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)
	user, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		return ""
	}
	return user.Username
}

func getFileGroup(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)
	group, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
	if err != nil {
		return ""
	}
	return group.Name
}
