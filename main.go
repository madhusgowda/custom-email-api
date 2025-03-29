package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type EmailRequest struct {
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

func sendEmail(to, subject, body string) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	// SMTP server configuration.
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Message.
	msg := []byte("From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/send-email", func(c *gin.Context) {
		var emailRequest EmailRequest

		if err := c.ShouldBindJSON(&emailRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := sendEmail(emailRequest.To, emailRequest.Subject, emailRequest.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully!"})
	})

	log.Println("Listening on port:" + os.Getenv("PORT"))
	r.Run(":8081")
}
