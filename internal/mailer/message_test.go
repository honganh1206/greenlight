package mailer

import (
	"bytes"
	"io"
	"testing"
)

type message struct {
	from    string
	to      []string
	content string
}

func TestMessage(t *testing.T) {
	m := NewMessage()
	m.SetAddressHeader("From", "from@example.com", "Señor From")
	m.SetHeader("To", m.FormatAddress("to@example.com", "Señor To"), "tobis@example.com")
	m.SetAddressHeader("Cc", "cc@example.com", "A, B")
	m.SetAddressHeader("X-To", "ccbis@example.com", "à, b")
	m.SetDateHeader("X-Date", now())
	m.SetHeader("X-Date-2", m.FormatDate(now()))
	m.SetHeader("Subject", "¡Hola, señor!")
	m.SetHeaders(map[string][]string{
		"X-Headers": {"Test", "Café"},
	})
	m.SetBody("text/plain", "¡Hola, señor!")

	want := &message{
		from: "from@example.com",
		to: []string{
			"to@example.com",
			"tobis@example.com",
			"cc@example.com",
		},
		content: "From: =?UTF-8?q?Se=C3=B1or_From?= <from@example.com>\r\n" +
			"To: =?UTF-8?q?Se=C3=B1or_To?= <to@example.com>, tobis@example.com\r\n" +
			"Cc: \"A, B\" <cc@example.com>\r\n" +
			"X-To: =?UTF-8?b?w6AsIGI=?= <ccbis@example.com>\r\n" +
			"X-Date: Wed, 25 Jun 2014 17:46:00 +0000\r\n" +
			"X-Date-2: Wed, 25 Jun 2014 17:46:00 +0000\r\n" +
			"X-Headers: Test, =?UTF-8?q?Caf=C3=A9?=\r\n" +
			"Subject: =?UTF-8?q?=C2=A1Hola,_se=C3=B1or!?=\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"Content-Transfer-Encoding: quoted-printable\r\n" +
			"\r\n" +
			"=C2=A1Hola, se=C3=B1or!",
	}

	testMessage(t, m, 0, want)
}

// ////////////////////////// HELPER FUNCTIONS
func testMessage(t *testing.T, m *Message, bCount int, want *message) {
	// stubSendMail satisfies the Sender interface
	// Thus we can pass stubSendEmail for the parameter s of type Sender
	err := Send(stubSendMail(t, bCount, want), m)
	if err != nil {
		t.Error(err)
	}
}

func stubSendMail(t *testing.T, bCount int, want *message) SendFunc {
	return func(from string, to []string, m io.WriterTo) error {
		if from != want.from {
			t.Fatalf("invalid from, got %q, want %q", from, want.from)
		}

		if len(to) != len(want.to) {
			t.Fatalf("invalid recipient count, \ngot %d: %q\nwant %d: %q", len(to), to, len(want.to), want.to)
		}

		for i := range want.to {
			if to[i] != want.to[i] {
				t.Fatalf("invalid recipient, \ngot: %q\nwant: %q", to, want.to)
			}
		}

		buf := new(bytes.Buffer)
		_, err := m.WriteTo(buf)
		if err != nil {
			t.Error(err)
		}

		// got := buf.String()

		// wantMsg := string("Mime-Version: 1.0\r\n" +
		// 	"Date: Wed, 25 Jun 2014 17:46:00 +0000\r\n" +
		// 	want.content)

		// When we need attachment, MIME boundaries will be > 0
		// TODO: Add this later when we do attachments
		// if bCount > 0 {
		// 	boundaries := getBoundaries(t, bCount, got)
		// 	for i, b := range boundaries {
		// 		wantMsg = strings.Replace(wantMsg, "_BOUNDARY_"+strconv.Itoa(i+1)+"_", b, -1)
		// 	}
		// }

		// TODO: Add this later
		// compareBodies(t, got, wantMsg)

		return nil
	}
}
