package logger

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	models "github.com/Wires-Solucao-e-Servicos/golang-logger-module/models"

	"github.com/jordan-wright/email"
)

func SendEmail(values models.Notification) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		config := GetSMTPConfig()
		clientName := GetClientName()

		date, event, location, details := values.Datetime, values.Code, values.Location, values.Details

		subject := fmt.Sprint("Defense Backup Service Notification - " + clientName)

		var message strings.Builder

		fmt.Fprintf(&message,
			"The following event occurred in the Defense Backup Service at %s:\n\n"+
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

		done <- e.Send(smtpAddress, smtpAuth)

	} ()

	select {
		case err := <- done:
			return err
		case <- ctx.Done():
			return fmt.Errorf("email sending timeout after 30s")
	}

}