package smtp_sender

import (
	"fmt"
	"net/smtp"
	"time"
)

func SmtpSslAuthCheck(host string, port string, username string, password string, proxyUrl string, timeout time.Duration) error {
	auth := smtp.PlainAuth("", username, password, host)

	c, err := CreateSmtpSslConnection(proxyUrl, host, port, timeout)
	if err != nil {
		return fmt.Errorf("create_smtp_ssl_error: %s", err)
	}

	fmt.Println("connected")
	exit := make(chan bool)
	defer close(exit)
	go func() {
		select {
		case <-exit:
			return
		case <-time.After(timeout):
			c.Close()
			return
		}
	}()
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("auth_error: %s", err)
	}
	fmt.Println("auth")
	c.Close()
	return nil
}
