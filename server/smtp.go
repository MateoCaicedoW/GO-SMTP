package server

import (
	"net/smtp"
)

type SMTPServer struct {
	host           string
	port           string
	serverMail     string
	serverPassword string
	identity       string
}

func (s *SMTPServer) ServerName() string {
	return s.host + ":" + s.port
}

func (s *SMTPServer) Auth() smtp.Auth {
	return smtp.PlainAuth(s.identity, s.serverMail, s.serverPassword, s.host)
}

func NewSMTP(host, port, serverMail, serverPassword, identity string) *SMTPServer {
	return &SMTPServer{
		host:           host,
		port:           port,
		serverMail:     serverMail,
		serverPassword: serverPassword,
		identity:       identity,
	}
}
