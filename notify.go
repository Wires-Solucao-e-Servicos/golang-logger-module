package logger

import (
	"fmt"
	"net/smtp"
	"strings"

	models "github.com/Wires-Solucao-e-Servicos/golang-logger-module/models"

	"github.com/jordan-wright/email"
)

func SendEmail(values models.Notification) error {

	clientName := ClientName

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

	e.From 		= SMTPConfig.From
	e.To 	 		= SMTPConfig.To
	e.Subject = subject
	e.Text 		= []byte(message.String())

	smtpAddress := fmt.Sprintf("%s:%d", SMTPConfig.Server, SMTPConfig.Port)
	smtpAuth := smtp.PlainAuth(
		"",
		SMTPConfig.Username,
		SMTPConfig.Password,
		SMTPConfig.Server,
	)

	err := e.Send(smtpAddress, smtpAuth)
	if err != nil {
		return err
	}
	
	return nil
}