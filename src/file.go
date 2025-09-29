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
  "strings"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/disk"
)

type FileInfo struct {
	Name            string `json:"name"`
	Type            string `json:"type"` 
	IsSymlink       bool   `json:"is_symlink"`
	SizeBytes       int64  `json:"size_bytes"`
	Permissions     string `json:"permissions"`
	HardLinks       uint64 `json:"hard_links"`
	Inode           uint64 `json:"inode"`
	OwnerUID        uint32 `json:"owner_uid"`
	OwnerGID        uint32 `json:"owner_gid"`
	BlockSize       int64  `json:"block_size"`
	BlocksAllocated int64  `json:"blocks_allocated"`
	LastModified    string `json:"last_modified"`
	LastAccess      string `json:"last_access"`
}

type PartitionContent struct {
	MountPath string     `json:"mount_path"`
	Files     []FileInfo `json:"files"`
}

func getMountPoint(partition string) (string, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return "", fmt.Errorf("Erreur lors de la récupération des partitions : %v", err)
	}

	for _, p := range partitions {
		if p.Device == partition {
			return p.Mountpoint, nil
		}
	}

	return "", fmt.Errorf("Aucun point de montage trouvé pour la partition %s", partition)
}

func getFileInfo(path string, entry os.DirEntry) (FileInfo, error) {
	info, err := entry.Info()
	if err != nil {
		return FileInfo{}, fmt.Errorf("Impossible d'obtenir les informations de %s : %v", path, err)
	}

	isSymlink := info.Mode()&os.ModeSymlink != 0

	var fileType string
	switch {
	case isSymlink:
		fileType = "symlink"
	case info.IsDir():
		fileType = "directory"
	default:
		fileType = "file"
	}

	var stat syscall.Stat_t
	if err := syscall.Lstat(path, &stat); err != nil {
		return FileInfo{}, fmt.Errorf("Impossible d'obtenir les métadonnées avancées de %s : %v", path, err)
	}

	formatDate := func(ts syscall.Timespec) string {
		if ts.Sec == 0 {
			return "Unknown"
		}
		return time.Unix(ts.Sec, 0).Format("2006-01-02 15:04:05")
	}

	return FileInfo{
		Name:            info.Name(),
		Type:            fileType,
		IsSymlink:       isSymlink,
		SizeBytes:       info.Size(),
		Permissions:     info.Mode().String(),
		HardLinks:       stat.Nlink,
		Inode:           stat.Ino,
		OwnerUID:        stat.Uid,
		OwnerGID:        stat.Gid,
		BlockSize:       stat.Blksize,
		BlocksAllocated: stat.Blocks,
		LastModified:    formatDate(stat.Mtim),
		LastAccess:      formatDate(stat.Atim),
	}, nil
}

func AddTrailingSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		return path + "/"
	}
	return path
}

func AddSlash(s string) string {
	if len(s) == 0 || s[0] != '/' {
		return "/" + s
	}
	return s
}

func RemoveDoubleSlashes(s string) string {
	for strings.Contains(s, "//") {
		s = strings.ReplaceAll(s, "//", "/")
	}
	return s
}

func listRootFiles(mountPath string, path string) (string, error) {
  path = AddTrailingSlash(path)
  path = AddSlash(path)

  var content PartitionContent
	content.MountPath = mountPath

	entries, err := os.ReadDir(RemoveDoubleSlashes(mountPath + path))
	if err != nil {
		return "", fmt.Errorf("Le dossier %s n'existe pas : %v", RemoveDoubleSlashes(mountPath + path), err)
	}

	for _, entry := range entries {
		fileInfo, err := getFileInfo(RemoveDoubleSlashes(mountPath + path) + entry.Name(), entry)
		if err != nil {
			fmt.Println(err) 
			continue
		}
		content.Files = append(content.Files, fileInfo)
	}

	jsonData, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Erreur lors de la génération du JSON : %v", err)
	}

	return string(jsonData), nil
}
