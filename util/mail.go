package util

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

type Mail struct {
	userName   string
	passWord   string
	smtpServer string
	smtpPort   int
	mailType   string
	ssl        bool
}

func (m *Mail) Send(to, subject, body string) (err error) {
	header := make(map[string]string)
	auth := smtp.PlainAuth("", m.userName, m.passWord, m.smtpServer)

	header["From"] = m.userName
	header["To"] = to
	header["Subject"] = subject
	if m.mailType == "html" {
		header["Content-Type"] = "text/html; charset=UTF-8"
	} else {
		header["Content-Type"] = "text/plain; charset=UTF-8"
	}
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	message += "\r\n" + body

	sendTo := strings.Split(to, ";")
	if m.ssl {
		err = m.sendMailUsingTLS(fmt.Sprintf("%s:%d", m.smtpServer, m.smtpPort), auth, m.userName, sendTo, []byte(message))
	} else {
		err = smtp.SendMail(fmt.Sprintf("%s:%d", m.smtpServer, m.smtpPort), auth, m.userName, sendTo, []byte(message))
	}

	return err
}

func (m *Mail) sendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	c, err := m.dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func (m *Mail) dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
