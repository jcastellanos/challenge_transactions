package adapter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/smtp"
	"path/filepath"
)

type EmailConfig struct {
	SMTPServer string
	Port       string
	Username   string
	Password   string
}

type gmailEmailAdapter struct {
	emailConfig EmailConfig
}

func NewGmailEmailAdapter(emailConfig EmailConfig) gmailEmailAdapter {
	return gmailEmailAdapter{
		emailConfig: emailConfig,
	}
}

func (ga gmailEmailAdapter) SendEmail(to string, subject string, variables map[string]string, attachmentPath string) error {
	html, err := htmlFromTemplate("template/template.html", variables)
	if err != nil {
		return err
	}
	// Leer la imagen
	imgData, err := ioutil.ReadFile("template/logo.png")
	if err != nil {
		return fmt.Errorf("error leyendo la imagen: %w", err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	// Crear un mensaje MIME multipart
	var msg bytes.Buffer
	writer := multipart.NewWriter(&msg)

	// Configuración de encabezados
	headers := map[string]string{
		"From":         ga.emailConfig.Username,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary()),
	}
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")

	// Parte del correo con HTML
	htmlPart, _ := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	htmlPart.Write([]byte(html))

	// Parte del correo con la imagen embebida
	imagePart, _ := writer.CreatePart(map[string][]string{
		"Content-Type":              {"image/png"}, // Cambiar según el formato de la imagen
		"Content-Transfer-Encoding": {"base64"},
		"Content-Disposition":       {"inline; filename=\"image.png\""},
		"Content-ID":                {"<embedded-img>"},
	})
	imagePart.Write([]byte(imgBase64))
	// Adjuntar archivo plano
	fileData, err := ioutil.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("error leyendo el archivo adjunto: %w", err)
	}
	attachmentPart, _ := writer.CreatePart(map[string][]string{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Disposition":       {fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(attachmentPath))},
		"Content-Transfer-Encoding": {"base64"},
	})
	attachmentPart.Write([]byte(base64.StdEncoding.EncodeToString(fileData)))

	// Finalizar el mensaje MIME
	writer.Close()

	// Conexión al servidor SMTP
	auth := smtp.PlainAuth("", ga.emailConfig.Username, ga.emailConfig.Password, ga.emailConfig.SMTPServer)
	serverAddr := fmt.Sprintf("%s:%s", ga.emailConfig.SMTPServer, ga.emailConfig.Port)
	err = smtp.SendMail(serverAddr, auth, ga.emailConfig.Username, []string{to}, msg.Bytes())
	if err != nil {
		return fmt.Errorf("error enviando el correo: %w", err)
	}

	return nil
}

func htmlFromTemplate(templatePath string, variables map[string]string) (string, error) {
	// Leer y procesar la plantilla HTML
	tplContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error leyendo la plantilla: %w", err)
	}
	tpl, err := template.New("email").Parse(string(tplContent))
	if err != nil {
		return "", fmt.Errorf("error parseando la plantilla: %w", err)
	}

	var tplBuffer bytes.Buffer
	if err := tpl.Execute(&tplBuffer, variables); err != nil {
		return "", fmt.Errorf("error ejecutando la plantilla: %w", err)
	}
	return tplBuffer.String(), nil
}
