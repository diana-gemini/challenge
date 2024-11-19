package mail

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/diana-gemini/doodocs/internal/models"
	"github.com/diana-gemini/doodocs/pkg"
	"github.com/gin-gonic/gin"
	"gopkg.in/mail.v2"
)

type MailController struct {
	Env *pkg.Env
}

func (h *MailController) SendEmail(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get file from request"})
		return
	}
	defer file.Close()

	validMimeTypes := []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/pdf",
	}

	mimeType, err := getFileTypeByContent(file, fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect file type"})
		return
	}

	valid := false
	for _, validMime := range validMimeTypes {
		if mimeType == validMime {
			valid = true
			break
		}
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One small file for mankind, one giant error for this system."})
		return
	}

	emails := c.Request.MultipartForm.Value["emails[]"]
	if len(emails) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No emails uploaded"})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file"})
		return
	}

	fileReader := bytes.NewReader(fileBytes)

	message := mail.NewMessage()

	message.SetHeader("From", "diana-test-project@mail.ru")

	var cleanedEmails []string
	for _, email := range emails {
		email = strings.TrimSpace(email)
		if email != "" {
			cleanedEmails = append(cleanedEmails, email)
		}
	}

	if len(cleanedEmails) > 0 {
		message.SetHeader("To", cleanedEmails...)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email addresses"})
		return
	}

	message.SetHeader("Subject", "Mission Complete: The Files Have Arrived!")

	message.SetBody("text/html", `
    <html>
        <body>
            <h1>Captain, the mission is a success!</h1>
            <p>We've successfully transmitted the requested file back to Earth. </p>
            <p><i>Your loyal server</i></p>
        </body>
    </html>
`)

	message.AttachReader(fileHeader.Filename, fileReader)

	smtpHost := h.Env.SMTPHost
	smtpPort := h.Env.SMTPPort
	smtpUsername := h.Env.SMTPUser
	smtpPassword := h.Env.SMTPPassword

	dialer := mail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	if err := dialer.DialAndSend(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: "Email sent successfully with attachments!"})
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
	if strings.HasSuffix(fileName, ".pdf") {
		mimeType = "application/pdf"
	}
	return mimeType, nil
}
