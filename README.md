# crostini-backup-restore
Simple tool to backup and restore ChromeOS Crostini containers

## Installation
* Open a Crosh termainal using ctrl+alt+p
* Start a new session `vsh termina`
* Copy the backup binary to /mnt/stateful/lxd_conf
```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ curl -l https://github.com/nicholasjackson/crostini-backup-restore/releases/download/v0.0.1/backup -o /mnt/stateful/lxd_conf/backup
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   597    0   597    0     0    736      0 --:--:-- --:--:-- --:--:--   790
```

## Using backup
A full list of options can be found by running backup with the help flag
```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ ./backup --help
Usage of ./backup:
  -backup-location string
        Location to backup container (default "/mnt/stateful/lxd_conf")
  -container string
        Name of the container to backup (default "penguin")
  -snapshot-prefix string
        Prefix of the snapshot to take, snapshots will have time appended to them (default "backup-snapshot")
  -username string
        Crostini username e.g. jacksonnic (default "chronos")
```

To backup a Crostini container use the following command, replace `container` with your own container name and `username` with your own username:

```
(termina) chronos@localhost /mnt/stateful/lxd_conf $ ./backup -container penguin --username jacksonnic
Crostini Backup and Restore tool version: dev
Starting backup, WARNING: This operation can take a long time

Stopping container penguin

Creating snapshot of container:penguin name:backup-snapshot-85563359

Publish container: penguin to penguin-backup
If the container publish is interupted, your container may be left in a bad state,
in this instance you can restore the snapshot using the command: lxc restore penguin backup-snapshot-85563359
Container published with fingerprint: d9468e72eeaa7d74a73eb39654fe278a14535a20b450077df672071c6b87d689

Exporting container to: /mnt/stateful/lxd_conf
Image exported successfully!           

Splitting backup into 3GB chunks

Starting container penguin
Mounting backup path /mnt/stateful/lxd_conf into container penguin
Device lxd-conf added to penguin

Moving backup files to /home/jacksonnic in container penguin
```

The backup files will be output into your `Linux Files` folder in 3GB chunks.  Once the backup has completed you can move these to external storage for safe storage.

## Testing
To test `backup` use an empty temporary container, this can be created using the following steps:

```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ lxc init 980e37d286ad tester
Creating tester
(termina) chronos@localhost /mnt/stateful/lxd_conf $ lxc start tester
```

You can then run a test backup on this container:
```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ ./backup -container tester
Crostini Backup and Restore tool version: dev
Starting backup, WARNING: This operation can take a long time

Stopping container tester

Creating snapshot of container:tester name:backup-snapshot-85563359

Publish container: tester to tester-backup
If the container publish is interupted, your container may be left in a bad state,
in this instance you can restore the snapshot using the command: lxc restore tester backup-snapshot-85563359
Container published with fingerprint: d9468e72eeaa7d74a73eb39654fe278a14535a20b450077df672071c6b87d689

Exporting container to: /mnt/stateful/lxd_conf
Image exported successfully!           

Splitting backup into 3GB chunks

Starting container tester
Mounting backup path /mnt/stateful/lxd_conf into container tester
Device lxd-conf added to tester

Moving backup files to /home/chronos in container tester
```

The backup files will be stored in the home folder `/home/chronos`, you can validate that these are present with the follwing command:

```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ lxc exec tester -- ls -las /home/chronos
total 318040
     0 drwxr-xr-x 1 root root        46 Jan 21 15:44 .
     0 drwxr-xr-x 1 root root        14 Jan 21 15:44 ..
318040 -rw-r--r-- 1 root root 325670095 Jan 21 15:44 tester-backup.tar.gz.aa
```

You can now delete the temporary container

```bash
lxc delete tester --force
```
