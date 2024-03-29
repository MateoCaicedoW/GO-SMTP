package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/MateoCaicedoW/GO-SMTP/email"
	"github.com/MateoCaicedoW/GO-SMTP/server"
	"github.com/joho/godotenv"
)

//go:embed file.pdf
var file []byte

func main() {
	smtpServer := server.NewSMTP("smtp.gmail.com", "587", os.Getenv("SENDER_EMAIL"), os.Getenv("SENDER_PASSWORD"), "")

	email := email.Params{
		SenderName:      "SMTP Server",
		Sender:          "example@example.com",
		To:              []string{"example@example.com"},
		Cc:              []string{"example@example.com"},
		Subject:         "Test email with attachment",
		Body:            "This is a test email with attachment.",
		BodyContentType: "text/plain",
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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
