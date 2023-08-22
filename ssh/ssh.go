package ssh

import (
	"fmt"
	"os"
	"os/exec"
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
