package ssh

import (
	"fmt"
	"go-ssh-util/types"
	"golang.org/x/crypto/ssh"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func CopyUsingSCP(source, destination, userInfo, port string, isRecursive bool) {
	cmdString := ""
	if isRecursive {
		cmdString = fmt.Sprintf("scp -P %s -r %s %s:%s", port, source, userInfo, destination)
	} else {
		cmdString = fmt.Sprintf("scp -P %s %s %s:%s", port, source, userInfo, destination)
	}

	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func ExecuteRemoteCommand(command, userInfo, port string) {
	sshCmd := exec.Command("ssh", "-p", port, userInfo, command)
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	err := sshCmd.Run()
	if err != nil {
		fmt.Println("Error executing remote command:", err)
	} else {
		fmt.Println("Command executed successfully on remote host")
	}
}

func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func SSHToRemoteHostWithKey(host types.SSHhost) error {

	// 相对路径
	relativePath := host.IdentityFile

	expandedPath, err := expandTilde(relativePath)
	if err != nil {
		return err
	}
	privateKey, err := os.ReadFile(expandedPath)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	address := fmt.Sprintf("%s:%s", host.Host, host.Port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return err
	}
	defer client.Close()

	// create session
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// setting stdin, out, err
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// create terminal
	termModes := ssh.TerminalModes{
		ssh.ECHO:          1,     // control signal
		ssh.TTY_OP_ISPEED: 14400, // input speed
		ssh.TTY_OP_OSPEED: 14400, // output speed
	}
	err = session.RequestPty("xterm", 80, 40, termModes)
	if err != nil {
		return err
	}

	// create Shell
	err = session.Shell()
	if err != nil {
		return err
	}

	// wait session
	session.Wait()
	return nil
}
