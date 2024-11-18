package archive

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateArchive_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockController := &ArchiveController{}
	r.POST("/api/archive/files", mockController.CreateArchive)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	file1, err := writer.CreateFormFile("files[]", "file1.docx")
	if err != nil {
		t.Fatal("Error creating form file:", err)
	}
	_, err = file1.Write([]byte("This is a DOCX file"))
	if err != nil {
		t.Fatal("Error writing file content:", err)
	}

	file2, err := writer.CreateFormFile("files[]", "file2.xml")
	if err != nil {
		t.Fatal("Error creating form file:", err)
	}
	_, err = file2.Write([]byte("<note><to>Tove</to><from>Jani</from></note>"))
	if err != nil {
		t.Fatal("Error writing file content:", err)
	}

	file3, err := writer.CreateFormFile("files[]", "file3.jpg")
	if err != nil {
		t.Fatal("Error creating form file:", err)
	}
	_, err = file3.Write([]byte("dummy image data"))
	if err != nil {
		t.Fatal("Error writing file content:", err)
	}

	file4, err := writer.CreateFormFile("files[]", "file4.png")
	if err != nil {
		t.Fatal("Error creating form file:", err)
	}
	_, err = file4.Write([]byte("dummy image data"))
	if err != nil {
		t.Fatal("Error writing file content:", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/files", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// if w.Code != http.StatusOK {
	// 	t.Errorf("Expected status 200 OK, got %d. Response body: %s", w.Code, w.Body.String())
	// }

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
}

func TestCreateArchiveWithInvalidFiles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockController := &ArchiveController{}
	r.POST("/api/archive/files", mockController.CreateArchive)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	file1, _ := writer.CreateFormFile("files[]", "file1.jpg") 
	file1.Write([]byte("This is an image file"))

	file2, _ := writer.CreateFormFile("files[]", "file2.txt") 
	file2.Write([]byte("This is a text file"))

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/files", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
