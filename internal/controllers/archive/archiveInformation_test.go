package archive

import (
	"archive/zip"
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestArchiveInformation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockController := &ArchiveController{}
	r.POST("/api/archive/information", mockController.ArchiveInformation)

	var zipFileBuf bytes.Buffer
	zipWriter := zip.NewWriter(&zipFileBuf)
	fileWriter, err := zipWriter.Create("test.txt")
	if err != nil {
		t.Fatal("Failed to create zip file:", err)
	}
	_, err = fileWriter.Write([]byte("dummy file content"))
	if err != nil {
		t.Fatal("Failed to write content to zip file:", err)
	}
	err = zipWriter.Close()
	if err != nil {
		t.Fatal("Failed to close zip file:", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	formFile, err := writer.CreateFormFile("file", "test.zip")
	if err != nil {
		t.Fatal("Failed to create form file:", err)
	}
	_, err = formFile.Write(zipFileBuf.Bytes())
	if err != nil {
		t.Fatal("Failed to write zip content to form file:", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/information", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArchiveInformation_NotAZipFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockController := &ArchiveController{}
	r.POST("/api/archive/information", mockController.ArchiveInformation)

	var buf bytes.Buffer
	_, err := buf.Write([]byte("This is just a regular text file, not a zip file."))
	if err != nil {
		t.Fatal("Failed to write content to the buffer:", err)
	}

	var multipartBuf bytes.Buffer
	writer := multipart.NewWriter(&multipartBuf)
	part, err := writer.CreateFormFile("file", "not_a_zip.txt")
	if err != nil {
		t.Fatal("Failed to create form file:", err)
	}
	_, err = part.Write(buf.Bytes())
	if err != nil {
		t.Fatal("Failed to write text file content to form file:", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/information", &multipartBuf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
