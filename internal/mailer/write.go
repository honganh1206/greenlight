package mailer

import (
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"strings"
)

// Different writers to handle different levels of MIME multipart content
type messageWriter struct {
	w          io.Writer            // Main writer for main message body (no boundaries)
	n          int64                // Total bytes written counter -- satisfy the io.Writer interface
	writers    [3]*multipart.Writer // Generate multipart messages
	partWriter io.Writer            // Part-specific writer for multipart message (yes boundaries)
	depth      uint8
	err        error
}

// Required by RFC 2045
// Legacy email systems and protocols have line length restrictions
// Old systems cannot handle lines longer than 78 characters
// So 76 it is to leave room for line wrapping/indentation
const maxLineLen = 76

type base64LineWriter struct {
	w       io.Writer
	lineLen int
}

func newBase64LineWriter(w io.Writer) *base64LineWriter {
	return &base64LineWriter{w: w}
}

func (w *base64LineWriter) Write(p []byte) (int, error) {
	n := 0
	for len(p)+w.lineLen > maxLineLen {
		// Write from current position up to the maximum line length,
		// taking into account how many characters are already on the current line (w.lineLen)
		w.w.Write(p[:maxLineLen-w.lineLen])
		// Add CRLF for line wrapping
		w.w.Write([]byte("\r\n"))
		// Update p to contain only the remaining unwritten bytes
		// by slicing off what we just wrote
		p = p[maxLineLen-w.lineLen:]
		// Add to our running count of bytes written
		// (this is how many bytes we just wrote before the line break)
		n += maxLineLen - w.lineLen
	}

	w.w.Write(p)
	w.lineLen += len(p)

	return n + len(p), nil
}

// Implementation of WriteTo for Message struct
// To satisfy the io.Writer interface
func (m *Message) WriteTo(w io.Writer) (int64, error) {
	mw := &messageWriter{w: w}
	mw.writeMessage(m)
	return mw.n, mw.err

}

// Implementation of Write for messageWriter struct
// To satisfy io.Writer interface
func (w *messageWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, errors.New("cannot write as writer has errors")
	}

	var n int
	n, w.err = w.w.Write(p) // Write to the underlying writer
	w.n += int64(n)

	return n, w.err
}

func (w *messageWriter) writeMessage(m *Message) {
	if _, ok := m.header["Mime-Version"]; !ok {
		w.writeString("Mime-Version: 1.0\r\n")
	}

	if _, ok := m.header["Date"]; !ok {
		w.writeHeader("Date", m.FormatDate(now()))
	}

	w.writeHeaders(m.header)

	// MIME level 2 is enough
	// multipart/alternative is also enough
	// Go form mixed (outermost) -> related -> alternative as the outer ones nesting the inner ones
	if m.hasAlternativePart() {
		w.openMultiPart("alternative")
	}

	for _, part := range m.parts {
		w.writePart(part, m.charset)
	}

	if m.hasAlternativePart() {
		w.closeMultipart()
	}

}

func (w *messageWriter) writePart(p *part, charset string) {
	w.writeHeaders(map[string][]string{
		"Content-Type":              {p.contentType + ";, charset=" + charset},
		"Content-Transfer-Encoding": {string(p.encoding)},
	})
	w.writeBody(p.copier, p.encoding)
}

func (w *messageWriter) writeHeaders(h map[string][]string) {
	// For top-level headers like from/to/Subject...
	if w.depth == 0 {
		for k, v := range h {
			if k != "Bcc" {
				w.writeHeader(k, v...)
			}
		}
	} else {
		// For MIME part headers like Content-Type/Content-Disposition/...
		w.createPart(h)
	}
}

func (w *messageWriter) writeBody(f func(io.Writer) error, enc Encoding) {
	var subWriter io.Writer
	if w.depth == 0 {
		// Separator between headers and body
		w.writeString("\r\n")
		subWriter = w.w
	} else {
		// Specifically for multipart message
		subWriter = w.partWriter
	}

	if enc == Base64 {
		wc := base64.NewEncoder(base64.StdEncoding, newBase64LineWriter(subWriter))
		w.err = f(wc)
		wc.Close()
	} else if enc == Unencoded {
		w.err = f(subWriter)
	} else {
		wc := newQPWriter(subWriter)
		w.err = f(wc)
		wc.Close()
	}
}

func (w *messageWriter) writeHeader(k string, v ...string) {
	w.writeString(k)
	if len(v) == 0 {
		w.writeString(":\r\n")
		return
	}
	w.writeString(": ")

	// For simplicity (and per RFC 5322) we use the 76-char limit
	charsLeft := 76 - len(k) - len(": ")

	for i, s := range v {
		// Insert a new line if the current line is too long already
		if charsLeft < 1 {
			if i == 0 {
				w.writeString("\r\n ")
			} else {
				w.writeString(", \r\n")
			}
			charsLeft = 75
		} else if i != 0 {
			// Otherwise we just list the values
			w.writeString(", ")
			charsLeft -= 2
		}

		// Fold the header content if it is too long
		for len(s) > charsLeft {
			s = w.writeLine(s, charsLeft)
			charsLeft = 75
		}
		w.writeString(s)

	}
}

func (w *messageWriter) writeLine(s string, charsLeft int) string {
	// Scenario 1: Handle existing newlines within the character limit
	// Example input: "Hello\nWorld" with charsLeft = 10
	// Result: Writes "Hello\n" and returns "World"
	if i := strings.IndexByte(s, '\n'); i != -1 && i < charsLeft {
		w.writeString(s[:i+1])
		return s[i+1:]
	}

	// Scenario 2: Break at word boundaries within the character limit
	// Searches backwards from the character limit to find a space
	// Example input: "This is a long sentence" with charsLeft = 10
	// Result: Writes "This is a" + "\r\n " and returns "long sentence"
	for i := charsLeft - 1; i >= 0; i-- {
		if s[i] == ' ' {
			w.writeString(s[:i])
			w.writeString("\r\n ")
			return s[i+1:]
		}
	}

	// Scenario 3: Look for space or newline beyond the standard limit (75 chars)
	// This handles cases where no suitable break point was found within charsLeft
	// Example 1: "ThisIsAVeryLongWordFollowedBy Space"
	// Result: Breaks at the first space even if it's beyond the limit
	//
	// Example 2: "VeryLongWord\nNextLine"
	// Result: Breaks at the newline character
	for i := 75; i < len(s); i++ {
		if s[i] == ' ' {
			w.writeString(s[:i])
			w.writeString("\r\n ")
			return s[i+1:]
		}
		if s[i] == '\n' {
			w.writeString(s[:i+1])
			return s[i+1:]
		}
	}

	// Scenario 4: No break points found in the entire string
	// Example input: "ThisIsOneVeryLongWordWithNoSpacesOrNewlines"
	// Result: Writes the entire string as is
	w.writeString(s)
	return ""
}

func (w *messageWriter) writeString(s string) {
	n, _ := io.WriteString(w.w, s)
	w.n += int64(n)
}

func (w *messageWriter) createPart(h map[string][]string) {
	w.partWriter, w.err = w.writers[w.depth-1].CreatePart(h)
}

// Delimit where each part begins and ends in the message body
func (w *messageWriter) openMultiPart(mimeType string) {
	mw := multipart.NewWriter(w)
	contentType := "multipart/" + mimeType + ";\r\n boundary=" + mw.Boundary()
	w.writers[w.depth] = mw // Store the writer at the current depth level

	if w.depth == 0 {
		// At this depth level it is all about setting up the message structure
		//
		w.writeHeader("Content-Type", contentType)
		w.writeString("\r\n")
	} else {
		// At this depth we need boundary markers, so we need to create a proper MIME part
		w.createPart(map[string][]string{
			"Content-Type": {contentType},
		})
	}
	w.depth++

}

func (m *Message) hasAlternativePart() bool {
	return len(m.parts) > 1
}

// Close the multipart writer at the current depth level
// From innermost part to the outermost one
func (w *messageWriter) closeMultipart() {
	if w.depth > 0 {
		w.writers[w.depth-1].Close()
		w.depth--
	}
}
