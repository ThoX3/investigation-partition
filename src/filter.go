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
)

type FilesResponse struct {
	MountPath string     `json:"mount_path"`
	Files     []FileInfo `json:"files"`
}

func FilterFilesByName(jsonData []byte, filter string) ([]byte, error) {
	var response FilesResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, err
	}

	var filteredFiles []FileInfo
	for _, file := range response.Files {
		if strings.Contains(strings.ToLower(file.Name), strings.ToLower(filter)) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	response.Files = filteredFiles
	return json.Marshal(response)
}
