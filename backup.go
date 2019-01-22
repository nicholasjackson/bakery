package baker

import (
	"fmt"
	"time"
)

type Backup struct {
}

func (b *Backup) Execute() {
	backupName = *container + "-backup"

	fmt.Println("")
	b.createSnapShot()
	b.publishBackup()
	b.exportImage()
	b.splitBackup()
	b.startContainer()
	b.mountBackupPath()
	b.moveBackupFiles()
	b.deleteBackupImage()

	fmt.Println("")
}

func (b *Backup) createSnapShot() {
	*snapShotName = fmt.Sprintf("%s-%d", *snapShotName, time.Now().Nanosecond())
	fmt.Printf("Creating snapshot of container:%s name:%s\n", *container, *snapShotName)

	executeCommand("lxc", "snapshot", *container, *snapShotName)
	fmt.Println()
}

func (b *Backup) publishBackup() {
	fmt.Printf("Publish container: %s to %s\n", *container, backupName)
	fmt.Println("If the container publish is interupted, your container may be left in a bad state,")
	fmt.Printf("in this instance you can restore the snapshot using the command: lxc restore %s %s\n", *container, *snapShotName)

	executeCommand("lxc", "publish", *container, "--alias", backupName)
	fmt.Println()
}

func (b *Backup) exportImage() {
	fmt.Printf("Exporting container to: %s\n", *backupLocation)
	executeCommand("lxc", "image", "export", backupName, *backupLocation+"/"+backupName)
	fmt.Println()
}

func (b *Backup) splitBackup() {
	backupFileLocation := *backupLocation + "/" + backupName + ".tar.gz"

	fmt.Println("Splitting backup into 3GB chunks")
	executeCommand("split", "-b", "3GB", backupFileLocation, backupFileLocation+".")
	executeCommand("rm", backupFileLocation)
	fmt.Println()
}

func (b *Backup) moveBackupFiles() {
	fmt.Printf("Moving backup files to /home/%s in container %s\n", *username, *container)
	// create the user folder if required
	executeCommand("lxc", "exec", *container, "--", "mkdir", "-p", "/home/"+*username)
	// move the files
	executeCommand("lxc", "exec", *container, "--", "find", "/mnt/lxd_conf", "-name", "*.tar.gz.*", "-exec", "cp", "{}", "-t", "/home/"+*username, ";")
	// delete the backups
	executeCommand("find", *backupLocation, "-name", "*.tar.gz.*", "-exec", "rm", "{}", ";")
	fmt.Println("")
}

func (b *Backup) deleteBackupImage() {
	fmt.Printf("Deleting temporary image %s", backupName)
	executeCommand("lxc", "image", "delete", backupName)
	fmt.Println("")
}
