package tasks

import (
	"context"
	"fmt"
)

const TypeSendWelcomeEmail = "send_welcome_email"

type SendWelcomeEmailPayload struct {
	UserID uint
	Email  string
	Name   string
}

func ExecuteSendWelcomeEmail(ctx context.Context, payload SendWelcomeEmailPayload) error {
	//Just printing to the console, later will change to SendGrid for real email sending
	fmt.Printf("📧 Sending welcome email to %s (%s)\n", payload.Name, payload.Email)

	// TODO: integrate real SMTP later
	return nil
}
