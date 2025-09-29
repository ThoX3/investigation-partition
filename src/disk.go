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
	"os/exec"
	"runtime"
	"strconv"
	"strings"


	"github.com/shirou/gopsutil/v3/disk"
)

type DiskInfo struct {
	Name     string  `json:"name"`
	Capacity float64 `json:"capacity_gb"`
}

func GetDisksInfoJSON() (string, error) {
	var disksInfo []DiskInfo
	disks, err := getPhysicalDisks()
	if err != nil {
		return "", err
	}

	for _, diskName := range disks {
		capacity, err := getDiskCapacity(diskName)
		if err != nil {
			continue
		}

		disksInfo = append(disksInfo, DiskInfo{
			Name:     diskName,
			Capacity: roundFloat(float64(capacity)/(1024*1024*1024), 2),
		})
	}

	jsonData, err := json.MarshalIndent(disksInfo, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func getPhysicalDisks() ([]string, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []string
	seenDisks := make(map[string]bool)

	for _, part := range partitions {
		diskName := part.Device

		if strings.Contains(diskName, "loop") {
			continue
		}

		if strings.Contains(diskName, "nvme") && strings.Contains(diskName, "p") {
			diskName = strings.Split(diskName, "p")[0]

		} else if strings.HasPrefix(diskName, "/dev/sd") || strings.HasPrefix(diskName, "/dev/vd") {
			diskName = trimLastDigit(diskName)

		} else if strings.HasPrefix(diskName, "/dev/mmcblk") && strings.Contains(diskName, "p") {
			diskName = strings.Split(diskName, "p")[0]

		} else if strings.HasPrefix(diskName, "/dev/md") && strings.Contains(diskName, "p") {
			diskName = strings.Split(diskName, "p")[0]
		}

		if !seenDisks[diskName] {
			seenDisks[diskName] = true
			disks = append(disks, diskName)
		}
	}
	return disks, nil
}

func getDiskCapacity(diskName string) (uint64, error) {
	if runtime.GOOS == "linux" || isWSL() {
		return getLinuxDiskCapacity(diskName)
	}
	switch runtime.GOOS {
	case "windows":
		return getWindowsDiskCapacity(diskName)
	case "darwin":
		return getMacOSDiskCapacity(diskName)
	default:
		return 0, fmt.Errorf("OS non supporté")
	}
}

func getLinuxDiskCapacity(diskName string) (uint64, error) {
	diskName = strings.TrimPrefix(diskName, "/dev/")
	out, err := exec.Command("lsblk", "-b", "-dn", "-o", "SIZE", "/dev/"+diskName).Output()
	if err != nil {
		return 0, err
	}

	size, err := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func getWindowsDiskCapacity(diskName string) (uint64, error) {
	out, err := exec.Command("wmic", "diskdrive", "get", "size").Output()
	if err != nil {
		return 0, fmt.Errorf("Erreur lors de l'exécution de WMIC : %v", err)
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		sizeStr := strings.TrimSpace(line)
		if sizeStr != "" && sizeStr != "Size" {
			size, err := strconv.ParseUint(sizeStr, 10, 64)
			if err != nil {
				continue
			}
			return size, nil
		}
	}
	return 0, fmt.Errorf("Capacité non trouvée sur Windows")
}

func getMacOSDiskCapacity(diskName string) (uint64, error) {
	out, err := exec.Command("diskutil", "info", diskName).Output()
	if err != nil {
		return 0, fmt.Errorf("Erreur lors de l'exécution de diskutil : %v", err)
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Total Size") {
			parts := strings.Fields(line)
			if len(parts) > 2 {
				size, err := strconv.ParseUint(parts[len(parts)-2], 10, 64)
				if err != nil {
					return 0, fmt.Errorf("Erreur conversion en uint64 : %v", err)
				}
				return size, nil
			}
		}
	}
	return 0, fmt.Errorf("Capacité non trouvée sur macOS")
}

func trimLastDigit(diskName string) string {
	return strings.TrimRightFunc(diskName, func(r rune) bool {
		return r >= '0' && r <= '9'
	})
}

func isWSL() bool {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "WSL")
}
