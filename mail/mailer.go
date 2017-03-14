package mail

import (
    "bytes"
    "net/smtp"
	"text/template"
	
    log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/SCKelemen/Cassius/common"
	"github.com/SCKelemen/Cassius/log"
)

type Mailer interface {
	SendPasswordResetMail(to, token string) error
}

var passwordResetMailTmpl = template.Must(template.New("passwordResetMailTemplate").Parse("To: {{.To}}\r\nSubject: The HackStack Password Reset\r\n\r\nClick the following link to reset password: {{.RootURL}}/#resetPassword?token={{.Token}}"))

type SMTPMailer struct {
	ServerAddr string
	Auth       smtp.Auth
	From       string
	rootURL    string
	logger     log15.Logger
}

func NewMailer(config common.AppConfig, logger log15.Logger) (Mailer, error) {
    if config.SmtpActive {
        logger = logger.New("module", "mail")
        log.SetFilterHandler("warn", logger, log15.StdoutHandler)

        auth := smtp.PlainAuth("", config.SmtpUser, config.SmtpPassword, config.SmtpServer)

        mailer := &SMTPMailer{
            ServerAddr: config.SmtpServer + ":" + config.SmtpPort,
            Auth:       auth,
            From:       config.SmtpFromAddress,
            rootURL:    config.SmtpRootUrl,
            logger:     logger,
	    }

        return mailer, nil
    }

    return nil, nil
}

func (m *SMTPMailer) SendPasswordResetMail(to, token string) error {
	var data = struct {
		RootURL string
		To      string
		Token   string
	}{
		RootURL: m.rootURL,
		To:      to,
		Token:   token,
	}

	buf := &bytes.Buffer{}
	err := passwordResetMailTmpl.Execute(buf, data)
	if err != nil {
		return err
	}

	err = smtp.SendMail(m.ServerAddr, m.Auth, m.From, []string{to}, buf.Bytes())
	if err != nil {
		m.logger.Error("SendPasswordResetEmail failed", "to", to, "error", err)
		return err
	}

	m.logger.Info("SendPasswordResetEmail", "to", to)
	return nil
}
