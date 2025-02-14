package mailer

import (
	"errors"
	"fmt"
	"io"
	"net/mail"
)

type Sender interface {
	Send(from string, to []string, msg io.WriterTo) error
}

// Group the Send and Close methods
type SendCloser interface {
	Sender
	Close() error
}

// Function type that sends emails to the given address
// This type satisfies the Sender interface
type SendFunc func(from string, to []string, msg io.WriterTo) error

// Explicitly make the `SendFunc` type to implement the Sender interface
func (f SendFunc) Send(from string, to []string, msg io.WriterTo) error {
	return f(from, to, msg)
}

func Send(s Sender, msg ...*Message) error {
	for i, m := range msg {
		if err := send(s, m); err != nil {
			return fmt.Errorf("could not send email %d: %v", i+1, err)
		}
	}

	return nil
}

func send(s Sender, m *Message) error {
	from, err := m.getFrom()

	if err != nil {
		return err
	}

	to, err := m.getRecipients()

	if err != nil {
		return err
	}

	if err := s.Send(from, to, m); err != nil {
		return err
	}

	return nil
}

func (m *Message) getFrom() (string, error) {
	from := m.header["Sender"]
	if len(from) == 0 {
		from = m.header["From"]
		if len(from) == 0 {
			return "", errors.New(`invalid message, "From" field is absent`)
		}
	}

	return parseAddress(from[0])
}

func parseAddress(field string) (string, error) {
	addr, err := mail.ParseAddress(field)
	if err != nil {
		return "", fmt.Errorf("invalid address %q: %v", field, err)
	}
	return addr.Address, nil
}

func (m *Message) getRecipients() ([]string, error) {
	fields := []string{"To", "Cc", "Bcc"}

	if m.header == nil {
		return nil, errors.New("no header information available")
	}

	n := 0
	for _, field := range fields {
		addresses, exists := m.header[field]
		if !exists {
			continue
		}

		n += len(addresses)
	}

	if n == 0 {
		return nil, errors.New("no recipients specified")
	}

	list := make([]string, 0, n)

	for _, field := range fields {
		addresses, exists := m.header[field]
		if !exists {
			continue
		}

		for _, a := range addresses {
			addr, err := parseAddress(a)
			if err != nil {
				return nil, err
			}
			list = addAddress(list, addr)
		}
	}

	return list, nil
}

func addAddress(list []string, addr string) []string {
	// Check for duplication
	for _, a := range list {
		if addr == a {
			return list
		}
	}
	return append(list, addr)
}
