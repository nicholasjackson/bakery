package baker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func stopContainer() {
	fmt.Printf("Stopping container %s\n", *container)
	executeCommand("lxc", "stop", *container, "--force")
	fmt.Println("")
}

func startContainer() {
	fmt.Printf("Starting container %s\n", *container)
	executeCommand("lxc", "start", *container)
}

func mountBackupPath() {
	fmt.Printf("Mounting backup path %s into container %s\n", *backupLocation, *container)

	// check if backup path mounted
	output := bytes.NewBuffer([]byte{})
	cmd := exec.Command("lxc", "config", "device", "show", *container)
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Run()

	if err != nil {
		fmt.Println(output.String())
		os.Exit(1)
	}

	validMount := regexp.MustCompile("source: " + *backupLocation)
	if validMount.MatchString(output.String()) {
		fmt.Println("Mount path already exists")
		fmt.Println("")
		return
	}

	executeCommand("lxc", "config", "device", "add", *container, "lxd-conf", "disk", "source="+*backupLocation, "path=/mnt/lxd_conf")
	fmt.Println("")
}

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

	return nil
}
