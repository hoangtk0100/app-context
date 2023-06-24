package main

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/mail"
	"github.com/hoangtk0100/app-context/core"
)

func main() {
	const cmpId = "gmail-sender"
	appCtx := appctx.NewAppContext(
		appctx.WithName("Demo Sending Email"),
		appctx.WithComponent(mail.NewEmailSender(cmpId)),
	)

	log := appCtx.Logger("service")

	if err := appCtx.Load(); err != nil {
		log.Error(err)
	}

	sender := appCtx.MustGet(cmpId).(core.EmailComponent)

	subject := "A test mail"
	content := `
	<h1>Hello world!</h1>
	<p>I'm coming :v</p>
	`
	to := []string{"<your-email>"}

	if err := sender.SendEmail(subject, content, to, nil, nil, nil); err != nil {
		log.Errorf(err, "Cannot send email ", err)
	}

	log.Infof("Sent email to %s", to)
}
