package main

import (
	"flag"
	"fmt"
	"os"

	baker "github.com/nicholasjackson/crostini-backup"
)

var version = "dev"
var backupName = ""

// flags
var container = flag.String("container", "penguin", "Name of the container to backup or restore")
var snapShotName = flag.String("snapshot-prefix", "backup-snapshot", "Prefix of the snapshot to take, snapshots will have time appended to them")
var terminaLocation = flag.String("termina-location", "/mnt/stateful/lxd_conf", "Location to store temporary backup files in termina")
var archiveLocation = flag.String("archive-location", "/home/chronos", "Location of archive files, this is generally the home folder of your current container")
var help = flag.Bool("help", true, "Show help menu")

func main() {
	fmt.Println("Crostini Backup and Restore tool version:", version)
	flag.Parse()

	if len(os.Args) < 2 || *help == true {
		printUsage()
	}

	// get the operation
	op := os.Args[1]
	switch op {
	case "backup":
		doBackup()
	case "restore":
		doRestore()
	default:
		printUsage()
	}

}

func printUsage() {
	fmt.Println("Usage: ./bakery backup|restore [flags]")
	flag.Usage()
	os.Exit(0)
}

func doBackup() {
	fmt.Println("Starting backup, WARNING: This operation can take a long time")
	b := baker.Backup{
		BackupName:              "backup",
		BackupLocation:          *terminaLocation,
		ContainerBackupLocation: *archiveLocation,
		ContainerName:           *container,
		SnapShotName:            *snapShotName,
	}
	b.Execute()
}

func doRestore() {
	fmt.Println("Starting restore, WARNING: This operation can take a long time")
}
