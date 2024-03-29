package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime"
	"net/smtp"
	"path/filepath"

	"strings"
	"time"

	"github.com/MateoCaicedoW/GO-SMTP/server"
)

type Params struct {
	Sender          string
	SenderName      string
	Subject         string
	To              []string
	Cc              []string
	ReplyTo         []string
	Bcc             []string
	Body            string
	BodyContentType string
	Attachments     []Attachment
}

func (m *Params) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + m.SenderName + " <" + m.Sender + ">" + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC1123Z) + "\r\n")

	buf.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")
	if len(m.Cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(m.Cc, ",") + "\r\n")
	}

	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(m.Subject)) + "?="
	buf.WriteString("Subject: " + subject + "\r\n")

	if len(m.ReplyTo) > 0 {
		buf.WriteString("Cc: " + strings.Join(m.ReplyTo, ",") + "\r\n")
	}

	buf.WriteString("MIME-Version: 1.0\r\n")

	boundary := time.Now().Format("20060102-150405")

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.BodyContentType))
	buf.WriteString(m.Body)
	buf.WriteString("\r\n")

	if len(m.Attachments) > 0 {
		for _, attachment := range m.Attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			ext := filepath.Ext(attachment.FileName)
			mimetype := mime.TypeByExtension(ext)
			if mimetype != "" {
				mime := fmt.Sprintf("Content-Type: %s\r\n", mimetype)
				buf.WriteString(mime)
			} else {
				buf.WriteString("Content-Type: application/octet-stream\r\n")
			}
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")

			buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
			buf.WriteString(coder.EncodeToString([]byte(attachment.FileName)))
			buf.WriteString("?=\"\r\n\r\n")

			b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Content)))
			base64.StdEncoding.Encode(b, attachment.Content)

			// write base64 content in lines of up to 76 chars
			for i, l := 0, len(b); i < l; i++ {
				buf.WriteByte(b[i])
				if (i+1)%76 == 0 {
					buf.WriteString("\r\n")
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
func (e *Params) Send(smtpServer *server.SMTPServer) error {
	auth := smtpServer.Auth()
	serverName := smtpServer.ServerName()
	bytes := e.Bytes()

	err := smtp.SendMail(serverName, auth, e.Sender, e.To, bytes)
	if err != nil {
		return err
	}

	if len(e.Cc) > 0 {
		err = smtp.SendMail(serverName, auth, e.Sender, e.Cc, bytes)
		if err != nil {
			return err
		}
	}

	if len(e.Bcc) > 0 {
		err = smtp.SendMail(serverName, auth, e.Sender, e.Bcc, bytes)
		if err != nil {
			return err
		}
	}

	return nil
}
