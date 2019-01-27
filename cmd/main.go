package main

import (
	"flag"
	"fmt"
	"os"

	bakery "github.com/nicholasjackson/crostini-backup"
)

var version = "dev"
var backupName = ""

var bfContainer *string
var bfArchiveLocation *string
var bfArchiveContainer *string
var bfSnapShotName *string
var bfTerminaLocation *string

var rfContainer *string
var rfArchiveLocation *string
var rfArchiveContainer *string
var rfTerminaLocation *string

func main() {

	bf := flag.NewFlagSet("backup", flag.ExitOnError)
	bfContainer = bf.String("container", "penguin", "Name of the container to backup")
	bfArchiveLocation = bf.String("archive-location", "/home/chronos", "Location of archive files, this is generally the home folder of your current container")
	bfArchiveContainer = bf.String("archive-container", "penguin", "Name of the container to write backup files")
	bfSnapShotName = bf.String("snapshot-prefix", "backup-snapshot", "Prefix of the snapshot to take, snapshots will have time appended to them")
	bfTerminaLocation = bf.String("termina-location", "/mnt/stateful/lxd_conf", "Location to store temporary backup files in termina")

	rf := flag.NewFlagSet("restore", flag.ExitOnError)
	rfContainer = rf.String("container", "penguin", "Name of the container to restore")
	rfArchiveLocation = rf.String("archive-location", "/home/chronos", "Location of archive files, this is generally the home folder of your current container")
	rfArchiveContainer = rf.String("archive-container", "penguin", "Name of the container to read backup files from")
	rfTerminaLocation = rf.String("termina-location", "/mnt/stateful/lxd_conf", "Location to store temporary backup files in termina")

	fmt.Println("Bakery - Crostini Backup and Restore tool")
	fmt.Println("version:", version)
	fmt.Println("")

	if len(os.Args) == 1 {
		printUsage()
	}

	// get the operation
	op := os.Args[1]
	switch op {
	case "backup":
		if len(os.Args) >= 2 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
			bf.Usage()
			os.Exit(0)
		}

		bf.Parse(os.Args[2:])
		doBackup()
	case "restore":
		if len(os.Args) >= 2 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
			rf.Usage()
			os.Exit(0)
		}

		rf.Parse(os.Args[2:])
		doRestore()
	default:
		printUsage()
	}

}

func printUsage() {
	fmt.Println("Usage: ./bakery backup|restore [flags]")
	fmt.Println("  --help")
	fmt.Println("        Show help menu")
	os.Exit(2)
}

func doBackup() {
	fmt.Println("Starting backup, WARNING: This operation can take a long time")
	b := bakery.Backup{
		BackupName:       "backup",
		BackupLocation:   *bfTerminaLocation,
		ArchiveLocation:  *bfArchiveLocation,
		ArchiveContainer: *bfArchiveContainer,
		ContainerName:    *bfContainer,
		SnapShotName:     *bfSnapShotName,
	}
	b.Execute()
}

func doRestore() {
	fmt.Println("Starting restore, WARNING: This operation can take a long time")
	r := bakery.Restore{
		BackupName:       "backup",
		TerminaLocation:  *rfTerminaLocation,
		ArchiveLocation:  *rfArchiveLocation,
		ArchiveContainer: *rfArchiveContainer,
		ContainerName:    *rfContainer,
	}
	r.Execute()
}
