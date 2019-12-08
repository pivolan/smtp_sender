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
func SmtpTlsAuthCheck(host string, port string, username string, password string, proxyUrl string, timeout time.Duration) error {
	c, err := CreateSmtpTlsConnection(proxyUrl, host, port, timeout)
	if err != nil {
		return fmt.Errorf("create_smtp_tls_error: %s", err)
	}

	fmt.Println("connected")
	auth := smtp.PlainAuth("", username, password, host)
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
func SmtpPlainAuthCheck(host string, port string, username string, password string, proxyUrl string, timeout time.Duration) error {
	c, err := CreateSmtpPlainConnection(proxyUrl, host, port, timeout)
	if err != nil {
		return fmt.Errorf("create_smtp_plain_error: %s", err)
	}
	fmt.Println("connected")
	auth := smtp.PlainAuth("", username, password, host)
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
	c.Close()
	fmt.Println("auth")
	return nil
}
