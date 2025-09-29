package main

import (
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPhysicalDisks(t *testing.T) {

	imagePath := "/tmp/testdisk.img"
	mountPath := os.Getenv("HOME") + "/testdisk"

	// Create the mount point directory if it doesn't exist
	if _, err := os.Stat(mountPath); os.IsNotExist(err) {
		err = os.Mkdir(mountPath, 0755)
		if err != nil {
			log.Fatalf("Error creating mount point: %v", err)
		}
	}

	// Create the disk image
	cmd := exec.Command("dd", "if=/dev/zero", "of="+imagePath, "bs=1M", "count=1024")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error creating disk image: %v", err)
	}

	// Ext4 filesystem
	cmd = exec.Command("mkfs.ext4", imagePath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error formatting disk image: %v", err)
	}

	// Mount the disk image
	cmd = exec.Command("sudo", "mount", "-o", "loop", imagePath, mountPath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error mounting disk image: %v", err)
	}

	// Wait for the disk to be mounted correctly
	time.Sleep(2 * time.Second)

	disks, err := getPhysicalDisks()

	assert.NoError(t, err)
	assert.Contains(t, disks, "/dev/loop0", "Disk name should contain '/dev/loop0'")

	// Clean up
	cmd = exec.Command("sudo", "umount", mountPath)
	err = cmd.Run()
	assert.NoError(t, err)
	err = os.Remove(imagePath)
	assert.NoError(t, err)
	err = os.Remove(mountPath)
	assert.NoError(t, err)
}

func TestGetDiskCapacity(t *testing.T) {
	// Define paths for the disk image and mount point
	imagePath := "/tmp/testdisk.img"
	mountPath := os.Getenv("HOME") + "/testdisk"

	// Create the mount point directory if it doesn't exist
	if _, err := os.Stat(mountPath); os.IsNotExist(err) {
		err = os.Mkdir(mountPath, 0755)
		if err != nil {
			log.Fatalf("Error creating mount point: %v", err)
		}
	}

	// Create the disk image
	cmd := exec.Command("dd", "if=/dev/zero", "of="+imagePath, "bs=1M", "count=1024")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error creating disk image: %v", err)
	}

	// Ext4 filesystem
	cmd = exec.Command("mkfs.ext4", imagePath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error formatting disk image: %v", err)
	}

	// Mount the disk image
	cmd = exec.Command("sudo", "mount", "-o", "loop", imagePath, mountPath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error mounting disk image: %v", err)
	}

	// Wait for the disk to be mounted correctly
	time.Sleep(2 * time.Second)

	capacity, err := getDiskCapacity("/dev/loop0")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1073741824), capacity) // Verify that the capacity is 1 GB

	// Clean up
	cmd = exec.Command("sudo", "umount", mountPath)
	err = cmd.Run()
	assert.NoError(t, err)
	err = os.Remove(imagePath)
	assert.NoError(t, err)
	err = os.Remove(mountPath)
	assert.NoError(t, err)
}

func TestTrimLastDigit(t *testing.T) {
	t.Run("Trim digits from disk name", func(t *testing.T) {
		result := trimLastDigit("sda1234")
		assert.Equal(t, "sda", result)
	})
	t.Run("No digits to trim", func(t *testing.T) {
		result := trimLastDigit("sda")
		assert.Equal(t, "sda", result)
	})

	t.Run("Empty disk name", func(t *testing.T) {
		result := trimLastDigit("")
		assert.Equal(t, "", result)
	})
}
