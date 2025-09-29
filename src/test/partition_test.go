package main

import (
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDevPrefix(t *testing.T) {
	assert.Equal(t, "sda1", removeDevPrefix("/dev/sda1"))
	assert.Equal(t, "sda1", removeDevPrefix("sda1"))
}

func TestEnsureLeadingSlash(t *testing.T) {
	assert.Equal(t, "/path", ensureLeadingSlash("path"))
	assert.Equal(t, "/path", ensureLeadingSlash("/path"))
}

func TestMountUnmountedPartitions(t *testing.T) {
	imagePath := "/tmp/testdisk.img"
	mountPath := os.Getenv("HOME") + "/testdisk"

	if _, err := os.Stat(mountPath); os.IsNotExist(err) {
		err = os.Mkdir(mountPath, 0755)
		if err != nil {
			log.Fatalf("Error creating mount point: %v", err)
		}
	}

	cmd := exec.Command("dd", "if=/dev/zero", "of="+imagePath, "bs=1M", "count=1024")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error creating disk image: %v", err)
	}

	cmd = exec.Command("mkfs.ext4", imagePath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error formatting disk image: %v", err)
	}

	cmd = exec.Command("sudo", "mount", "-o", "loop", imagePath, mountPath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error mounting disk image: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = MountUnmountedPartitions("/dev/loop0")
	assert.NoError(t, err)

	cmd = exec.Command("sudo", "umount", mountPath)
	err = cmd.Run()
	assert.NoError(t, err)

	err = os.Remove(imagePath)
	assert.NoError(t, err)

	err = os.Remove(mountPath)
	assert.NoError(t, err)
}

func TestGetPartitionsInfo(t *testing.T) {
	imagePath := "/tmp/testdisk.img"
	mountPath := os.Getenv("HOME") + "/testdisk"

	if _, err := os.Stat(mountPath); os.IsNotExist(err) {
		err = os.Mkdir(mountPath, 0755)
		if err != nil {
			log.Fatalf("Error creating mount point: %v", err)
		}
	}

	cmd := exec.Command("dd", "if=/dev/zero", "of="+imagePath, "bs=1M", "count=1024")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error creating disk image: %v", err)
	}

	cmd = exec.Command("mkfs.ext4", imagePath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error formatting disk image: %v", err)
	}

	cmd = exec.Command("sudo", "mount", "-o", "loop", imagePath, mountPath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error mounting disk image: %v", err)
	}

	time.Sleep(2 * time.Second)

	partitions, err := GetPartitionsInfo("/dev/loop0")
	assert.NoError(t, err)
	assert.NotEmpty(t, partitions)

	cmd = exec.Command("sudo", "umount", mountPath)
	err = cmd.Run()
	assert.NoError(t, err)

	err = os.Remove(imagePath)
	assert.NoError(t, err)

	err = os.Remove(mountPath)
	assert.NoError(t, err)
}

func TestGetPartitionsForDiskJSON(t *testing.T) {
	imagePath := "/tmp/testdisk.img"
	mountPath := os.Getenv("HOME") + "/testdisk"

	if _, err := os.Stat(mountPath); os.IsNotExist(err) {
		err = os.Mkdir(mountPath, 0755)
		if err != nil {
			log.Fatalf("Error creating mount point: %v", err)
		}
	}

	cmd := exec.Command("dd", "if=/dev/zero", "of="+imagePath, "bs=1M", "count=1024")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error creating disk image: %v", err)
	}

	cmd = exec.Command("mkfs.ext4", imagePath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error formatting disk image: %v", err)
	}

	cmd = exec.Command("sudo", "mount", "-o", "loop", imagePath, mountPath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error mounting disk image: %v", err)
	}

	time.Sleep(2 * time.Second)

	jsonData, err := GetPartitionsForDiskJSON("/dev/loop0")
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	cmd = exec.Command("sudo", "umount", mountPath)
	err = cmd.Run()
	assert.NoError(t, err)

	err = os.Remove(imagePath)
	assert.NoError(t, err)

	err = os.Remove(mountPath)
	assert.NoError(t, err)
}
