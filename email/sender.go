package email

import "gopkg.in/gomail.v2"

const (
	SmtpVmailHost = "smtp.viettel.com.vn"
	SmtpGmailPort = 465
)

type EmailSender interface {
	Send(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachments []string,
	) error
}

type VmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewVmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &VmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *VmailSender) Send(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachments []string,
) error {
	m := gomail.NewMessage()
	m.SetHeader("From", sender.fromEmailAddress)
	m.SetHeader("To", to...)
	m.SetHeader("Cc", cc...)
	m.SetHeader("Bcc", bcc...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	// attachment files
	for _, filepath := range attachments {
		m.Attach(filepath)
	}

	dial := gomail.NewDialer(SmtpVmailHost, SmtpGmailPort, sender.fromEmailAddress, sender.fromEmailPassword)
	return dial.DialAndSend(m)
}
