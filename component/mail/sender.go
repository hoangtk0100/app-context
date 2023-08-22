package mail

import (
	"fmt"
	"net/smtp"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/jordan-wright/email"
	"github.com/spf13/pflag"
)

const (
	defaultSMTPServer = "smtp.gmail.com"
	defaultSMTPPort   = 587
)

type emailOpt struct {
	smtpServer string
	smtpPort   int
}

type emailSender struct {
	id       string
	name     string
	address  string
	password string
	logger   appctx.Logger
	*emailOpt
}

func NewEmailSender(id string) *emailSender {
	return &emailSender{
		id:       id,
		emailOpt: new(emailOpt),
	}
}

func (es *emailSender) ID() string {
	return es.id
}

func (es *emailSender) InitFlags() {
	pflag.StringVar(&es.name,
		"email-sender-name",
		"",
		"Email sender name",
	)

	pflag.StringVar(&es.address,
		"email-sender-address",
		"",
		"Email sender address",
	)

	pflag.StringVar(&es.password,
		"email-sender-password",
		"",
		"Email sender password",
	)

	pflag.StringVar(&es.smtpServer,
		"email-smtp-server",
		defaultSMTPServer,
		fmt.Sprintf("Email SMTP server - Default: %s (Gmail)", defaultSMTPServer),
	)

	pflag.IntVar(&es.smtpPort,
		"email-smtp-port",
		defaultSMTPPort,
		fmt.Sprintf("Email SMTP port - Default: %d", defaultSMTPPort),
	)
}

func (es *emailSender) Run(ac appctx.AppContext) error {
	es.logger = ac.Logger(es.id)
	return nil
}

func (es *emailSender) Stop() error {
	return nil
}

func (es *emailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachments []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", es.name, es.address)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachments {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpServerAddress := fmt.Sprintf("%s:%d", es.smtpServer, es.smtpPort)
	smtpAuth := smtp.PlainAuth("", es.address, es.password, es.smtpServer)

	return e.Send(smtpServerAddress, smtpAuth)
}
