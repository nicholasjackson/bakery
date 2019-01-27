package bakery

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

// General backup and restore functions
type General struct {
}

func (g *General) stopContainer(containerName string) {
	fmt.Printf("Stopping container %s\n", containerName)

	// check if container is running
	output := bytes.NewBuffer([]byte{})
	cmd := exec.Command("lxc", "info", containerName)
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Run()

	if err != nil {
		fmt.Println(output.String())
		os.Exit(1)
	}

	running := regexp.MustCompile("Status: Running")
	if !running.MatchString(output.String()) {
		fmt.Println("Container already stopped")
		fmt.Println("")
		return
	}

	g.executeCommand("lxc", "stop", containerName, "--force")
	fmt.Println("")
}

func (g *General) startContainer(containerName string) {
	fmt.Printf("Starting container %s\n", containerName)
	g.executeCommand("lxc", "start", containerName)
}

func (g *General) mountBackupPath(backupLocation, containerName string) {
	fmt.Printf("Mounting backup path %s into container %s\n", backupLocation, containerName)

	// check if backup path mounted
	output := bytes.NewBuffer([]byte{})
	cmd := exec.Command("lxc", "config", "device", "show", containerName)
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Run()

	if err != nil {
		fmt.Println(output.String())
		os.Exit(1)
	}

	validMount := regexp.MustCompile("source: " + backupLocation)
	if validMount.MatchString(output.String()) {
		fmt.Println("Mount path already exists")
		fmt.Println("")
		return
	}

	g.executeCommand("lxc", "config", "device", "add", containerName, "lxd-conf", "disk", "source="+backupLocation, "path=/mnt/lxd_conf", "readonly=false")
	fmt.Println("")
}

func (g *General) executeCommand(command string, args ...string) error {
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
