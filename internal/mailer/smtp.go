package mailer

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// Establish and maintain connections
type Dialer struct {
	Host      string
	Port      int
	Username  string
	Password  string
	LocalName string // hostname to use in HELO/EHLO command
	// If SSL, make sure port == 465
	SSL       bool // Should be false in most cases. Prefer STARTTLS instead
	TLSConfig *tls.Config
	Auth      smtp.Auth
	Timeout   time.Duration
}

type smtpSender struct {
	smtpClient // struct inside another struct
	d          *Dialer
}

type smtpClient interface {
	Hello(string) error              // HELO/EHLO command
	Extension(string) (bool, string) // Check if server supports specific SMTP extension
	StartTLS(*tls.Config) error      // Upgrade connection to StartTLS
	Auth(smtp.Auth) error            // Authenticate with the server
	Mail(string) error               // Set sender address
	Rcpt(string) error               // Add recipient address
	Data() (io.WriteCloser, error)   // Open data stream for message content
	Quit() error                     // Send QUIT command to end session
	Close() error                    // Close the connection
}

// Stub for testing
var (
	netDialTimeout = net.DialTimeout
	smtpNewClient  = func(conn net.Conn, host string) (smtpClient, error) {
		return smtp.NewClient(conn, host)
	}
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

func (d *Dialer) DialAndSend(m ...*Message) error {
	s, err := d.Dial()

	if err != nil {
		return err
	}

	defer s.Close()

	return Send(s, m...)
}

// Dial and authenticate to an SMTP server
// Auto close the sender when done using
func (d *Dialer) Dial() (SendCloser, error) {
	conn, err := netDialTimeout("tcp", addr(d.Host, d.Port), d.Timeout)
	if err != nil {
		return nil, err
	}

	c, err := smtpNewClient(conn, d.Host)
	if err != nil {
		return nil, err
	}

	if d.LocalName != "" {
		if err := c.Hello(d.LocalName); err != nil {
			return nil, err
		}
	}

	if !d.SSL {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(d.tlsConfig()); err != nil {
				c.Close()
				return nil, err
			}
		}
	}

	if d.Auth == nil && d.Username != "" {
		if ok, auths := c.Extension("AUTH"); ok {
			if strings.Contains(auths, "CRAM-MD5") {
				d.Auth = smtp.CRAMMD5Auth(d.Username, d.Password)
			} else if strings.Contains(auths, "LOGIN") &&
				!strings.Contains(auths, "PLAIN") {
				d.Auth = &loginAuth{
					username: d.Username,
					password: d.Password,
					host:     d.Host,
				}
			} else {
				d.Auth = smtp.PlainAuth("", d.Username, d.Password, d.Host)
			}
		}
	}

	if d.Auth != nil {
		if err = c.Auth(d.Auth); err != nil {
			c.Close()
			return nil, err
		}
	}

	return &smtpSender{c, d}, nil
}

// Implementation of Send for smtpSender
func (c *smtpSender) Send(from string, to []string, msg io.WriterTo) error {
	if err := c.Mail(from); err != nil {
		if err == io.EOF {
			// Possible timeout
			// Try to establish a new connection
			sc, derr := c.d.Dial()
			if derr == nil {
				if s, ok := sc.(*smtpSender); ok {
					// Replace the old connection with the new one
					// And retry sending the message recursively
					*c = *s
					return c.Send(from, to, msg)
				}
			}
		}
		return err
	}

	// Check for valid recipients
	for _, addr := range to {
		if err := c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	// Write the message to the data stream
	if _, err = msg.WriteTo(w); err != nil {
		w.Close()
		return err
	}

	return w.Close()
}

func (d *Dialer) tlsConfig() *tls.Config {
	if d.TLSConfig == nil {
		return &tls.Config{ServerName: d.Host}
	}
	return d.TLSConfig
}
