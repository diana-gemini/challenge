package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/diana-gemini/doodocs/internal/models"
	"github.com/diana-gemini/doodocs/pkg"
	"github.com/gin-gonic/gin"
)

type ArchiveController struct {
	Env *pkg.Env
}

func (h *ArchiveController) ArchiveInformation(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	var result models.ArchiveInformation

	result.Filename = fileHeader.Filename
	result.ArchiveSize = float64(fileHeader.Size)

	if !isZipArchive(file) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Elementary, my dear user... This file is not a zip archive."})
		return
	}

	file.Seek(0, io.SeekStart)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file into buffer"})
		return
	}

	totalSize, files, err := countFilesInZipArchive(&buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read zip archive"})
		return
	}
	result.TotalFiles = float64(len(files))
	result.TotalSize = float64(totalSize)
	result.Files = files

	c.JSON(http.StatusOK, models.SuccessResponse{Result: result})
}

func isZipArchive(file io.Reader) bool {
	buf := make([]byte, 4)
	_, err := file.Read(buf)
	if err != nil {
		return false
	}
	return buf[0] == 'P' && buf[1] == 'K' && buf[2] == 0x03 && buf[3] == 0x04
}

func countFilesInZipArchive(file io.Reader) (float64, []models.File, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(file.(*bytes.Buffer).Bytes()), int64(file.(*bytes.Buffer).Len()))
	if err != nil {
		return 0, nil, err
	}

	var files []models.File
	var totalSize float64
	for _, zipFile := range zipReader.File {
		var file models.File
		if !zipFile.FileInfo().IsDir() {
			file.FilePath = zipFile.Name
			file.Size = float64(zipFile.FileInfo().Size())
			file.Mimetype = mime.TypeByExtension(filepath.Ext(file.FilePath))
			totalSize += file.Size
			files = append(files, file)
		}
	}

	return totalSize, files, nil
}
