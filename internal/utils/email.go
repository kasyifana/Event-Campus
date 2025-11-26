package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"gopkg.in/gomail.v2"
)

// EmailSender handles email sending
type EmailSender struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
}

// NewEmailSender creates a new email sender
func NewEmailSender(host string, port int, user, password string) *EmailSender {
	return &EmailSender{
		smtpHost:     host,
		smtpPort:     port,
		smtpUser:     user,
		smtpPassword: password,
	}
}

// SendEmail sends an email
func (e *EmailSender) SendEmail(to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.smtpUser, e.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Email templates

// SendRegistrationConfirmation sends registration confirmation email
func (e *EmailSender) SendRegistrationConfirmation(to, userName, eventTitle string, eventDate time.Time, registrationID string) error {
	subject := fmt.Sprintf("Konfirmasi Pendaftaran: %s", eventTitle)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.info-box { background-color: white; padding: 15px; margin: 15px 0; border-left: 4px solid #4CAF50; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>✅ Pendaftaran Berhasil!</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.UserName}}</strong>,</p>
			<p>Terima kasih telah mendaftar untuk event:</p>
			
			<div class="info-box">
				<h2>{{.EventTitle}}</h2>
				<p>📅 <strong>Tanggal:</strong> {{.EventDate}}</p>
				<p>🎫 <strong>ID Pendaftaran:</strong> {{.RegistrationID}}</p>
			</div>

			<p>Anda akan menerima email reminder H-1 sebelum event dimulai.</p>
			<p><strong>Simpan ID pendaftaran ini untuk keperluan check-in.</strong></p>

			<p>Sampai jumpa di event!</p>
		</div>
		<div class="footer">
			<p>Event Campus - Platform Manajemen Event Kampus</p>
		</div>
	</div>
</body>
</html>
	`

	data := struct {
		UserName       string
		EventTitle     string
		EventDate      string
		RegistrationID string
	}{
		UserName:       userName,
		EventTitle:     eventTitle,
		EventDate:      eventDate.Format("Monday, 02 January 2006 - 15:04 WIB"),
		RegistrationID: registrationID,
	}

	var body bytes.Buffer
	t := template.Must(template.New("email").Parse(tmpl))
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return e.SendEmail(to, subject, body.String())
}

// SendWaitlistNotification sends waitlist notification email
func (e *EmailSender) SendWaitlistNotification(to, userName, eventTitle string, position int) error {
	subject := fmt.Sprintf("Waiting List: %s", eventTitle)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #FF9800; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.info-box { background-color: white; padding: 15px; margin: 15px 0; border-left: 4px solid #FF9800; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>⏳ Anda Masuk Waiting List</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.UserName}}</strong>,</p>
			<p>Event <strong>{{.EventTitle}}</strong> sudah penuh.</p>
			
			<div class="info-box">
				<p>📊 <strong>Posisi Anda:</strong> Nomor {{.Position}} di waiting list</p>
			</div>

			<p>Jika ada peserta yang membatalkan pendaftaran, kami akan segera menghubungi Anda!</p>

			<p>Terima kasih atas kesabaran Anda.</p>
		</div>
		<div class="footer">
			<p>Event Campus - Platform Manajemen Event Kampus</p>
		</div>
	</div>
</body>
</html>
	`

	data := struct {
		UserName   string
		EventTitle string
		Position   int
	}{
		UserName:   userName,
		EventTitle: eventTitle,
		Position:   position,
	}

	var body bytes.Buffer
	t := template.Must(template.New("email").Parse(tmpl))
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return e.SendEmail(to, subject, body.String())
}

// SendWaitlistPromotion sends waitlist promotion email
func (e *EmailSender) SendWaitlistPromotion(to, userName, eventTitle string, eventDate time.Time, registrationID string) error {
	subject := fmt.Sprintf("🎉 Promosi dari Waiting List: %s", eventTitle)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.info-box { background-color: white; padding: 15px; margin: 15px 0; border-left: 4px solid #4CAF50; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎉 Selamat! Pendaftaran Dikonfirmasi</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.UserName}}</strong>,</p>
			<p>Kabar baik! Sekarang Anda terdaftar untuk event:</p>
			
			<div class="info-box">
				<h2>{{.EventTitle}}</h2>
				<p>📅 <strong>Tanggal:</strong> {{.EventDate}}</p>
				<p>🎫 <strong>ID Pendaftaran:</strong> {{.RegistrationID}}</p>
			</div>

			<p>Anda dipromosikan dari waiting list karena ada pembatalan.</p>
			<p><strong>Simpan ID pendaftaran ini untuk keperluan check-in.</strong></p>

			<p>Sampai jumpa di event!</p>
		</div>
		<div class="footer">
			<p>Event Campus - Platform Manajemen Event Kampus</p>
		</div>
	</div>
</body>
</html>
	`

	data := struct {
		UserName       string
		EventTitle     string
		EventDate      string
		RegistrationID string
	}{
		UserName:       userName,
		EventTitle:     eventTitle,
		EventDate:      eventDate.Format("Monday, 02 January 2006 - 15:04 WIB"),
		RegistrationID: registrationID,
	}

	var body bytes.Buffer
	t := template.Must(template.New("email").Parse(tmpl))
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return e.SendEmail(to, subject, body.String())
}

// SendCancellationConfirmation sends cancellation confirmation email
func (e *EmailSender) SendCancellationConfirmation(to, userName, eventTitle string) error {
	subject := fmt.Sprintf("Pembatalan Pendaftaran: %s", eventTitle)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #f44336; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>❌ Pembatalan Dikonfirmasi</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.UserName}}</strong>,</p>
			<p>Pendaftaran Anda untuk event <strong>{{.EventTitle}}</strong> telah dibatalkan.</p>

			<p>Jika Anda berubah pikiran, silakan daftar kembali (jika masih ada slot tersedia).</p>

			<p>Terima kasih.</p>
		</div>
		<div class="footer">
			<p>Event Campus - Platform Manajemen Event Kampus</p>
		</div>
	</div>
</body>
</html>
	`

	data := struct {
		UserName   string
		EventTitle string
	}{
		UserName:   userName,
		EventTitle: eventTitle,
	}

	var body bytes.Buffer
	t := template.Must(template.New("email").Parse(tmpl))
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return e.SendEmail(to, subject, body.String())
}

// SendReminderEmail sends H-1 reminder email
func (e *EmailSender) SendReminderEmail(to, userName, eventTitle string, eventDate time.Time, location string, zoomLink *string, registrationID string) error {
	subject := fmt.Sprintf("[Reminder] Event Besok: %s", eventTitle)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #2196F3; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.info-box { background-color: white; padding: 15px; margin: 15px 0; border-left: 4px solid #2196F3; }
		.zoom-link { background-color: #4CAF50; color: white; padding: 10px; border-radius: 5px; text-align: center; margin: 10px 0; }
		.zoom-link a { color: white; text-decoration: none; font-weight: bold; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>⏰ Reminder: Event Besok!</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.UserName}}</strong>,</p>
			<p>Ini adalah pengingat bahwa kamu terdaftar untuk event:</p>
			
			<div class="info-box">
				<h2>{{.EventTitle}}</h2>
				<p>📅 <strong>Waktu:</strong> {{.EventDate}}</p>
				{{if .ZoomLink}}
				<p>💻 <strong>Event Online</strong></p>
				<div class="zoom-link">
					<a href="{{.ZoomLink}}" target="_blank">🔗 Klik Disini Untuk Join Zoom</a>
				</div>
				<p><small>💡 Link akan aktif 15 menit sebelum event dimulai</small></p>
				{{else}}
				<p>📍 <strong>Lokasi:</strong> {{.Location}}</p>
				{{end}}
				<p>🎫 <strong>Registration ID:</strong> {{.RegistrationID}}</p>
			</div>

			<p><strong>Jangan lupa untuk hadir!</strong></p>
			{{if not .ZoomLink}}
			<p>Tunjukkan Registration ID saat check-in.</p>
			{{end}}

			<p>Sampai jumpa besok!</p>
		</div>
		<div class="footer">
			<p>Event Campus - Platform Manajemen Event Kampus</p>
		</div>
	</div>
</body>
</html>
	`

	data := struct {
		UserName       string
		EventTitle     string
		EventDate      string
		Location       string
		ZoomLink       *string
		RegistrationID string
	}{
		UserName:       userName,
		EventTitle:     eventTitle,
		EventDate:      eventDate.Format("Monday, 02 January 2006 - 15:04 WIB"),
		Location:       location,
		ZoomLink:       zoomLink,
		RegistrationID: registrationID,
	}

	var body bytes.Buffer
	t := template.Must(template.New("email").Parse(tmpl))
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return e.SendEmail(to, subject, body.String())
}

// SendWhitelistApproval sends whitelist approval email
func (e *EmailSender) SendWhitelistApproval(to, userName, orgName string) error {
	subject := "✅ Pengajuan Organisasi Disetujui"

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.info-box { background-color: white; padding: 15px; margin: 15px 0; border-left: 4px solid #4CAF50; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎉 Selamat!</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.UserName}}</strong>,</p>
			
			<div class="info-box">
				<p>Pengajuan organisasi <strong>{{.OrgName}}</strong> telah disetujui!</p>
			</div>

			<p>Sekarang Anda memiliki akses untuk:</p>
			<ul>
				<li>✅ Membuat event baru</li>
				<li>✅ Mengelola event yang Anda buat</li>
				<li>✅ Melihat daftar peserta</li>
				<li>✅ Mark attendance</li>
			</ul>

			<p>Silakan login kembali untuk mulai membuat event!</p>
		</div>
		<div class="footer">
			<p>Event Campus - Platform Manajemen Event Kampus</p>
		</div>
	</div>
</body>
</html>
	`

	data := struct {
		UserName string
		OrgName  string
	}{
		UserName: userName,
		OrgName:  orgName,
	}

	var body bytes.Buffer
	t := template.Must(template.New("email").Parse(tmpl))
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return e.SendEmail(to, subject, body.String())
}

// SendWhitelistRejection sends whitelist rejection email
func (e *EmailSender) SendWhitelistRejection(to, userName, reason string) error {
	subject := "❌ Pengajuan Organisasi Ditolak"

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #f44336; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.info-box { background-color: white; padding: 15px; margin: 15px 0; border-left: 4px solid #f44336; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Pengajuan Ditolak</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.UserName}}</strong>,</p>
			<p>Mohon maaf, pengajuan organisasi Anda ditolak.</p>
			
			<div class="info-box">
				<p><strong>Alasan:</strong></p>
				<p>{{.Reason}}</p>
			</div>

			<p>Anda dapat mengajukan kembali dengan melengkapi persyaratan yang diminta.</p>
		</div>
		<div class="footer">
			<p>Event Campus - Platform Manajemen Event Kampus</p>
		</div>
	</div>
</body>
</html>
	`

	data := struct {
		UserName string
		Reason   string
	}{
		UserName: userName,
		Reason:   reason,
	}

	var body bytes.Buffer
	t := template.Must(template.New("email").Parse(tmpl))
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return e.SendEmail(to, subject, body.String())
}
