# Bakery - Crostini backup and restore tool
Simple tool to backup and restore ChromeOS Crostini containers, influenced by the excelent readme on Reddit
[https://www.reddit.com/r/Crostini/wiki/howto/backup](https://www.reddit.com/r/Crostini/wiki/howto/backup)

[![CircleCI](https://circleci.com/gh/nicholasjackson/crostini-backup-restore.svg?style=svg)](https://circleci.com/gh/nicholasjackson/crostini-backup-restore)

## Demo Video
[https://www.useloom.com/share/71cdc4055744465f8f467f65cd26db44](https://www.useloom.com/share/71cdc4055744465f8f467f65cd26db44)

## Installation
* Open a Crosh terminal using ctrl+alt+t
* Start a new session `vsh termina`
* Copy the backup binary to /mnt/stateful/lxd_conf

The latest release can be found, in the `Releases` section, select the correct file for your archictecture.

Example Linux AMD64:
```bash
(termina) chronos@localhost ~ $ curl -L https://github.com/nicholasjackson/bakery/releases/download/v0.1.2/bakery_0.1.2_Linux_amd64.tar.gz -o /mnt/stateful/lxd_conf/backup.tar.gz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   639    0   639    0     0    193      0 --:--:--  0:00:03 --:--:--  2158
100  752k  100  752k    0     0   138k      0  0:00:05  0:00:05 --:--:--  660k

cd /mnt/stateful/lxd_conf
tar -zxf backup.tar.gz 
```

## Backup a Crostini container
A full list of options can be found by running bakery with the help flag

```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ ./bakery backup --help
Bakery - Crostini Backup and Restore tool
version: v0.1.2

Usage of backup:
  -archive-container string
        Name of the container to write backup files (default "penguin")
  -archive-location string
        Location of archive files, this is generally the home folder of your current container (default "/home/chronos")
  -container string
        Name of the container to backup (default "penguin")
  -snapshot-prefix string
        Prefix of the snapshot to take, snapshots will have time appended to them (default "backup-snapshot")
  -termina-location string
```

To backup a Crostini container use the following command, replace `container` with your own container name and `archive-location` with the location to store the backup files.  Generally this should be something like your home folder so you can copy the files using ChromeOS files:

```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ ./bakery backup -container tester -archive-location /home/jacksonnic
Bakery - Crostini Backup and Restore tool
version: v0.1.2

Starting backup, WARNING: This operation can take a long time

Creating snapshot of container:tester name:backup-snapshot-52723394

Stopping container tester

Publish container: tester to backup
If the container publish is interupted, your container may be left in a bad state,
in this instance you can restore the snapshot using the command: lxc restore tester backup-snapshot-52723394
Container published with fingerprint: 3b46b83105b3f2da09e70531b41705187cf58cc1e015eb14d1b1a778ef4b962f

Exporting container to: /mnt/stateful/lxd_conf
Image exported successfully!           

Splitting backup into 3GB chunks

Starting container tester
Mounting backup path /mnt/stateful/lxd_conf into container penguin
Mount path already exists

Moving backup files to /home/jacksonnic in container penguin

Deleting temporary image backup
```

The backup files will be output into your `archive-location` folder in 3GB chunks.  Once the backup has completed you can move these to external storage for safe storage.

## Restore a Crostini container
A full list of options can be found by running bakery with the help flag

```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ ./bakery restore -help
Bakery - Crostini Backup and Restore tool
version: v0.1.2

Usage of restore:
  -archive-container string
        Name of the container to read backup files from (default "penguin")
  -archive-location string
        Location of archive files, this is generally the home folder of your current container (default "/home/chronos")
  -container string
        Name of the container to restore (default "penguin")
  -termina-location string
        Location to store temporary backup files in termina (default "/mnt/stateful/lxd_conf")
```

To restore a container, first copy your backup archive to a running Crostini container. You can then use the following command replacing the value of the container flag with the name to which you want to restore your backup and the archive-location to the location of your backup files in a running container. Generally this is the home folder. 

```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ ./bakery restore --container tester2 -archive-location /home/jacksonnic
Bakery - Crostini Backup and Restore tool
version: v0.1.2

Starting restore, WARNING: This operation can take a long time

Mounting backup path /mnt/stateful/lxd_conf into container penguin
Mount path already exists

Moving backup files from /home/jacksonnic in container penguin

Merge backup files back into a single archive /mnt/stateful/lxd_conf/backup.tar.gz

importing backup /mnt/stateful/lxd_conf/backup.tar.gz to image backup
Image imported with fingerprint: 3b46b83105b3f2da09e70531b41705187cf58cc1e015eb14d1b1a778ef4b962f

Initializing container tester2 from image backup
Creating tester2

Deleting backup image backup

```

## Testing
To test `backup` use an empty temporary container, this can be created using the following steps:

```bash
(termina) chronos@localhost /mnt/stateful/lxd_conf $ lxc init 980e37d286ad tester
Creating tester
(termina) chronos@localhost /mnt/stateful/lxd_conf $ lxc start tester
```

You can delete the temporary container using the following command

```bash
lxc delete tester --force
```
