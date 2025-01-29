package mailer

import (
	"bytes"
	"mime"
)

// TODO: Add SetHeader, SetBody, AddAlternative

type header map[string]string

type Encoding string

// This inherits all methods and fields of WordEncoder
type mimeEncoder struct {
	mime.WordEncoder
}

var (
	bEncoding = mimeEncoder{mime.BEncoding} // base64 encoding
	qEncoding = mimeEncoder{mime.QEncoding} // quoted-printable encoding
)

const (
	QuotedPrintable Encoding = "quoted-printable"
	Base64          Encoding = "base64"
	Unencoded       Encoding = "8bit" // Avoid encoding body of an email, but the headers will still be encoded
)

type Message struct {
	header   header
	encoding Encoding
	charset  string
	hEncoder mimeEncoder // Header encoder
	buf      bytes.Buffer
}

// Can be used as an ARGUMENT to configure an email
type MessageSetting func(m *Message)

func NewMessage(settings ...MessageSetting) *Message {
	m := &Message{
		header:   make(header),
		charset:  "UTF-8",
		encoding: QuotedPrintable,
	}
}
