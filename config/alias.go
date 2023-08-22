package config

import (
	"fmt"
	"github.com/kevinburke/ssh_config"
	"github.com/manifoldco/promptui"
	"go-ssh-util/types"
	"os"
	"path/filepath"
	"strings"
)

func ChooseAlias() (types.SSHhost, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return types.SSHhost{}, fmt.Errorf("user home dir failed: %v", err)
	}
	path := filepath.Join(home, ".ssh", "config")
	f, err := os.Open(path)
	if err != nil {
		return types.SSHhost{}, fmt.Errorf("open config [%s] failed: %v", path, err)
	}
	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return types.SSHhost{}, fmt.Errorf("decode config [%s] failed: %v", path, err)
	}
	hosts := []types.SSHhost{}
	for _, host := range cfg.Hosts {
		alias := host.Patterns[0].String()
		if alias == "*" {
			continue
		}
		host := ssh_config.Get(alias, "HostName")
		port := ssh_config.Get(alias, "Port")
		user := ssh_config.Get(alias, "User")
		hosts = append(hosts, types.SSHhost{alias, host, port, user})
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F4BB {{ .Alias | cyan }} ({{ .Host | red }})",
		Inactive: "  {{ .Alias | cyan }} ({{ .Host | red }})",
		Selected: "\U0001F4BB {{ .Alias | red | cyan }}",
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
	options := append(hosts, types.SSHhost{Alias: "Enter Manually"})
	prompt := promptui.Select{
		Label:     "SSH Alias",
		Items:     options,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return types.SSHhost{}, err
	}
	if i == len(options)-1 {
		// Last option is "Enter Manually"
		return EnterManualSSHHost()
	}
	return hosts[i], nil
}

func EnterManualSSHHost() (types.SSHhost, error) {
	fmt.Println("The SSH host alias is not found in the config. Please enter the host information manually.")

	// Prompt for manual entry
	prompt := []*promptui.Prompt{
		{
			Label: "Enter alias:",
		},
		{
			Label: "Enter host:",
		},
		{
			Label: "Enter port:",
		},
		{
			Label: "Enter user:",
		},
	}

	manualEntry := types.SSHhost{}
	for i, p := range prompt {
		result, err := p.Run()
		if err != nil {
			return types.SSHhost{}, err
		}

		switch i {
		case 0:
			manualEntry.Alias = result
		case 1:
			manualEntry.Host = result
		case 2:
			manualEntry.Port = result
		case 3:
			manualEntry.User = result
		}
	}

	return manualEntry, nil
}
