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
//   - Kalash Abdulaziz  (abdulaziz.kalash@ecole.ensicaen.fr)
//   - Yahya Chikar      (yahya.chikar@ecole.ensicaen.fr)
//   - Antony Huynh      (antony.huynh@ecole.ensicaen.fr)
//   - Maelys Sable      (maelys.sable@ecole.ensicaen.fr)
//   - Yam Pakzad        (yam.pakzad@ecole.ensicaen.fr)
// ============================================================================

package main

import (
  "encoding/json"
  "os"
  "net/http"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

  r.GET("/disks", func(c *gin.Context) {
    data, err := GetDisksInfoJSON()
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    c.JSON(http.StatusOK, json.RawMessage(data))
  })

  r.GET("/partitions", func(c *gin.Context) {
    disk := c.Query("disk")
    if disk == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "disk parameter is required"})
      return
    }
    data, err := GetPartitionsForDiskJSON(disk)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    c.JSON(http.StatusOK, json.RawMessage(data))
  })

  r.GET("/files", func(c *gin.Context) {
    partition := c.Query("partition")
    if partition == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "partition parameter is required"})
      return
    }

    path := c.Query("path")
    filter := c.Query("filter")

    mountPath, err := getMountPoint(partition)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "No mount point found for the given partition"})
      return
    }

    if _, err := os.Stat(mountPath); os.IsNotExist(err) {
      c.JSON(http.StatusBadRequest, gin.H{"error": "The mount point is not accessible"})
      return
    }

    data, err := listRootFiles(mountPath, path)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    if filter == "" {
      c.JSON(http.StatusOK, json.RawMessage(data))
    } else {
      filteredData, err := FilterFilesByName([]byte(data), filter)
      if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
      }

      c.JSON(http.StatusOK, json.RawMessage(filteredData))
    }
  })

  r.GET("/blocks", func(c *gin.Context) {
    partition := c.Query("partition")
    if partition == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "partition parameter is required"})
      return
    }

    path := c.Query("path")
    if path == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
      return
    }

    mountPath, err := getMountPoint(partition)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "No mount point found for the given partition"})
      return
    }

    if _, err := os.Stat(mountPath); os.IsNotExist(err) {
      c.JSON(http.StatusBadRequest, gin.H{"error": "The mount point is not accessible"})
      return
    }

    path = AddSlash(path)

	  blocks, err := GenerateJSON(RemoveDoubleSlashes(mountPath + path))

	  if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		  return
	  }
    
    c.JSON(http.StatusOK, json.RawMessage(blocks))
  })


  r.Run("0.0.0.0:8080")
}
