package baker

import (
	"fmt"
	"time"
)

// Backup a Crostini Container
type Backup struct {
	BackupName       string
	BackupLocation   string
	ArchiveLocation  string
	ArchiveContainer string
	ContainerName    string
	SnapShotName     string

	General
}

// Execute starts the backup
func (b *Backup) Execute() {
	fmt.Println("")
	b.createSnapShot()
	b.stopContainer(b.ContainerName)
	b.publishBackup()
	b.exportImage()
	b.splitBackup()
	b.startContainer(b.ContainerName)
	b.mountBackupPath(b.BackupLocation, b.ArchiveContainer)
	b.moveBackupFiles()
	b.deleteBackupImage()

	fmt.Println("")
}

func (b *Backup) createSnapShot() {
	b.SnapShotName = fmt.Sprintf("%s-%d", b.SnapShotName, time.Now().Nanosecond())
	fmt.Printf("Creating snapshot of container:%s name:%s\n", b.ContainerName, b.SnapShotName)

	b.executeCommand("lxc", "snapshot", b.ContainerName, b.SnapShotName)
	fmt.Println()
}

func (b *Backup) publishBackup() {
	fmt.Printf("Publish container: %s to %s\n", b.ContainerName, b.BackupName)
	fmt.Println("If the container publish is interupted, your container may be left in a bad state,")
	fmt.Printf(
		"in this instance you can restore the snapshot using the command: lxc restore %s %s\n",
		b.ContainerName,
		b.SnapShotName,
	)

	b.executeCommand("lxc", "publish", b.ContainerName, "--alias", b.BackupName)
	fmt.Println()
}

func (b *Backup) exportImage() {
	fmt.Printf("Exporting container to: %s\n", b.BackupLocation)
	b.executeCommand("lxc", "image", "export", b.BackupName, b.BackupLocation+"/"+b.BackupName)
	fmt.Println()
}

func (b *Backup) splitBackup() {
	backupFileLocation := b.BackupLocation + "/" + b.BackupName + ".tar.gz"

	fmt.Println("Splitting backup into 3GB chunks")
	b.executeCommand("split", "-b", "3GB", backupFileLocation, backupFileLocation+".")
	b.executeCommand("rm", backupFileLocation)
	fmt.Println()
}

func (b *Backup) moveBackupFiles() {
	fmt.Printf("Moving backup files to %s in container %s\n", b.ArchiveLocation, b.ArchiveContainer)
	// create the user folder if required
	b.executeCommand("lxc", "exec", b.ArchiveContainer, "--", "mkdir", "-p", b.ArchiveLocation)
	// move the files
	b.executeCommand("lxc", "exec", b.ArchiveContainer, "--", "find", "/mnt/lxd_conf", "-name", "*.tar.gz.*", "-exec", "cp", "{}", "-t", b.ArchiveLocation, ";")
	// delete the backups
	b.executeCommand("find", b.BackupLocation, "-name", "*.tar.gz.*", "-exec", "rm", "{}", ";")
	fmt.Println("")
}

func (b *Backup) deleteBackupImage() {
	fmt.Printf("Deleting temporary image %s", b.BackupName)
	b.executeCommand("lxc", "image", "delete", b.BackupName)
	fmt.Println("")
}
