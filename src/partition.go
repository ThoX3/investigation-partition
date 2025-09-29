// ============================================================================
// ENSICAEN
// 6 Boulevard Maréchal Juin
// F-14050 Caen Cedex
//
// Projet 2A
// Investigation de Partition HFS+ et Ext4 en Go
//
// 2025
//
// Auteurs :
//   - Kalash Abdulaziz (abdulaziz.kalash@ecole.ensicaen.fr)
//   - Yahya Chikar      (yahya.chikar@ecole.ensicaen.fr)
//   - Antony Huynh      (antony.huynh@ecole.ensicaen.fr)
//   - Maelys Sable      (maelys.sable@ecole.ensicaen.fr)
//   - Yam Pakzad        (yam.pakzad@ecole.ensicaen.fr)
// ============================================================================


package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

type PartitionInfo struct {
	Name         string          `json:"name"`
	FSType       string          `json:"fs_type,omitempty"`
	MountPoint   string          `json:"mount_point,omitempty"`
	Size         string          `json:"size"`
	FSSize       string          `json:"fs_size,omitempty"`
	FSUsed       string          `json:"fs_used,omitempty"`
	UUID         string          `json:"uuid,omitempty"`
	PartTypeName string          `json:"part_type_name,omitempty"`
	PartNumber   json.RawMessage `json:"part_number,omitempty"`
	PhySec       json.RawMessage `json:"physical_sector_size,omitempty"`
	LogSec       json.RawMessage `json:"logical_sector_size,omitempty"`
}

type LsblkDevice struct {
	Name         string          `json:"name"`
	FSType       string          `json:"fstype,omitempty"`
	MountPoint   string          `json:"mountpoint,omitempty"`
	Size         string          `json:"size"`
	FSSize       string          `json:"fssize,omitempty"`
	FSUsed       string          `json:"fsused,omitempty"`
	UUID         string          `json:"uuid,omitempty"`
	PartTypeName string          `json:"parttypename,omitempty"`
	PartNumber   json.RawMessage `json:"partn,omitempty"`
	PhySec       json.RawMessage `json:"phy-sec,omitempty"`
	LogSec       json.RawMessage `json:"log-sec,omitempty"`
	Children     []LsblkDevice   `json:"children,omitempty"`
}

type LsblkJSON struct {
	BlockDevices []LsblkDevice `json:"blockdevices"`
}

func mountPartition(device string, mountPoint string, fstype string) error {
	cmd := exec.Command("sudo", "mount", "-t", fstype, device, mountPoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Erreur lors du montage: %v\nSortie: %s", err, string(output))
	}
	return nil
}

func unmountPartition(device string) error {
	cmd := exec.Command("sudo", "umount", device)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Erreur lors du démontage: %v\nSortie: %s", err, string(output))
	}
	return nil
}

func GetMountPoint() string {
	paths := []string{"/media", "/mnt"}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func removeDevPrefix(name string) string {
	return strings.TrimPrefix(name, "/dev")
}

func ensureLeadingSlash(s string) string {
	if !strings.HasPrefix(s, "/") {
		return "/" + s
	}
	return s
}

func MountUnmountedPartitions(diskName string) error {
	partitions, err := GetPartitionsInfo(diskName)
	if err != nil {
		return err
	}

	mountBase := GetMountPoint()
	if mountBase == "" {
		return fmt.Errorf("Aucun point de montage valide trouvé (/media ou /mnt)")
	}

	for _, part := range partitions {
		if part.MountPoint == "" && part.FSType != "" {
			mountDir := fmt.Sprintf("%s%s", mountBase + "/", removeDevPrefix(part.Name))

      if _, err := os.Stat(mountDir); os.IsNotExist(err) {
        if err := os.Mkdir(mountDir, 0755); err != nil {
          return fmt.Errorf("Erreur lors de la création du dossier %s : %v", mountDir, err)
        }
      }

			if err := mountPartition(part.Name, mountDir, part.FSType); err != nil {
				return err
			}
		}
	}

	return nil
}

func GetPartitionsInfo(diskName string) ([]PartitionInfo, error) {
	out, err := exec.Command("lsblk", "-J", "-o", "NAME,FSTYPE,MOUNTPOINT,SIZE,FSSIZE,FSUSED,UUID,PHY-SEC,LOG-SEC,PARTTYPENAME").Output()
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de l'exécution de lsblk : %v", err)
	}

	var lsblkData LsblkJSON
	err = json.Unmarshal(out, &lsblkData)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture du JSON lsblk : %v", err)
	}

	var partitions []PartitionInfo

	found := false
	for _, device := range lsblkData.BlockDevices {
		if "/dev/"+device.Name == diskName {
			found = true
			for _, part := range device.Children {
				partitions = append(partitions, formatPartitionInfo(part))
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("Aucune partition trouvée pour le disque %s", diskName)
	}

	return partitions, nil
}

func formatPartitionInfo(dev LsblkDevice) PartitionInfo {
	fsUsed := dev.FSUsed
	if fsUsed == "" && dev.MountPoint != "" {
		usage, err := disk.Usage(dev.MountPoint)
		if err == nil {
			fsUsed = fmt.Sprintf("%v", usage.Used)
		}
	}

	return PartitionInfo{
		Name:         "/dev/" + dev.Name,
		FSType:       dev.FSType,
		MountPoint:   dev.MountPoint,
		Size:         dev.Size,
		FSSize:       dev.FSSize,
		FSUsed:       fsUsed,
		UUID:         dev.UUID,
		PartTypeName: dev.PartTypeName,
		PartNumber:   dev.PartNumber,
		PhySec:       dev.PhySec,
		LogSec:       dev.LogSec,
	}
}

func GetPartitionsForDiskJSON(diskName string) (string, error) {
	MountUnmountedPartitions(diskName)
	partitions, err := GetPartitionsInfo(diskName)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(partitions, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
