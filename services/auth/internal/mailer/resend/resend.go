package resend

import (
	"context"
	"fmt"
	"github.com/resend/resend-go/v3"
)

type Resend struct {
	From   string
	ApiKey string
	Client *resend.Client
}

func NewResendMailer(from, apikey string) *Resend {
	return &Resend{
		From:   from,
		ApiKey: apikey,
		Client: resend.NewClient(apikey),
	}
}

func (r *Resend) Send(ctx context.Context, to string, subject string, html string, userID string) error {

	options := &resend.SendEmailOptions{
		IdempotencyKey: r.CreateIdempotencyKey("activation-user", userID),
	}

	params := &resend.SendEmailRequest{
		From:    r.From,
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	sent, err := r.Client.Emails.SendWithOptions(ctx, params, options)
	if err != nil {
		return err
	}

	fmt.Println(sent)

	return nil
}

func (r *Resend) CreateIdempotencyKey(emailType, userID string) string {
	return emailType + "/" + userID
}
