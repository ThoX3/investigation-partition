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
	"encoding/hex"
	"encoding/json"
	"os"
)

func ReadBlock(file *os.File, blockNumber int) (string, error) {
	const blockSize = 512

	offset := int64(blockNumber * blockSize)
	_, err := file.Seek(offset, 0)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, blockSize)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return "", err
	}

	hexData := hex.EncodeToString(buffer[:n])
	return hexData, nil
}

func GenerateJSON(filePath string) (string, error) {
	const blockSize = 512

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()
	maxBlocks := (fileSize + blockSize - 1) / blockSize 

	blocks := make(map[int]string)
	for i := int64(0); i < maxBlocks; i++ {
		hexData, err := ReadBlock(file, int(i))
		if err != nil {
			return "", err
		}
		blocks[int(i)] = hexData
	}

	jsonData, err := json.MarshalIndent(blocks, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
