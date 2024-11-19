package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *ArchiveController) CreateArchive(c *gin.Context) {

	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := c.Request.MultipartForm.File["files[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer file.Close()

		mimeType, err := getFileTypeByContent(file, fileHeader.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect file type"})
			return
		}

		if !isValidFileMimeType(mimeType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Houston, we have a problem. One of these files is an alien."})
			return
		}

		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file in zip"})
			return
		}

		_, err = io.Copy(zipFile, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file to zip"})
			return
		}
	}

	err = zipWriter.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close zip archive"})
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

func getFileTypeByContent(file io.Reader, fileName string) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)

	if strings.HasSuffix(fileName, ".docx") {
		mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	if strings.HasSuffix(fileName, ".xml") {
		mimeType = "application/xml"
	}
	if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		mimeType = "image/jpeg"
	}
	if strings.HasSuffix(fileName, ".png") {
		mimeType = "image/png"
	}

	return mimeType, nil
}

func isValidFileMimeType(mimeType string) bool {
	return mimeType == "image/jpeg" || mimeType == "image/png" ||
		mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" ||
		mimeType == "application/xml"
}
