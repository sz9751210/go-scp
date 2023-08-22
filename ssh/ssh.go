package ssh

import (
	"fmt"
	"os"
	"os/exec"
)

func CopyFileUsingSCP(sourceFile, userInfo, port, remoteDestination string) {
	cmdString := fmt.Sprintf("scp -P %s %s %s:%s", port, sourceFile, userInfo, remoteDestination)

	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
