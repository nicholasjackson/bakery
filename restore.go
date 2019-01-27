package bakery

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Restore a Crostini Container
type Restore struct {
	BackupName       string
	TerminaLocation  string
	ArchiveLocation  string
	ArchiveContainer string
	ContainerName    string

	General
}

// Execute starts the restore
func (r *Restore) Execute() {
	fmt.Println("")

	r.mountBackupPath(r.TerminaLocation, r.ArchiveContainer)
	f := r.moveBackupFiles()
	r.joinBackupFiles(f)
	r.importBackup()
	r.initContainer()
	r.removeBackupImage()

	fmt.Println("")
}

func (r *Restore) moveBackupFiles() []string {
	fmt.Printf("Moving backup files from %s in container %s\n", r.ArchiveLocation, r.ArchiveContainer)

	// list the backup files
	output := bytes.NewBuffer([]byte{})
	errOutput := bytes.NewBuffer([]byte{})
	cmd := exec.Command("lxc", "exec", r.ArchiveContainer, "--", "find", r.ArchiveLocation, "-name", "*.tar.gz.*")
	cmd.Stdout = output
	cmd.Stderr = errOutput
	err := cmd.Run()

	if err != nil {
		fmt.Println("Unable to find archive files to restore")
		fmt.Println(errOutput.String())
		os.Exit(1)
	}

	out := output.String()
	if len(out) < 1 {
		fmt.Println("Unable to find archive files to restore")
		os.Exit(1)
	}

	files := strings.Split(out, "\n")
	for _, f := range files {
		if len(f) > 1 {
			r.executeCommand("lxc", "file", "pull", r.ArchiveContainer+f, r.TerminaLocation)
		}
	}

	fmt.Println("")
	return files
}

func (r *Restore) joinBackupFiles(files []string) {
	bn := r.TerminaLocation + "/" + r.BackupName + ".tar.gz"
	fmt.Println("Merge backup files back into a single archive", bn)

	// delete the existing backup if present
	r.executeCommand("rm", "-f", bn)

	for _, f := range files {
		if len(f) > 1 {
			fn := strings.Replace(f, r.ArchiveLocation, "", -1)
			r.executeCommand("bash", "-c", "cat "+r.TerminaLocation+fn+" >> "+bn)
			r.executeCommand("rm", "-f", r.TerminaLocation+fn)
		}
	}

	fmt.Println("")
}

func (r *Restore) importBackup() {
	bf := r.TerminaLocation + "/" + r.BackupName + ".tar.gz"
	fmt.Printf("Importing backup %s to image %s\n", bf, r.BackupName)

	r.executeCommand("lxc", "image", "import", bf, "--alias", r.BackupName)
	r.executeCommand("rm", "-f", bf)

	fmt.Println("")
}

func (r *Restore) initContainer() {
	fmt.Printf("Initializing container %s from image %s\n", r.ContainerName, r.BackupName)

	r.executeCommand("lxc", "init", r.BackupName, r.ContainerName)

	fmt.Println("")
}

func (r *Restore) removeBackupImage() {
	fmt.Printf("Deleting backup image %s\n", r.BackupName)

	r.executeCommand("lxc", "image", "delete", r.BackupName)

	fmt.Println("")
}
