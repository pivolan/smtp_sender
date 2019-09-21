package smtp_sender

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/smtp"
	"net/url"
	"time"
)

func CreateProxyConnection(proxyUrl string, host string, port string, timeout time.Duration) (conn net.Conn, err error) {
	servername := host + ":" + port
	uri, err := url.Parse(proxyUrl)
	dialer, err := proxy.FromURL(uri, proxy.Direct)
	if err != nil {
		fmt.Printf("FromURL failed: %v\n", err)
		return nil, fmt.Errorf("proxy_parse_error: %s", err)
	}
	fmt.Println("start dial", servername, proxyUrl, port, time.Now().Format(time.StampMilli))
	if f, ok := dialer.(proxy.ContextDialer); ok {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		conn, err = f.DialContext(ctx, "tcp", servername)
	} else {
		conn, err = dialer.Dial("tcp", servername)
	}
	if err != nil {
		return nil, fmt.Errorf("dial_error: %s", err)
	}
	fmt.Println("end dial", servername, proxyUrl, port, time.Now().Format(time.StampMilli))
	return
}

func CreateSmtpPlainConnection(proxyUrl string, host string, port string, timeout time.Duration) (c *smtp.Client, err error) {
	conn, err := CreateProxyConnection(proxyUrl, host, port, timeout)
	if err != nil {
		return
	}
	exit := make(chan bool)
	defer close(exit)
	go func() {
		for {
			select {
			case <-exit:
				return
			case <-time.After(timeout):
				conn.Close()
				return
			}
		}
	}()
	servername := host + ":" + port
	c, err = smtp.NewClient(conn, servername)
	return
}
func CreateSmtpSslConnection(proxyUrl string, host string, port string, timeout time.Duration) (c *smtp.Client, err error) {
	conn, err := CreateProxyConnection(proxyUrl, host, port, timeout)
	if err != nil {
		return
	}
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	exit := make(chan bool)
	go func() {
		for {
			select {
			case <-exit:
				return
			case <-time.After(timeout):
				conn.Close()
				return
			}
		}
	}()
	fmt.Println("start ssl", proxyUrl, port, host, time.Now().Format(time.StampMilli))
	conn = tls.Client(conn, tlsconfig)
	fmt.Println("end ssl", proxyUrl, port, host, time.Now().Format(time.StampMilli))
	close(exit)
	conn.SetDeadline(time.Now().Add(timeout * 2))
	c, err = smtp.NewClient(conn, host)
	return
}
func CreateSmtpTlsConnection(proxyUrl string, host string, port string, timeout time.Duration) (c *smtp.Client, err error) {
	conn, err := CreateProxyConnection(proxyUrl, host, port, timeout)
	if err != nil {
		return
	}
	conn.SetDeadline(time.Now().Add(timeout * 2))
	c, err = smtp.NewClient(conn, host)
	if err != nil {
		return nil, fmt.Errorf("smtp_connect error: %s", err)
	}
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	exit2 := make(chan bool)
	go func() {
		for {
			select {
			case <-exit2:
				return
			case <-time.After(timeout):
				c.Close()
				return
			}
		}
	}()
	fmt.Println("start tls", port, host, proxyUrl, time.Now().Format(time.StampMilli))
	err = c.StartTLS(tlsconfig)
	close(exit2)
	fmt.Println("end tls", port, host, proxyUrl, time.Now().Format(time.StampMilli))
	if err != nil {
		return nil, fmt.Errorf("tls_error: %s", err)
	}

	return
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
func SmtpSslAuth(host string, port string, username string, password string, proxyUrl string, timeout time.Duration) (c *smtp.Client, err error) {
	auth := smtp.PlainAuth("", username, password, host)

	c, err = CreateSmtpSslConnection(proxyUrl, host, port, timeout)
	if err != nil {
		return nil, fmt.Errorf("create_smtp_ssl_error: %s", err)
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
		return nil, fmt.Errorf("auth_error: %s", err)
	}
	fmt.Println("auth")
	return
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
