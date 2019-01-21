package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var version = "dev"
var backupName = ""

// flags
var container = flag.String("container", "penguin", "Name of the container to backup")
var snapShotName = flag.String("snapshot-prefix", "backup-snapshot", "Prefix of the snapshot to take, snapshots will have time appended to them")
var backupLocation = flag.String("backup-location", "/mnt/stateful/lxd_conf", "Location to backup container")
var username = flag.String("username", "chronos", "Crostini username e.g. jacksonnic")

func main() {
	flag.Parse()
	backupName = *container + "-backup"

	fmt.Println("Crostini Backup and Restore tool version:", version)
	fmt.Println("Starting backup, WARNING: This operation can take a long time")
	fmt.Println("")

	stopContainer()
	createSnapShot()
	publishBackup()
	exportImage()
	splitBackup()
	startContainer()
	mountBackupPath()
	moveBackupFiles()
	deleteBackupImage()

	fmt.Println("")
}

func stopContainer() {
	fmt.Printf("Stopping container %s\n", *container)
	executeCommand("lxc", "stop", *container, "--force")
	fmt.Println("")
}

func createSnapShot() {
	*snapShotName = fmt.Sprintf("%s-%d", *snapShotName, time.Now().Nanosecond())
	fmt.Printf("Creating snapshot of container:%s name:%s\n", *container, *snapShotName)

	executeCommand("lxc", "snapshot", *container, *snapShotName)
	fmt.Println()
}

func publishBackup() {
	fmt.Printf("Publish container: %s to %s\n", *container, backupName)
	fmt.Println("If the container publish is interupted, your container may be left in a bad state,")
	fmt.Printf("in this instance you can restore the snapshot using the command: lxc restore %s %s\n", *container, *snapShotName)

	executeCommand("lxc", "publish", *container, "--alias", backupName)
	fmt.Println()
}

func exportImage() {
	fmt.Printf("Exporting container to: %s\n", *backupLocation)
	executeCommand("lxc", "image", "export", backupName, *backupLocation+"/"+backupName)
	fmt.Println()
}

func splitBackup() {
	backupFileLocation := *backupLocation + "/" + backupName + ".tar.gz"

	fmt.Println("Splitting backup into 3GB chunks")
	executeCommand("split", "-b", "3GB", backupFileLocation, backupFileLocation+".")
	executeCommand("rm", backupFileLocation)
	fmt.Println()
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

func moveBackupFiles() {
	fmt.Printf("Moving backup files to /home/%s in container %s\n", *username, *container)
	// create the user folder if required
	executeCommand("lxc", "exec", *container, "--", "mkdir", "-p", "/home/"+*username)
	// move the files
	executeCommand("lxc", "exec", *container, "--", "find", "/mnt/lxd_conf", "-name", "*.tar.gz.*", "-exec", "cp", "{}", "-t", "/home/"+*username, ";")
	// delete the backups
	executeCommand("find", *backupLocation, "-name", "*.tar.gz.*", "-exec", "rm", "{}", ";")
	fmt.Println("")
}

func deleteBackupImage() {
	fmt.Printf("Deleting temporary image %s", backupName)
	executeCommand("lxc", "image", "delete", backupName)
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
