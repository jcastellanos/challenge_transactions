package usecase

type EmailPort interface {
	SendEmail(to string, subject string, variables map[string]string, attachmentPath string) error
}
