package logger

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/jordan-wright/email"
)

func SendEmail(values Notification) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		config := GetSMTPConfig()
		clientName := GetClientName()

		date, event, location, details := values.Datetime, values.Code, values.Location, values.Details

		subject := fmt.Sprint("Logger Notification Service - " + clientName)

		var message strings.Builder

		fmt.Fprintf(&message,
			"The following event occurred at %s:\n\n"+
			"Date: %s\n"+
			"Event: %s\n"+
			"Location: %s\n"+
			"Details: %s\n", clientName, date, event, location, details,
		)

		e := email.NewEmail()

		e.From 		= config.From
		e.To 	 		= config.To
		e.Subject = subject
		e.Text 		= []byte(message.String())

		smtpAddress := fmt.Sprintf("%s:%d", config.Server, config.Port)
		smtpAuth := smtp.PlainAuth(
			"",
			config.Username,
			config.Password,
			config.Server,
		)

		tlsConfig := &tls.Config{
			ServerName: config.Server,
			InsecureSkipVerify: false,
		}

		done <- e.SendWithTLS(smtpAddress, smtpAuth, tlsConfig)

	} ()

	select {
		case err := <- done:
			return err
		case <- ctx.Done():
			return fmt.Errorf("email sending timeout after 30s")
	}

}