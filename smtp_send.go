package smtp_sender

import (
	"encoding/base64"
	"fmt"
	"time"
)

func SendSsl(host string, port string, username string, password string, proxyUrl string, timeout time.Duration, from string, to string, subject string, bodyHtml string, filename string, fileBody []byte) (err error) {
	c, err := SmtpSslAuth(host, port, username, password, proxyUrl, timeout)
	if err != nil {
		return
	}
	err = c.Mail(from)
	if err != nil {
		return
	}
	err = c.Rcpt(to)
	if err != nil {
		return
	}
	wc, err := c.Data()
	if err != nil {
		return
	}

	if len(filename) > 0 && len(fileBody) > 10 {
		_, err = fmt.Fprintf(wc, BODY_TPL_WITH_FILE, to, subject, bodyHtml, filename, base64.StdEncoding.EncodeToString(fileBody))
		if err != nil {
			return
		}
	} else {
		_, err = fmt.Fprintf(wc, BODY_TPL, to, subject, bodyHtml)
		if err != nil {
			return
		}
	}
	err = wc.Close()
	if err != nil {
		return
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		return
	}
	return
}

const BODY_TPL_WITH_FILE = `To: %s
Subject: %s
Content-Type: multipart/mixed; boundary="_===86128918====mx66.intranet.ru===_"

--_===86128918====mx66.intranet.ru===_
Content-Type: text/html; charset="utf-8"
Content-Transfer-Encoding: 8bit

%s

--_===86128918====mx66.intranet.ru===_
Content-Type: application/octet-stream
Content-Disposition: attachment;
 filename="%s"
Content-Transfer-Encoding: base64

%s`
const BODY_TPL = `To: %s
Subject: %s
Content-Type: multipart/mixed; boundary="_===86128918====mx66.intranet.ru===_"

--_===86128918====mx66.intranet.ru===_
Content-Type: text/html; charset="utf-8"
Content-Transfer-Encoding: 8bit

%s`
