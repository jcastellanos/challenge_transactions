package usecase

import "github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"

type EmailPort interface {
	SendEmail(to string, subject string, statistics model.Statistics, attachmentPath string) error
}
