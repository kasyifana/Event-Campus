package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Get SMTP configuration
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpUser == "" || smtpPassword == "" {
		log.Fatal("SMTP configuration not found in .env file. Please set SMTP_HOST, SMTP_USERNAME, and SMTP_PASSWORD")
	}

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", "kasyifana09@gmail.com")
	m.SetHeader("Subject", "Test Email dari Event Campus")
	m.SetBody("text/html", `
		<h1>✅ Email Test Berhasil!</h1>
		<p>Halo dari Event Campus Backend!</p>
		<p>Jika kamu menerima email ini, berarti konfigurasi SMTP sudah bekerja dengan baik.</p>
		<hr>
		<p><strong>SMTP Configuration:</strong></p>
		<ul>
			<li>Host: `+smtpHost+`</li>
			<li>Port: `+fmt.Sprintf("%d", smtpPort)+`</li>
			<li>User: `+smtpUser+`</li>
		</ul>
		<p style="color: #888; font-size: 12px;">Sent from Event Campus API</p>
	`)

	// Send email
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

	fmt.Println("🔄 Sending test email to kasyifana09@gmail.com...")
	fmt.Printf("📧 SMTP Host: %s:%d\n", smtpHost, smtpPort)
	fmt.Printf("👤 From: %s\n", smtpUser)

	if err := d.DialAndSend(m); err != nil {
		log.Fatal("❌ Failed to send email: ", err)
	}

	fmt.Println("✅ Email sent successfully!")
	fmt.Println("📬 Please check kasyifana09@gmail.com inbox (including spam folder)")
}
