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

	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
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

func (ga gmailEmailAdapter) SendEmail(to string, subject string, statistics model.Statistics, attachmentPath string) error {
	html, err := htmlFromTemplate("template/template.html", statistics)
	if err != nil {
		return err
	}
	imgData, err := ioutil.ReadFile("template/logo.png")
	if err != nil {
		return fmt.Errorf("error reading image: %w", err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)
	var msg bytes.Buffer
	writer := multipart.NewWriter(&msg)
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
	htmlPart, _ := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	htmlPart.Write([]byte(html))
	imagePart, _ := writer.CreatePart(map[string][]string{
		"Content-Type":              {"image/png"},
		"Content-Transfer-Encoding": {"base64"},
		"Content-Disposition":       {"inline; filename=\"image.png\""},
		"Content-ID":                {"<embedded-img>"},
	})
	imagePart.Write([]byte(imgBase64))
	fileData, err := ioutil.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("error reading the file to attach: %w", err)
	}
	attachmentPart, _ := writer.CreatePart(map[string][]string{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Disposition":       {fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(attachmentPath))},
		"Content-Transfer-Encoding": {"base64"},
	})
	attachmentPart.Write([]byte(base64.StdEncoding.EncodeToString(fileData)))
	writer.Close()
	auth := smtp.PlainAuth("", ga.emailConfig.Username, ga.emailConfig.Password, ga.emailConfig.SMTPServer)
	serverAddr := fmt.Sprintf("%s:%s", ga.emailConfig.SMTPServer, ga.emailConfig.Port)
	err = smtp.SendMail(serverAddr, auth, ga.emailConfig.Username, []string{to}, msg.Bytes())
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

func htmlFromTemplate(templatePath string, statistics model.Statistics) (string, error) {
	variables := map[string]interface{}{
		"totalBalance":  fmt.Sprintf("%.2f", statistics.TotalBalance()),
		"averageCredit": fmt.Sprintf("%.2f", statistics.AverageCredit()),
		"averageDebit":  fmt.Sprintf("%.2f", statistics.AverageDebit()),
	}
	tplContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading template: %w", err)
	}
	tpl, err := template.New("email").Parse(string(tplContent))
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var tplBuffer bytes.Buffer
	if err := tpl.Execute(&tplBuffer, variables); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}
	return tplBuffer.String(), nil
}
