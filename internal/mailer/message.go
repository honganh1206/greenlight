package mailer

import (
	"bytes"
	"io"
	"mime"
	"time"
)

// TODO: Add SetBody, AddAlternative

type header map[string][]string

type Encoding string

// This inherits all methods and fields of WordEncoder
type mimeEncoder struct {
	mime.WordEncoder
}

var (
	bEncoding = mimeEncoder{mime.BEncoding} // base64 encoding
	qEncoding = mimeEncoder{mime.QEncoding} // quoted-printable encoding
)

// For testing
var now = time.Now

const (
	QuotedPrintable Encoding = "quoted-printable"
	Base64          Encoding = "base64"
	Unencoded       Encoding = "8bit" // Avoid encoding body of an email, but the headers will still be encoded
)

type Message struct {
	header   header
	parts    []*part
	encoding Encoding
	charset  string
	hEncoder mimeEncoder // Header encoder
	buf      bytes.Buffer
}

// Represent different parts/sections of an email message in MIME format
// message := NewMessage()
// message.SetBody("text/html", `
//
//	<html>
//	    <body>
//	        <h1>Welcome</h1>
//	        <img src="cid:logo.png">
//	    </body>
//	</html>
//
// `)
// message.Embed("logo.png")
// message.Attach("report.pdf")
type part struct {
	contentType string
	copier      func(io.Writer) error
	encoding    Encoding
}

// Configure the part added to the message
type PartSetting func(*part)

// A function type that takes a pointer and modify it
// Can be used as an ARGUMENT to configure an email
type MessageSetting func(m *Message)

func NewMessage(settings ...MessageSetting) *Message {
	m := &Message{
		header:   make(header),
		charset:  "UTF-8",
		encoding: QuotedPrintable,
	}

	m.applySettings(settings)

	if m.encoding == Base64 {
		m.hEncoder = bEncoding
	} else {
		m.hEncoder = qEncoding
	}

	return m
}

func (m *Message) SetHeader(field string, value ...string) {
	m.encodeHeader(value)
	m.header[field] = value
}

func (m *Message) SetHeaders(fields map[string][]string) {
	for k, v := range fields {
		m.SetHeader(k, v...)
	}
}

func (m *Message) SetAddressHeader(field, address, name string) {
	m.header[field] = []string{m.FormatAddress(address, name)}
}

func (m *Message) FormatAddress(address, name string) string {
	if name == "" {
		return address
	}

	enc := m.encodeString(name)

	if enc == name {
		// Wrap the name around the ""
		// `"John Doe" <john@example.com>`
		m.buf.WriteByte('"')
		for i := 0; i < len(name); i++ {
			b := name[i]
			if b == '\\' || b == '"' {
				m.buf.WriteByte('\\')
			}
			m.buf.WriteByte(b)
		}
		m.buf.WriteByte('"')
	} else if hasSpecials(name) {
		// `=?UTF-8?B?Sm9zw6k=?= <jose@example.com>`
		m.buf.WriteString(bEncoding.Encode(m.charset, name))
	} else {
		// `john@example.com` (when name is empty)
		m.buf.WriteString(enc)
	}
	m.buf.WriteString(" <")
	m.buf.WriteString(address)
	m.buf.WriteByte('>')

	addr := m.buf.String()
	m.buf.Reset()

	return addr
}

func (m *Message) SetDateHeader(field string, date time.Time) {
	m.header[field] = []string{m.FormatDate(date)}
}

func (m *Message) FormatDate(date time.Time) string {
	return date.Format(time.RFC1123Z)
}

// Replace any content set by SetBody, AddAlternative or AddAlternativeWriter
func (m *Message) SetBody(contentType, body string, settings ...PartSetting) {
	m.parts = []*part{m.newPart(contentType, newCopier(body), settings)}
}

// TODO: Add AddAlternative

//////////////// HELPER FUNCTIONS

func (m *Message) encodeHeader(values []string) {
	for i := range values {
		values[i] = m.encodeString(values[i])
	}
}

func (m *Message) encodeString(value string) string {
	return m.hEncoder.Encode(m.charset, value)
}

func hasSpecials(text string) bool {
	for i := 0; i < len(text); i++ {
		switch c := text[i]; c {
		case '(', ')', '<', '>', '[', ']', ':', ';', '@', '\\', ',', '.', '"':
			return true
		}
	}

	return false
}

func (m *Message) applySettings(settings []MessageSetting) {
	for _, s := range settings {
		s(m)
	}
}

func (m *Message) newPart(contentType string, f func(io.Writer) error, settings []PartSetting) *part {
	p := &part{
		contentType: contentType,
		copier:      f,
		encoding:    m.encoding,
	}

	for _, s := range settings {
		// Modify the part here
		s(p)
	}

	return p
}

func newCopier(s string) func(io.Writer) error {
	return func(w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	}
}
