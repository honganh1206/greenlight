package mailer

import (
	"fmt"
	"net"
)

// Establish and maintain connections
type Dialer struct {
	Host     string
	Port     int
	Username string
	Password string
	// TODO: Add SSL, Auth and TLSConfig if needed
	// If SSL, make sure SSL if port == 465
}

// type smtpClient interface {

// Stub for testing
var (
	netDialTimeout = net.DialTimeout
)

func NewDialer(host string, port int, username, password string) *Dialer {
	return &Dialer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func addr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
