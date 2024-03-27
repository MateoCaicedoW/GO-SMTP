package main

import (
	_ "embed"
	"smtp/email"
	"smtp/server"
)

const (
	senderEmail    = "smtpmateo@gmail.com"
	senderPassword = "rdej kqnl pczk ixve"
	host           = "smtp.gmail.com"
	port           = "587"
)

//go:embed file.pdf
var file []byte

func main() {
	smtpServer := server.NewSMTP(host, port, senderEmail, senderPassword, "")

	email := email.Params{
		SenderName: "Mateo Caicedo",
		Sender:     "caicedomateo9@gmail.com",
		To:         []string{"caicedomateo9@gmail.com"},
		Cc:         []string{"mcaicedo@wawand.co"},
		Subject:    "La chapa que vibra",
		Body:       "Hola, este es un mensaje de prueba",
		Attachments: []email.Attachment{
			{
				FileName: "file.pdf",
				Content:  file,
			},
		},
	}

	err := email.Send(smtpServer)
	if err != nil {
		panic(err)
	}

}
